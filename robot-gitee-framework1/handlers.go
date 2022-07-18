package framework

import (
	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/sirupsen/logrus"
)

const (
	logFieldOrg    = "org"
	logFieldRepo   = "repo"
	logFieldURL    = "url"
	logFieldAction = "action"
)

// IssueEventHandler defines handler for sdk.IssueEvent.
type IssueEventHandler interface {
	HandleIssueEvent(e *sdk.IssueEvent, log *logrus.Entry) error
}

// PullRequestHandler defines handler for sdk.PullRequestEvent.
type PullRequestEventHandler interface {
	HandlePullRequestEvent(e *sdk.PullRequestEvent, log *logrus.Entry) error
}

// PushEventHandler defines handler for sdk.PushEvent.
type PushEventHandler interface {
	HandlePushEvent(e *sdk.PushEvent, log *logrus.Entry) error
}

// NoteEventHandler defines handler for sdk.NoteEvent.
type NoteEventHandler interface {
	HandleNoteEvent(e *sdk.NoteEvent, log *logrus.Entry) error
}

type handlers struct {
	pullRequestEventHandler PullRequestEventHandler
	issueEventHandler       IssueEventHandler
	pushEventHandler        PushEventHandler
	noteEventHandler        NoteEventHandler
}

// registerHandler registers a robot's each handlers.
func (h *handlers) registerHandler(robot interface{}) {
	if v, ok := robot.(IssueEventHandler); ok {
		h.issueEventHandler = v
	}

	if v, ok := robot.(PullRequestEventHandler); ok {
		h.pullRequestEventHandler = v
	}

	if v, ok := robot.(PushEventHandler); ok {
		h.pushEventHandler = v
	}

	if v, ok := robot.(NoteEventHandler); ok {
		h.noteEventHandler = v
	}
}

func (h *handlers) getHandler() (r map[string]func([]byte, *logrus.Entry)) {
	r = make(map[string]func([]byte, *logrus.Entry))

	if h.issueEventHandler != nil {
		r[sdk.EventTypeIssue] = h.handleIssueEvent
	}

	if h.noteEventHandler != nil {
		r[sdk.EventTypeNote] = h.handleNoteEvent
	}

	if h.pushEventHandler != nil {
		r[sdk.EventTypePush] = h.handlePushEvent
	}

	if h.pullRequestEventHandler != nil {
		r[sdk.EventTypePR] = h.handlePullRequestEvent
	}

	return
}

func (h *handlers) handlePullRequestEvent(payload []byte, l *logrus.Entry) {
	e, err := sdk.ConvertToPREvent(payload)
	if err != nil {
		l.Errorf("convert to PREvent, err: ", err.Error())

		return
	}

	l = l.WithFields(logrus.Fields{
		logFieldURL:    e.PullRequest.HtmlUrl,
		logFieldAction: e.GetActionDesc(),
	})

	if err := h.pullRequestEventHandler.HandlePullRequestEvent(&e, l); err != nil {
		l.Error(err.Error())
	} else {
		l.Info()
	}
}

func (h *handlers) handleIssueEvent(payload []byte, l *logrus.Entry) {
	e, err := sdk.ConvertToIssueEvent(payload)
	if err != nil {
		l.Errorf("convert to IssueEvent, err: ", err.Error())

		return
	}

	l = l.WithFields(logrus.Fields{
		logFieldURL:    e.Issue.HtmlUrl,
		logFieldAction: *e.Action,
	})

	if err := h.issueEventHandler.HandleIssueEvent(&e, l); err != nil {
		l.Error(err.Error())
	} else {
		l.Info()
	}
}

func (h *handlers) handlePushEvent(payload []byte, l *logrus.Entry) {
	e, err := sdk.ConvertToPushEvent(payload)
	if err != nil {
		l.Errorf("convert to PushEvent, err: ", err.Error())

		return
	}

	org, repo := e.GetOrgRepo()

	l = l.WithFields(logrus.Fields{
		logFieldOrg:  org,
		logFieldRepo: repo,
		"ref":        e.Ref,
		"head":       e.After,
	})

	if err := h.pushEventHandler.HandlePushEvent(&e, l); err != nil {
		l.Error(err.Error())
	} else {
		l.Info()
	}
}

func (h *handlers) handleNoteEvent(payload []byte, l *logrus.Entry) {
	e, err := sdk.ConvertToNoteEvent(payload)
	if err != nil {
		l.Errorf("convert to NoteEvent, err: ", err.Error())

		return
	}

	l = l.WithFields(logrus.Fields{
		"commenter":    e.Comment.User.Login,
		logFieldURL:    e.Comment.HtmlUrl,
		logFieldAction: *e.Action,
	})

	if err := h.noteEventHandler.HandleNoteEvent(&e, l); err != nil {
		l.Error(err.Error())
	} else {
		l.Info()
	}
}
