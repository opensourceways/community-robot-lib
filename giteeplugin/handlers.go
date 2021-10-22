package giteeplugin

import (
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/community-robot-lib/config"
)

// IssueHandler defines the function contract for a gitee.IssueEvent handler.
type IssueHandler func(e *gitee.IssueEvent, cfg config.PluginConfig, log *logrus.Entry) error

// PullRequestHandler defines the function contract for a gitee.PullRequestEvent handler.
type PullRequestHandler func(e *gitee.PullRequestEvent, cfg config.PluginConfig, log *logrus.Entry) error

// PushEventHandler defines the function contract for a gitee.PushEvent handler.
type PushEventHandler func(e *gitee.PushEvent, cfg config.PluginConfig, log *logrus.Entry) error

// NoteEventHandler defines the function contract for a gitee.NoteEvent handler.
type NoteEventHandler func(e *gitee.NoteEvent, cfg config.PluginConfig, log *logrus.Entry) error

type handlers struct {
	issueHandlers      IssueHandler
	pullRequestHandler PullRequestHandler
	pushEventHandler   PushEventHandler
	noteEventHandler   NoteEventHandler
}

// RegisterIssueHandler registers a plugin's gitee.IssueEvent handler.
func (h *handlers) RegisterIssueHandler(fn IssueHandler) {
	h.issueHandlers = fn
}

// RegisterPullRequestHandler registers a plugin's gitee.PullRequestEvent handler.
func (h *handlers) RegisterPullRequestHandler(fn PullRequestHandler) {
	h.pullRequestHandler = fn
}

// RegisterPushEventHandler registers a plugin's gitee.PushEvent handler.
func (h *handlers) RegisterPushEventHandler(fn PushEventHandler) {
	h.pushEventHandler = fn
}

// RegisterNoteEventHandler registers a plugin's gitee.NoteEvent handler.
func (h *handlers) RegisterNoteEventHandler(fn NoteEventHandler) {
	h.noteEventHandler = fn
}
