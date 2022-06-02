package framework

import (
	"github.com/opensourceways/community-robot-lib/config"
	"github.com/opensourceways/community-robot-lib/gitlabclient"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const (
	logFieldOrg    = "org"
	logFieldRepo   = "repo"
	logFieldURL    = "url"
	logFieldAction = "action"
)

type dispatcher struct {
	agent *config.ConfigAgent

	h handlers

	// Tracks running handlers for graceful shutdown
	wg sync.WaitGroup
}

func (d *dispatcher) Wait() {
	d.wg.Wait() // Handle remaining requests
}

func (d *dispatcher) Dispatch(eventType string, payload []byte, l *logrus.Entry) error {
	eventType2 := gitlab.EventType(eventType)
	hook, err := gitlab.ParseHook(eventType2, payload)
	if err != nil {
		return err
	}

	switch hook := hook.(type) {
	case *gitlab.IssueEvent:
		d.wg.Add(1)
		go d.handleIssueEvent(hook, l)
	case *gitlab.MergeEvent:
		d.wg.Add(1)
		go d.handleMergeEvent(hook, l)
	case *gitlab.PushEvent:
		d.wg.Add(1)
		go d.handlePushEvent(hook, l)
	case *gitlab.IssueCommentEvent:
		d.wg.Add(1)
		go d.handleIssueCommentEvent(hook, l)
	case *gitlab.MergeCommentEvent:
		d.wg.Add(1)
		go d.handleMergeCommentEvent(hook, l)
	case *gitlab.CommitCommentEvent:
		d.wg.Add(1)
		go d.handleCommitCommentEvent(hook, l)
	default:
		l.Debug("Ignoring unknown event type")
	}

	return nil
}

func (d *dispatcher) getConfig() config.Config {
	_, c := d.agent.GetConfig()

	return c
}

func (d *dispatcher) handleIssueEvent(e *gitlab.IssueEvent, l *logrus.Entry) {
	defer d.wg.Done()

	l = l.WithFields(logrus.Fields{
		logFieldURL:    e.ObjectAttributes.URL,
		logFieldAction: e.ObjectAttributes.Action,
	})

	if err := d.h.issueHandlers(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) handleMergeEvent(e *gitlab.MergeEvent, l *logrus.Entry) {
	defer d.wg.Done()

	l = l.WithFields(logrus.Fields{
		logFieldURL:    e.ObjectAttributes.URL,
		logFieldAction: e.ObjectAttributes.Action,
	})

	if err := d.h.mergeEventHandler(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) handlePushEvent(e *gitlab.PushEvent, l *logrus.Entry) {
	defer d.wg.Done()
	l = l.WithFields(logrus.Fields{
		logFieldOrg:  strings.Split(e.Project.PathWithNamespace, "/")[0],
		logFieldRepo: e.Repository.Name,
		"ref":        e.Ref,
		"head":       e.After,
	})

	if err := d.h.pushEventHandler(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) handleIssueCommentEvent(e *gitlab.IssueCommentEvent, l *logrus.Entry) {
	defer d.wg.Done()

	l = l.WithFields(logrus.Fields{
		logFieldURL:    e.Issue.URL,
		logFieldAction: e.Issue.State,
	})

	if err := d.h.issueCommentHandler(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) handleMergeCommentEvent(e *gitlab.MergeCommentEvent, l *logrus.Entry) {
	defer d.wg.Done()

	org, repo := gitlabclient.GetOrgRepo(e.Project.PathWithNamespace)
	l = l.WithFields(logrus.Fields{
		logFieldOrg:  org,
		logFieldRepo: repo,
		"review":     e.MergeRequest.LastCommit.ID,
		"reviewer":   e.MergeRequest.LastCommit.Author.Name,
		"url":        e.MergeRequest.LastCommit.URL,
	})

	if err := d.h.mergeCommentEventHandler(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) handleCommitCommentEvent(e *gitlab.CommitCommentEvent, l *logrus.Entry) {
	defer d.wg.Done()

	org, repo := gitlabclient.GetOrgRepo(e.Project.PathWithNamespace)
	l = l.WithFields(logrus.Fields{
		logFieldOrg:  org,
		logFieldRepo: repo,
		"commit":     e.Commit.ID,
		"reviewer":   e.User.Username,
		"url":        e.Commit.URL,
	})

	if err := d.h.commitCommentEventHandler(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eventType, eventGUID, payload, ok := parseRequest(w, r)
	if !ok {
		return
	}

	l := logrus.WithFields(
		logrus.Fields{
			"event-type": eventType,
			"event_id":   eventGUID,
		},
	)

	if err := d.Dispatch(eventType, payload, l); err != nil {
		l.WithError(err).Error()
	}
}

func parseRequest(w http.ResponseWriter, r *http.Request) (eventType string, uuid string, payload []byte, ok bool) {
	defer r.Body.Close()

	resp := func(code int, msg string) {
		http.Error(w, msg, code)
	}

	if r.Header.Get("User-Agent") != "Robot-Gitlab-Access" {
		resp(http.StatusBadRequest, "400 Bad Request: unknown User-Agent Header")
		return
	}

	if eventType = r.Header.Get("X-Gitlab-Event"); eventType == "" {
		resp(http.StatusBadRequest, "400 Bad Request: Missing X-Gitlab-Event Header")
		return
	}

	if uuid = r.Header.Get("X-Gitlab-Event-UUID"); uuid == "" {
		resp(http.StatusBadRequest, "400 Bad Request: Missing X-Gitlab-Event-UUID Header")
		return
	}

	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp(http.StatusInternalServerError, "500 Internal Server Error: Failed to read request body")
		return
	}

	payload = v
	ok = true

	return
}
