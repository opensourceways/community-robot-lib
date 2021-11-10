package giteeclient

import sdk "gitee.com/openeuler/go-gitee/gitee"

// NewIssueEventWrapper creates a wrapper of issue event.
func NewIssueEventWrapper(e *sdk.IssueEvent) IssueEventWrapper {
	return IssueEventWrapper{IssueEvent: e}
}

// IssueEventWrapper is a wrapper of the issue event to
// provide methods to obtain information about issue safely.
type IssueEventWrapper struct {
	*sdk.IssueEvent
}

// GetOrgRepo returns the action
func (i IssueEventWrapper) GetAction() string {
	if i.Action != nil {
		return *i.Action
	}
	return ""
}

// GetOrgRepo returns the org and repo
func (i IssueEventWrapper) GetOrgRep() (string, string) {
	return getOrgRepo(i.Repository)
}

// GetIssueAuthor returns the author of the issue
func (i IssueEventWrapper) GetIssueAuthor() string {
	if i.Issue != nil && i.Issue.User != nil {
		return i.Issue.User.Login
	}
	return ""
}

// GetIssueNumber returns the number of the issue
func (i IssueEventWrapper) GetIssueNumber() string {
	if i.Issue != nil {
		return i.Issue.Number
	}
	return ""
}
