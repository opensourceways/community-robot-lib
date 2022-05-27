package gitlabclient

import (
	"github.com/xanzy/go-gitlab"
)

const (
	ActionOpened  = "opened"
	ActionCreated = "created"
	ActionReopen  = "reopened"
	ActionClosed  = "closed"
)

// GetOrgRepo return the owner and name of the repository
func GetOrgRepo(repo *gitlab.Repository) (string, string) {
	return repo.Namespace, repo.Name
}
