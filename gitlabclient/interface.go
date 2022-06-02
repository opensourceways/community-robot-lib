package gitlabclient

import "github.com/xanzy/go-gitlab"

type Client interface {
	UpdateMergeRequest(projectID interface{}, mrID int, options gitlab.UpdateMergeRequestOptions) (gitlab.MergeRequest, error)
	GetMergeRequest(projectID interface{}, mrID int) (gitlab.MergeRequest, error)

	ListCollaborators(projectID interface{}) ([]*gitlab.ProjectMember, error)
	IsCollaborator(projectID interface{}, loginID int) (bool, error)
	AddProjectMember(projectID interface{}, loginID interface{}) error
	RemoveProjectMember(projectID interface{}, loginID int) error
	IsMember(groupID interface{}, userID int) (bool, error)

	GetMergeRequestChanges(projectID interface{}, mrID int) ([]string, error)
	GetMergeRequestLabels(projectID interface{}, mrID int) (gitlab.Labels, error)
	ListMergeRequestComments(projectID interface{}, mrID int) ([]*gitlab.Note, error)
	ListIssues(projectID interface{}) ([]*gitlab.Issue, error)
	ListIssueRelatedMergeRequest(projectID interface{}, issueID int) ([]*gitlab.MergeRequest, error)
	UpdateMergeRequestComment(projectID interface{}, mrID, noteID int, comment string) error
	CreateMergeRequestComment(projectID interface{}, mrID int, comment string) error
	DeleteMergeRequestComment(projectID interface{}, mrID int, noteID int) error
	AddMergeRequestLabel(projectID interface{}, mrID int, labels gitlab.Labels) error
	RemoveMergeRequestLabel(projectID interface{}, mrID int, labels gitlab.Labels) error
}
