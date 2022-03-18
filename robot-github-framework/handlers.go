package framework

import (
	"github.com/google/go-github/v36/github"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/community-robot-lib/config"
)

// IssueHandler defines the function contract for a github.IssuesEvent handler.
type IssueHandler func(e *github.IssuesEvent, cfg config.Config, log *logrus.Entry) error

// IssueCommentHandler defines the function contract for a github.IssueCommentEvent handler.
type IssueCommentHandler func(e *github.IssueCommentEvent, cfg config.Config, log *logrus.Entry) error

// PullRequestHandler defines the function contract for a github.PullRequestEvent handler.
type PullRequestHandler func(e *github.PullRequestEvent, cfg config.Config, log *logrus.Entry) error

// StatusEventHandler defines the function contract for a github.StatusEvent handler.
type StatusEventHandler func(e *github.StatusEvent, cfg config.Config, log *logrus.Entry) error

// PushEventHandler defines the function contract for a github.PushEvent handler.
type PushEventHandler func(e *github.PushEvent, cfg config.Config, log *logrus.Entry) error

// ReviewEventHandler defines the function contract for a github.PullRequestReviewEvent handler.
type ReviewEventHandler func(e *github.PullRequestReviewEvent, cfg config.Config, log *logrus.Entry) error

// ReviewCommentEventHandler defines the function contract for a github.PullRequestReviewCommentEvent handler.
type ReviewCommentEventHandler func(e *github.PullRequestReviewCommentEvent, cfg config.Config, log *logrus.Entry) error

// CommitCommentEventHandler defines the function contract for a github.CommitCommentEvent handler.
type CommitCommentEventHandler func(e *github.CommitCommentEvent, cfg config.Config, log *logrus.Entry) error

type handlers struct {
	issueHandlers             IssueHandler
	pullRequestHandler        PullRequestHandler
	pushEventHandler          PushEventHandler
	issueCommentHandler       IssueCommentHandler
	statusEventHandler        StatusEventHandler
	reviewEventHandler        ReviewEventHandler
	reviewCommentEventHandler ReviewCommentEventHandler
	commitCommentEventHandler CommitCommentEventHandler
}

// RegisterIssueHandler registers a plugin's github.IssueEvent handler.
func (h *handlers) RegisterIssueHandler(fn IssueHandler) {
	h.issueHandlers = fn
}

// RegisterPullRequestHandler registers a plugin's github.PullRequestEvent handler.
func (h *handlers) RegisterPullRequestHandler(fn PullRequestHandler) {
	h.pullRequestHandler = fn
}

// RegisterPushEventHandler registers a plugin's github.PushEvent handler.
func (h *handlers) RegisterPushEventHandler(fn PushEventHandler) {
	h.pushEventHandler = fn
}

// RegisterIssueCommentHandler registers a plugin's github.IssueCommentEvent handler.
func (h *handlers) RegisterIssueCommentHandler(fn IssueCommentHandler) {
	h.issueCommentHandler = fn
}

// RegisterStatusEventHandler registers a plugin's github.StatusEvent handler.
func (h *handlers) RegisterStatusEventHandler(fn StatusEventHandler) {
	h.statusEventHandler = fn
}

// RegisterReviewEventHandler registers a plugin's github.ReviewEvent handler.
func (h *handlers) RegisterReviewEventHandler(fn ReviewEventHandler) {
	h.reviewEventHandler = fn
}

// RegisterReviewCommentEventHandler registers a plugin's github.ReviewCommentEvent handler.
func (h *handlers) RegisterReviewCommentEventHandler(fn ReviewCommentEventHandler) {
	h.reviewCommentEventHandler = fn
}

// RegisterCommitCommentEventHandler registers a plugin's github.CommitCommentEvent handler.
func (h *handlers) RegisterCommitCommentEventHandler(fn CommitCommentEventHandler) {
	h.commitCommentEventHandler = fn
}
