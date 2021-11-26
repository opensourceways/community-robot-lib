package giteeplugin

import (
	"sync"

	sdk "gitee.com/openeuler/go-gitee/gitee"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/community-robot-lib/config"
	"github.com/opensourceways/community-robot-lib/giteeclient"
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
	switch eventType {
	case giteeclient.EventTypeNote:
		if d.h.noteEventHandler == nil {
			return nil
		}

		e, err := giteeclient.ConvertToNoteEvent(payload)
		if err != nil {
			return err
		}

		d.wg.Add(1)
		go d.handleNoteEvent(&e, l)

	case giteeclient.EventTypeIssue:
		if d.h.issueHandlers == nil {
			return nil
		}

		e, err := giteeclient.ConvertToIssueEvent(payload)
		if err != nil {
			return err
		}

		d.wg.Add(1)
		go d.handleIssueEvent(&e, l)

	case giteeclient.EventTypePR:
		if d.h.pullRequestHandler == nil {
			return nil
		}

		e, err := giteeclient.ConvertToPREvent(payload)
		if err != nil {
			return err
		}

		d.wg.Add(1)
		go d.handlePullRequestEvent(&e, l)

	case giteeclient.EventTypePush:
		if d.h.pushEventHandler == nil {
			return nil
		}

		e, err := giteeclient.ConvertToPushEvent(payload)
		if err != nil {
			return err
		}

		d.wg.Add(1)
		go d.handlePushEvent(&e, l)

	default:
		l.Debug("Ignoring unknown event type")
	}
	return nil
}

func (d *dispatcher) getConfig() config.PluginConfig {
	_, c := d.agent.GetConfig()
	return c
}

func (d *dispatcher) handlePullRequestEvent(e *sdk.PullRequestEvent, l *logrus.Entry) {
	defer d.wg.Done()

	l = l.WithFields(logrus.Fields{
		logFieldURL:    e.PullRequest.HtmlUrl,
		logFieldAction: e.GetActionDesc(),
	})

	if err := d.h.pullRequestHandler(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) handleIssueEvent(e *sdk.IssueEvent, l *logrus.Entry) {
	defer d.wg.Done()

	l = l.WithFields(logrus.Fields{
		logFieldURL:    e.Issue.HtmlUrl,
		logFieldAction: *e.Action,
	})

	if err := d.h.issueHandlers(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) handlePushEvent(e *sdk.PushEvent, l *logrus.Entry) {
	defer d.wg.Done()

	org, repo := giteeclient.GetOwnerAndRepoByPushEvent(e)

	l = l.WithFields(logrus.Fields{
		logFieldOrg:  org,
		logFieldRepo: repo,
		"ref":        e.Ref,
		"head":       e.After,
	})

	if err := d.h.pushEventHandler(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (d *dispatcher) handleNoteEvent(e *sdk.NoteEvent, l *logrus.Entry) {
	defer d.wg.Done()

	l = l.WithFields(logrus.Fields{
		"commenter":    e.Comment.User.Login,
		logFieldURL:    e.Comment.HtmlUrl,
		logFieldAction: *e.Action,
	})

	if err := d.h.noteEventHandler(e, d.getConfig(), l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}
