package githubclient

import "github.com/google/go-github/v36/github"

const (
	ActionOpened  = "opened"
	ActionCreated = "created"
	ActionReopen  = "reopened"
	ActionClosed  = "closed"
)

// GetOrgRepo return the owner and name of the repository
func GetOrgRepo(repo *github.Repository) (string, string) {
	return repo.GetOwner().GetLogin(), repo.GetName()
}

// IsIssueOpened judge is issue create event
func IsIssueOpened(action string) bool {
	return action == ActionOpened
}

// IsIssueCommentCreated judge is issue comment create event
func IsIssueCommentCreated(action string) bool {
	return action == ActionCreated
}
