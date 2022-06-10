package gitlabclient

import (
	"encoding/json"
	"github.com/xanzy/go-gitlab"
	"strings"
)

const (
	ActionOpened  = "opened"
	ActionCreated = "created"
	ActionReopen  = "reopened"
	ActionClosed  = "closed"
)

type ErrorForbidden struct {
	err string
}

func (e ErrorForbidden) Error() string {
	return e.err
}

// GetOrgRepo return the owner and name of the project
func GetOrgRepo(orgPath string) (string, string) {
	org, repo := strings.Split(orgPath, "/")[0], strings.Split(orgPath, "/")[1]
	return org, repo
}

func ConvertToMergeEvent(payload []byte) (e *gitlab.MergeEvent, err error) {
	if err = json.Unmarshal(payload, &e); err != nil {
		return
	}
	return
}

func ConvertToIssueEvent(payload []byte) (e *gitlab.IssueEvent, err error) {
	if err = json.Unmarshal(payload, &e); err != nil {
		return
	}
	return
}

func ConvertToMergeCommentEvent(payload []byte) (e *gitlab.MergeCommentEvent, err error) {
	if err = json.Unmarshal(payload, &e); err != nil {
		return
	}
	return
}

func ConvertToIssueCommentEvent(payload []byte) (e *gitlab.IssueCommentEvent, err error) {
	if err = json.Unmarshal(payload, &e); err != nil {
		return
	}
	return
}

func ConvertToPushEvent(payload []byte) (e *gitlab.PushEvent, err error) {
	if err = json.Unmarshal(payload, &e); err != nil {
		return
	}
	return
}

func ConvertToCommitCommentEvent(payload []byte) (e *gitlab.CommitCommentEvent, err error) {
	if err = json.Unmarshal(payload, &e); err != nil {
		return
	}
	return
}
