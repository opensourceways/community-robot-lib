package giteeclient

import sdk "gitee.com/openeuler/go-gitee/gitee"

// Client interface for Gitee API
type Client interface {
	CreatePullRequest(org, repo, title, body, head, base string, canModify bool) (sdk.PullRequest, error)
	GetPullRequests(org, repo string, opts ListPullRequestOpt) ([]sdk.PullRequest, error)
	UpdatePullRequest(org, repo string, number int32, param sdk.PullRequestUpdateParam) (sdk.PullRequest, error)

	ListCollaborators(org, repo string) ([]sdk.ProjectMember, error)
	IsCollaborator(owner, repo, login string) (bool, error)
	IsMember(org, login string) (bool, error)
	RemoveRepoMember(org, repo, login string) error
	AddRepoMember(org, repo, login, permission string) error

	GetRef(org, repo, ref string) (string, error)
	GetPullRequestChanges(org, repo string, number int32) ([]sdk.PullRequestFiles, error)
	GetPRLabels(org, repo string, number int32) ([]sdk.Label, error)
	ListPRComments(org, repo string, number int32) ([]sdk.PullRequestComments, error)
	ListPrIssues(org, repo string, number int32) ([]sdk.Issue, error)
	DeletePRComment(org, repo string, ID int32) error
	CreatePRComment(org, repo string, number int32, comment string) error
	UpdatePRComment(org, repo string, commentID int32, comment string) error
	AddPRLabel(org, repo string, number int32, label string) error
	AddMultiPRLabel(org, repo string, number int32, label []string) error
	RemovePRLabel(org, repo string, number int32, label string) error
	RemovePRLabels(org, repo string, number int32, labels []string) error
	ReplacePRAllLabels(owner, repo string, number int32, labels []string) error
	ListPROperationLogs(org, repo string, number int32) ([]sdk.OperateLog, error)

	ClosePR(org, repo string, number int32) error
	AssignPR(owner, repo string, number int32, logins []string) error
	UnassignPR(owner, repo string, number int32, logins []string) error
	GetPRCommits(org, repo string, number int32) ([]sdk.PullRequestCommits, error)

	GetGiteePullRequest(org, repo string, number int32) (sdk.PullRequest, error)
	GetPRCommit(org, repo, SHA string) (sdk.RepoCommit, error)
	MergePR(owner, repo string, number int32, opt sdk.PullRequestMergePutParam) error

	GetRepos(org string) ([]sdk.Project, error)
	CreateRepo(org string, repo sdk.RepositoryPostParam) error
	UpdateRepo(org, repo string, info sdk.RepoPatchParam) error
	GetRepo(org, repo string) (sdk.Project, error)
	GetGiteeRepo(org, repo string) (sdk.Project, error)

	SetRepoReviewer(org, repo string, reviewer sdk.SetRepoReviewer) error
	CreateRepoLabel(org, repo, label, color string) error
	GetRepoLabels(owner, repo string) ([]sdk.Label, error)

	AssignGiteeIssue(org, repo string, number string, login string) error
	UnassignGiteeIssue(org, repo string, number string, login string) error
	CreateIssueComment(org, repo string, number string, comment string) error
	UpdateIssueComment(org, repo string, commentID int32, comment string) error
	ListIssueComments(org, repo, number string) ([]sdk.Note, error)
	GetIssueLabels(org, repo, number string) ([]sdk.Label, error)
	RemoveIssueLabel(org, repo, number, label string) error
	RemoveIssueLabels(org, repo, number string, label []string) error
	AddIssueLabel(org, repo, number, label string) error
	AddMultiIssueLabel(org, repo, number string, label []string) error
	CloseIssue(owner, repo string, number string) error
	ReopenIssue(owner, repo string, number string) error
	UpdateIssue(owner, number string, param sdk.IssueUpdateParam) (sdk.Issue, error)
	GetIssue(org, repo, number string) (sdk.Issue, error)

	CreateBranch(org, repo, branch, parentBranch string) error
	GetRepoAllBranch(org, repo string) ([]sdk.Branch, error)
	SetProtectionBranch(org, repo, branch string) error
	CancelProtectionBranch(org, repo, branch string) error

	CreateFile(org, repo, branch, path, content, commitMsg string) (sdk.CommitContent, error)
	GetPathContent(org, repo, path, ref string) (sdk.Content, error)
	GetDirectoryTree(org, repo, sha string, recursive int32) (sdk.Tree, error)

	GetBot() (sdk.User, error)
	GetUserPermissionsOfRepo(org, repo, login string) (sdk.ProjectMemberPermission, error)
}

type ListPullRequestOpt struct {
	State           string
	Head            string
	Base            string
	Sort            string
	Direction       string
	MilestoneNumber int32
	Labels          []string
}
