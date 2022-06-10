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
	ReplaceMergeRequestAllLabels(projectID interface{}, mrID int, labels gitlab.Labels) error
	ReopenMergeRequest(projectID interface{}, mrID int, state string) error
	CloseMergeRequest(projectID interface{}, mrID int, state string) error
	AssignMergeRequest(projectID interface{}, mrID int, ids []int) error
	UnAssignMergeRequest(projectID interface{}, mrID int, ids []int) error
	GetMergeRequestCommits(projectID interface{}, mrID int) ([]*gitlab.Commit, error)
	GetSingleRepoCommit(projectID interface{}, sha string) (*gitlab.Commit, error)
	MergeMergeRequest(projectID interface{}, mrID int) error
	GetGroups() ([]*gitlab.Group, error)
	GetProjects(gid interface{}) ([]*gitlab.Project, error)
	GetProject(projectID interface{}) (*gitlab.Project, error)
	CreateProject(opts gitlab.CreateProjectOptions) error
	UpdateProject(projectID interface{}, opts gitlab.EditProjectOptions) error
	AddProjectLabel(projectID interface{}, label, color string) error
	UpdateProjectLabel(projectID interface{}, oldLabel, label, color string) error
	GetProjectLabels(projectID interface{}) ([]*gitlab.Label, error)
	AssignIssue(projectID interface{}, issueID int, assignees []int) error
	UpdateIssue(projectID interface{}, issueID int, opts gitlab.UpdateIssueOptions) error
	CreateIssue(projectID interface{}, opts gitlab.CreateIssueOptions) error
	UnAssignIssue(projectID interface{}, issueID int, assignees []int) error
	RemoveAssignIssue(projectID interface{}, issueID int) error
	CreateIssueComment(projectID interface{}, issueID int, comment string) error
	ListIssueComments(projectID interface{}, issueID int) ([]*gitlab.Note, error)
	UpdateIssueComment(projectID interface{}, issueID, noteID int, comment string) error
	RemoveIssueComment(projectID interface{}, issueID, noteID int) error
	GetIssueLabels(projectID interface{}, issueID int) ([]string, error)
	RemoveIssueLabels(projectID interface{}, issueID int, labels gitlab.Labels) error
	AddIssueLabels(projectID interface{}, issueID int, labels gitlab.Labels) error
	CloseIssue(projectID interface{}, issueID int) error
	ReopenIssue(projectID interface{}, issueID int) error
	GetSingleIssue(projectID interface{}, issueID int) (*gitlab.Issue, error)
	CreateBranch(projectID interface{}, branch, ref string) error
	GetProjectAllBranches(projectID interface{}) ([]*gitlab.Branch, error)
	SetProtectionBranch(projectID interface{}, branch string) error
	UnProtectBranch(projectID interface{}, branch string) error
	CreateFile(projectID interface{}, file string, opts gitlab.CreateFileOptions) error
	GetPathContent(projectID interface{}, file, filePath string) (*gitlab.File, error)
	GetDirectoryTree(projectID interface{}, opts gitlab.ListTreeOptions) ([]*gitlab.TreeNode, error)
	GetUserPermissionOfProject(projectID interface{}, userID int) (bool, error)
	GetMergeRequestLabelChanges(projectID interface{}, mrID int) ([]*gitlab.LabelEvent, error)
	CreateProjectLabel(pid interface{}, label, color string) error
}
