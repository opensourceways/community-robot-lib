package giteeclient

import (
	sdk "gitee.com/openeuler/go-gitee/gitee"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	//StatusOpen gitee issue or pr status is open
	StatusOpen = "open"
	//StatusClosed gitee issue or pr status is closed
	StatusClosed = "closed"
)

//NoteEventWrapper a wrapper for the event of the comment to
//provide common methods for obtaining comment related information
type NoteEventWrapper struct {
	*sdk.NoteEvent
}

//IsCreatingCommentEvent Determine whether an note event is create a comment
func (ne NoteEventWrapper) IsCreatingCommentEvent() bool {
	return *(ne.Action) == "comment"
}

//GetCommenter Return to the author of the comment
func (ne NoteEventWrapper) GetCommenter() string {
	return ne.Comment.User.Login
}

//GetComment Return to the content of the comment
func (ne NoteEventWrapper) GetComment() string {
	return ne.Comment.Body
}

//GetOrgRepo Return to the org and repo
func (ne NoteEventWrapper) GetOrgRep() (string, string) {
	return getOrgRepo(ne.Repository)
}

//IsPullRequest Determine whether it is a PullRequest
func (ne NoteEventWrapper) IsPullRequest() bool {
	return *(ne.NoteableType) == "PullRequest"
}

//IsIssue Determine whether it is a issue
func (ne NoteEventWrapper) IsIssue() bool {
	return *(ne.NoteableType) == "Issue"
}

//IssueNoteEvent a wrapper for the event of the comment issue
//to provide methods for obtaining issue related information
type IssueNoteEvent struct {
	NoteEventWrapper
}

//IsIssueClosed whether the status of issue is close
func (ne IssueNoteEvent) IsIssueClosed() bool {
	return ne.Issue.State == StatusClosed
}

//IsIssueOpen whether the status of issue is open
func (ne IssueNoteEvent) IsIssueOpen() bool {
	return ne.Issue.State == StatusOpen
}

//GetIssueAuthor Return to the author of the issue
func (ne IssueNoteEvent) GetIssueAuthor() string {
	return ne.Issue.User.Login
}

//GetIssueNumber Return to the number of the issue
func (ne IssueNoteEvent) GetIssueNumber() string {
	return ne.Issue.Number
}

//GetIssueLabels returns the labels of the issue
func (ne IssueNoteEvent) GetIssueLabels() sets.String {
	return getLabelFromEvent(ne.Issue.Labels)
}

//PRNoteEvent a wrapper for the event of the comment pullrequest
//to provide methods for obtaining pullrequest related information
type PRNoteEvent struct {
	NoteEventWrapper
}

//GetPRNumber Return to the number of the PR
func (ne PRNoteEvent) GetPRNumber() int32 {
	return ne.PullRequest.Number
}

//GetPRAuthor Return to the author of the PR
func (ne PRNoteEvent) GetPRAuthor() string {
	return ne.PullRequest.User.Login
}

//IsPROpen whether the status of PR is open
func (ne PRNoteEvent) IsPROpen() bool {
	return ne.PullRequest.State == StatusOpen
}

// GetPRInfo gen PRInfo instance
func (ne PRNoteEvent) GetPRInfo() PRInfo {
	org, repo := ne.GetOrgRep()
	pr := ne.PullRequest
	return PRInfo{
		Org:     org,
		Repo:    repo,
		BaseRef: pr.Base.Ref,
		HeadSHA: pr.Head.Sha,
		Author:  ne.GetPRAuthor(),
		Number:  ne.GetPRNumber(),
		Labels:  ne.GetPRLabels(),
	}
}

// GetPRLabels returns the labels of the pr
func (ne PRNoteEvent) GetPRLabels() sets.String {
	return getLabelFromEvent(ne.PullRequest.Labels)
}

//NewNoteEventWrapper create a wrapper for comment events
func NewNoteEventWrapper(e *sdk.NoteEvent) NoteEventWrapper {
	return NoteEventWrapper{NoteEvent: e}
}

//NewIssueNoteEvent create a wrapper for the issue's comment event
func NewIssueNoteEvent(e *sdk.NoteEvent) IssueNoteEvent {
	return IssueNoteEvent{
		NoteEventWrapper: NoteEventWrapper{NoteEvent: e},
	}
}

//NewPRNoteEvent create a wrapper for the pr's comment event
func NewPRNoteEvent(e *sdk.NoteEvent) PRNoteEvent {
	return PRNoteEvent{
		NoteEventWrapper: NoteEventWrapper{NoteEvent: e},
	}
}
