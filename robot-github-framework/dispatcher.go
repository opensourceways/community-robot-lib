package framework

import (
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/google/go-github/v36/github"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/community-robot-lib/config"
	"github.com/opensourceways/community-robot-lib/githubclient"
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
	hook, err := github.ParseWebHook(eventType, payload)
	if err != nil {
		return err
	}

	switch hook := hook.(type) {
	case *github.IssuesEvent:
		d.wg.Add(1)
		go d.handleIssueEvent(hook, l)
	case *github.PullRequestEvent:
		d.wg.Add(1)
		go d.handlePullRequestEvent(hook, l)
	case *github.PushEvent:
		d.wg.Add(1)
		go d.handlePushEvent(hook, l)
	case *github.IssueCommentEvent:
		d.wg.Add(1)
		go d.handleIssueCommentEvent(hook, l)
	case *github.PullRequestReviewEvent:
		d.wg.Add(1)
		go d.handleReviewEvent(hook, l)
	case *github.PullRequestReviewCommentEvent:
		d.wg.Add(1)
		go d.handleReviewCommentEvent(hook, l)
	case *github.StatusEvent:
		d.wg.Add(1)
		go d.handleStatusEvent(hook, l)
	case *github.CommitCommentEvent:
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

func (d *dispatcher) handleIssueEvent(e *github.IssuesEvent, l *logrus.Entry) {
	defer d.wg.Done()

	l = l.WithFields(logrus.Fields{
		logFieldURL:    e.GetIssue().GetHTMLURL(),
		logFieldAction: e.GetAction(),
	})

	if err := d.h.issueHandlers(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) handlePullRequestEvent(e *github.PullRequestEvent, l *logrus.Entry) {
	defer d.wg.Done()

	l = l.WithFields(logrus.Fields{
		logFieldURL:    e.GetPullRequest().GetHTMLURL(),
		logFieldAction: e.GetAction(),
	})

	if err := d.h.pullRequestHandler(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) handlePushEvent(e *github.PushEvent, l *logrus.Entry) {
	defer d.wg.Done()
	l = l.WithFields(logrus.Fields{
		logFieldOrg:  e.GetRepo().GetOwner().GetLogin(),
		logFieldRepo: e.GetRepo().GetName(),
		"ref":        e.GetRef(),
		"head":       e.GetAfter(),
	})

	if err := d.h.pushEventHandler(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) handleIssueCommentEvent(e *github.IssueCommentEvent, l *logrus.Entry) {
	defer d.wg.Done()

	l = l.WithFields(logrus.Fields{
		logFieldURL:    e.GetIssue().GetHTMLURL(),
		logFieldAction: e.GetAction(),
	})

	if err := d.h.issueCommentHandler(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) handleStatusEvent(e *github.StatusEvent, l *logrus.Entry) {
	defer d.wg.Done()

	org, repo := githubclient.GetOrgRepo(e.GetRepo())
	l = l.WithFields(logrus.Fields{
		logFieldOrg:  org,
		logFieldRepo: repo,
		"context":    e.GetContext(),
		"sha":        e.GetSHA(),
		"state":      e.GetState(),
		"id":         e.GetID(),
	})

	if err := d.h.statusEventHandler(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) handleReviewEvent(e *github.PullRequestReviewEvent, l *logrus.Entry) {
	defer d.wg.Done()

	org, repo := githubclient.GetOrgRepo(e.GetRepo())
	l = l.WithFields(logrus.Fields{
		logFieldOrg:  org,
		logFieldRepo: repo,
		"review":     e.GetReview().GetID(),
		"reviewer":   e.GetReview().GetUser().GetLogin(),
		"url":        e.GetReview().GetHTMLURL(),
	})

	if err := d.h.reviewEventHandler(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) handleReviewCommentEvent(e *github.PullRequestReviewCommentEvent, l *logrus.Entry) {
	defer d.wg.Done()

	org, repo := githubclient.GetOrgRepo(e.GetRepo())
	l = l.WithFields(logrus.Fields{
		logFieldOrg:  org,
		logFieldRepo: repo,
		"review":     e.GetComment().GetPullRequestReviewID(),
		"reviewer":   e.GetComment().GetUser().GetLogin(),
		"url":        e.GetComment().GetHTMLURL(),
	})

	if err := d.h.reviewCommentEventHandler(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) handleCommitCommentEvent(e *github.CommitCommentEvent, l *logrus.Entry) {
	defer d.wg.Done()

	org, repo := githubclient.GetOrgRepo(e.GetRepo())
	l = l.WithFields(logrus.Fields{
		logFieldOrg:  org,
		logFieldRepo: repo,
		"commit":     e.GetComment().GetCommitID(),
		"reviewer":   e.GetComment().GetUser().GetLogin(),
		"url":        e.GetComment().GetHTMLURL(),
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

	if r.Header.Get("User-Agent") != "Robot-Github-Access" {
		resp(http.StatusBadRequest, "400 Bad Request: unknown User-Agent Header")
		return
	}

	if eventType = r.Header.Get("X-GitHub-Event"); eventType == "" {
		resp(http.StatusBadRequest, "400 Bad Request: Missing X-GitHub-Event Header")
		return
	}

	if uuid = r.Header.Get("X-GitHub-Delivery"); uuid == "" {
		resp(http.StatusBadRequest, "400 Bad Request: Missing X-GitHub-Delivery Header")
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
