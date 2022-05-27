package framework

import (
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"

	"github.com/opensourceways/community-robot-lib/config"
)

// IssueHandler defines the function contract for a gitlab.IssuesEvent handler.
type IssueHandler func(e *gitlab.IssueEvent, cfg config.Config, log *logrus.Entry) error

// IssueCommentHandler defines the function contract for a gitlab.IssueCommentEvent handler.
type IssueCommentHandler func(e *gitlab.IssueCommentEvent, cfg config.Config, log *logrus.Entry) error

// MergeEventHandler defines the function contract for a gitlab.MergeEvent handler.
type MergeEventHandler func(e *gitlab.MergeEvent, cfg config.Config, log *logrus.Entry) error

// PushEventHandler defines the function contract for a gitlab.PushEvent handler.
type PushEventHandler func(e *gitlab.PushEvent, cfg config.Config, log *logrus.Entry) error

// MergeCommentEventEventHandler defines the function contract for a gitlab.MergeCommentEventEventHandler handler.
type MergeCommentEventEventHandler func(e *gitlab.MergeCommentEvent, cfg config.Config, log *logrus.Entry) error

// CommitCommentEventHandler defines the function contract for a gitlab.CommitCommentEvent handler.
type CommitCommentEventHandler func(e *gitlab.CommitCommentEvent, cfg config.Config, log *logrus.Entry) error

type handlers struct {
	issueHandlers             IssueHandler
	mergeEventHandler         MergeEventHandler
	pushEventHandler          PushEventHandler
	issueCommentHandler       IssueCommentHandler
	mergeCommentEventHandler  MergeCommentEventEventHandler
	commitCommentEventHandler CommitCommentEventHandler
}

// RegisterIssueHandler registers a plugin's gitlab.IssueEvent handler.
func (h *handlers) RegisterIssueHandler(fn IssueHandler) {
	h.issueHandlers = fn
}

// RegisterMergeEventHandler registers a plugin's gitlab.PullRequestEvent handler.
func (h *handlers) RegisterMergeEventHandler(fn MergeEventHandler) {
	h.mergeEventHandler = fn
}

// RegisterPushEventHandler registers a plugin's gitlab.PushEvent handler.
func (h *handlers) RegisterPushEventHandler(fn PushEventHandler) {
	h.pushEventHandler = fn
}

// RegisterIssueCommentHandler registers a plugin's gitlab.IssueCommentEvent handler.
func (h *handlers) RegisterIssueCommentHandler(fn IssueCommentHandler) {
	h.issueCommentHandler = fn
}

// RegisterMergeCommentEventHandler registers a plugin's gitlab.MergeCommentEvent handler.
func (h *handlers) RegisterMergeCommentEventHandler(fn MergeCommentEventEventHandler) {
	h.mergeCommentEventHandler = fn
}

// RegisterCommitCommentEventHandler registers a plugin's gitlab.CommitCommentEvent handler.
func (h *handlers) RegisterCommitCommentEventHandler(fn CommitCommentEventHandler) {
	h.commitCommentEventHandler = fn
}
