package giteeclient

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/antihax/optional"
	sdk "github.com/opensourceways/go-gitee/gitee"
	"golang.org/x/oauth2"
)

var _ Client = (*client)(nil)

type client struct {
	ac *sdk.APIClient
}

func NewClient(getToken func() []byte) Client {
	token := string(getToken())

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	conf := sdk.NewConfiguration()
	conf.HTTPClient = oauth2.NewClient(context.Background(), ts)

	c := sdk.NewAPIClient(conf)
	return &client{ac: c}
}

func (c *client) CreatePullRequest(org, repo, title, body, head, base string, canModify bool) (sdk.PullRequest, error) {
	opts := sdk.CreatePullRequestParam{
		Title:             title,
		Head:              head,
		Base:              base,
		Body:              body,
		PruneSourceBranch: true,
	}

	pr, _, err := c.ac.PullRequestsApi.PostV5ReposOwnerRepoPulls(
		context.Background(), org, repo, opts)

	return pr, formatErr(err, "create pull request")
}

func (c *client) GetPullRequests(org, repo string, opts ListPullRequestOpt) ([]sdk.PullRequest, error) {

	setStr := func(t *optional.String, v string) {
		if v != "" {
			*t = optional.NewString(v)
		}
	}

	opt := sdk.GetV5ReposOwnerRepoPullsOpts{}
	setStr(&opt.State, opts.State)
	setStr(&opt.Head, opts.Head)
	setStr(&opt.Base, opts.Base)
	setStr(&opt.Sort, opts.Sort)
	setStr(&opt.Direction, opts.Direction)
	if opts.MilestoneNumber > 0 {
		opt.MilestoneNumber = optional.NewInt32(opts.MilestoneNumber)
	}
	if opts.Labels != nil && len(opts.Labels) > 0 {
		opt.Labels = optional.NewString(strings.Join(opts.Labels, ","))
	}

	var r []sdk.PullRequest
	p := int32(1)
	for {
		opt.Page = optional.NewInt32(p)
		prs, _, err := c.ac.PullRequestsApi.GetV5ReposOwnerRepoPulls(context.Background(), org, repo, &opt)
		if err != nil {
			return nil, formatErr(err, "get pull requests")
		}

		if len(prs) == 0 {
			break
		}

		r = append(r, prs...)
		p++
	}

	return r, nil
}

func (c *client) UpdatePullRequest(org, repo string, number int32, param sdk.PullRequestUpdateParam) (sdk.PullRequest, error) {
	pr, _, err := c.ac.PullRequestsApi.PatchV5ReposOwnerRepoPullsNumber(context.Background(), org, repo, number, param)
	return pr, formatErr(err, "update pull request")
}

func (c *client) GetGiteePullRequest(org, repo string, number int32) (sdk.PullRequest, error) {
	pr, _, err := c.ac.PullRequestsApi.GetV5ReposOwnerRepoPullsNumber(
		context.Background(), org, repo, number, nil)
	return pr, formatErr(err, "get pull request")
}

func (c *client) GetBot() (sdk.User, error) {
	u, _, err := c.ac.UsersApi.GetV5User(context.Background(), nil)
	if err != nil {
		return u, formatErr(err, "fetch bot name")
	}
	return u, nil
}

func (c *client) ListCollaborators(org, repo string) ([]sdk.ProjectMember, error) {
	var r []sdk.ProjectMember

	opt := sdk.GetV5ReposOwnerRepoCollaboratorsOpts{}
	p := int32(1)
	for {
		opt.Page = optional.NewInt32(p)
		cs, _, err := c.ac.RepositoriesApi.GetV5ReposOwnerRepoCollaborators(context.Background(), org, repo, &opt)
		if err != nil {
			return nil, formatErr(err, "list collaborators")
		}
		if len(cs) == 0 {
			break
		}

		r = append(r, cs...)
		p++
	}
	return r, nil
}

func (c *client) GetRef(org, repo, ref string) (string, error) {
	branch := strings.TrimPrefix(ref, "heads/")
	b, _, err := c.ac.RepositoriesApi.GetV5ReposOwnerRepoBranchesBranch(context.Background(), org, repo, branch, nil)
	if err != nil {
		return "", formatErr(err, "get branch")
	}

	return b.Commit.Sha, nil
}

func (c *client) GetPullRequestChanges(org, repo string, number int32) ([]sdk.PullRequestFiles, error) {
	fs, _, err := c.ac.PullRequestsApi.GetV5ReposOwnerRepoPullsNumberFiles(
		context.Background(), org, repo, number, nil)
	if err != nil {
		return nil, formatErr(err, "list files of pr")
	}

	return fs, nil
}

func (c *client) GetPRLabels(org, repo string, number int32) ([]sdk.Label, error) {
	var r []sdk.Label

	p := int32(1)
	opt := sdk.GetV5ReposOwnerRepoPullsNumberLabelsOpts{}
	for {
		opt.Page = optional.NewInt32(p)
		ls, _, err := c.ac.PullRequestsApi.GetV5ReposOwnerRepoPullsNumberLabels(
			context.Background(), org, repo, number, &opt)
		if err != nil {
			return nil, formatErr(err, "list labels of pr")
		}

		if len(ls) == 0 {
			break
		}

		r = append(r, ls...)
		p++
	}

	return r, nil
}

func (c *client) ListPRComments(org, repo string, number int32) ([]sdk.PullRequestComments, error) {
	var r []sdk.PullRequestComments

	p := int32(1)
	opt := sdk.GetV5ReposOwnerRepoPullsNumberCommentsOpts{}
	for {
		opt.Page = optional.NewInt32(p)
		cs, _, err := c.ac.PullRequestsApi.GetV5ReposOwnerRepoPullsNumberComments(
			context.Background(), org, repo, number, &opt)
		if err != nil {
			return nil, formatErr(err, "list comments of pr")
		}

		if len(cs) == 0 {
			break
		}

		r = append(r, cs...)
		p++
	}

	return r, nil
}

func (c *client) ListPROperationLogs(org, repo string, number int32) ([]sdk.OperateLog, error) {
	v, _, err := c.ac.PullRequestsApi.GetV5ReposOwnerRepoPullsNumberOperateLogs(
		context.Background(), org, repo, number, nil)
	if err != nil {
		return nil, formatErr(err, "list operation logs of pr")
	}
	return v, nil
}

func (c *client) ListPrIssues(org, repo string, number int32) ([]sdk.Issue, error) {
	var issues []sdk.Issue
	p := int32(1)
	opt := sdk.GetV5ReposOwnerRepoPullsNumberIssuesOpts{}
	for {
		opt.Page = optional.NewInt32(p)
		iss, _, err := c.ac.PullRequestsApi.GetV5ReposOwnerRepoPullsNumberIssues(context.Background(), org, repo, number, &opt)
		if err != nil {
			return nil, formatErr(err, "list issues of pr")
		}
		if len(iss) == 0 {
			break
		}
		issues = append(issues, iss...)
		p++
	}
	return issues, nil
}

func (c *client) DeletePRComment(org, repo string, ID int32) error {
	_, err := c.ac.PullRequestsApi.DeleteV5ReposOwnerRepoPullsCommentsId(
		context.Background(), org, repo, ID, nil)
	return formatErr(err, "delete comment of pr")
}

func (c *client) CreatePRComment(org, repo string, number int32, comment string) error {
	opt := sdk.PullRequestCommentPostParam{Body: comment}
	_, _, err := c.ac.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberComments(
		context.Background(), org, repo, number, opt)
	return formatErr(err, "create comment of pr")
}

func (c *client) UpdatePRComment(org, repo string, commentID int32, comment string) error {
	opt := sdk.PullRequestCommentPatchParam{Body: comment}
	_, _, err := c.ac.PullRequestsApi.PatchV5ReposOwnerRepoPullsCommentsId(
		context.Background(), org, repo, commentID, opt)
	return formatErr(err, "update comment of pr")
}

func (c *client) AddPRLabel(org, repo string, number int32, label string) error {
	return c.AddMultiPRLabel(org, repo, number, []string{label})
}

func (c *client) AddMultiPRLabel(org, repo string, number int32, label []string) error {
	opt := sdk.PullRequestLabelPostParam{Body: label}
	_, _, err := c.ac.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberLabels(
		context.Background(), org, repo, number, opt)
	return formatErr(err, "add multi label for pr")
}

func (c *client) RemovePRLabel(org, repo string, number int32, label string) error {
	// gitee's bug, it can't deal with the label which includes '/'
	label = strings.Replace(label, "/", "%2F", -1)

	v, err := c.ac.PullRequestsApi.DeleteV5ReposOwnerRepoPullsLabel(
		context.Background(), org, repo, number, label, nil)

	if err == nil || (v != nil && v.StatusCode == 404) {
		return nil
	}
	return formatErr(err, "remove label of pr")
}

func (c *client) RemovePRLabels(org, repo string, number int32, labels []string) error {
	return c.RemovePRLabel(org, repo, number, strings.Join(labels, ","))
}

func (c *client) ClosePR(org, repo string, number int32) error {
	opt := sdk.PullRequestUpdateParam{State: sdk.StatusClosed}
	_, err := c.UpdatePullRequest(org, repo, number, opt)
	return formatErr(err, "close pr")
}

func (c *client) AssignPR(org, repo string, number int32, logins []string) error {
	opt := sdk.PullRequestAssigneePostParam{Assignees: strings.Join(logins, ",")}
	_, _, err := c.ac.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberAssignees(
		context.Background(), org, repo, number, opt)
	return formatErr(err, "assign reviewer to pr")
}

func (c *client) UnassignPR(org, repo string, number int32, logins []string) error {
	_, _, err := c.ac.PullRequestsApi.DeleteV5ReposOwnerRepoPullsNumberAssignees(
		context.Background(), org, repo, number, strings.Join(logins, ","), nil)
	return formatErr(err, "unassign reviewer from pr")
}

func (c *client) GetPRCommits(org, repo string, number int32) ([]sdk.PullRequestCommits, error) {
	commits, _, err := c.ac.PullRequestsApi.GetV5ReposOwnerRepoPullsNumberCommits(
		context.Background(), org, repo, number, nil)
	return commits, formatErr(err, "get pr commits")
}

func (c *client) AssignGiteeIssue(org, repo string, number string, login string) error {
	opt := sdk.IssueUpdateParam{Repo: repo, Assignee: login}
	_, v, err := c.ac.IssuesApi.PatchV5ReposOwnerIssuesNumber(
		context.Background(), org, number, opt)

	if err != nil {
		if v.StatusCode == 403 {
			return ErrorForbidden{err: formatErr(err, "assign assignee to issue").Error()}
		}
	}
	return formatErr(err, "assign assignee to issue")
}

func (c *client) UnassignGiteeIssue(org, repo string, number string, login string) error {
	return c.AssignGiteeIssue(org, repo, number, " ")
}

func (c *client) RemoveIssueAssignee(org, repo string, number string) error {
	return c.AssignGiteeIssue(org, repo, number, " ")
}

func (c *client) CreateIssueComment(org, repo string, number string, comment string) error {
	opt := sdk.IssueCommentPostParam{Body: comment}
	_, _, err := c.ac.IssuesApi.PostV5ReposOwnerRepoIssuesNumberComments(
		context.Background(), org, repo, number, opt)
	return formatErr(err, "create issue comment")
}

func (c *client) IsCollaborator(owner, repo, login string) (bool, error) {
	v, err := c.ac.RepositoriesApi.GetV5ReposOwnerRepoCollaboratorsUsername(
		context.Background(), owner, repo, login, nil)
	if err == nil {
		return true, nil
	}

	if v != nil && v.StatusCode == 404 {
		return false, nil
	}
	return false, formatErr(err, "get collaborator of pr")
}

func (c *client) IsMember(org, login string) (bool, error) {
	_, v, err := c.ac.OrganizationsApi.GetV5OrgsOrgMembershipsUsername(
		context.Background(), org, login, nil)
	if err == nil {
		return true, nil
	}

	if v != nil && v.StatusCode == 404 {
		return false, nil
	}
	return false, formatErr(err, "get member of org")
}

func (c *client) GetPRCommit(org, repo, SHA string) (sdk.RepoCommit, error) {
	v, _, err := c.ac.RepositoriesApi.GetV5ReposOwnerRepoCommitsSha(
		context.Background(), org, repo, SHA, nil)
	if err != nil {
		return v, formatErr(err, "get commit info")
	}

	return v, nil
}

func (c *client) MergePR(owner, repo string, number int32, opt sdk.PullRequestMergePutParam) error {
	_, err := c.ac.PullRequestsApi.PutV5ReposOwnerRepoPullsNumberMerge(
		context.Background(), owner, repo, number, opt)
	return formatErr(err, "merge pr")
}

func (c *client) GetGiteeRepo(org, repo string) (sdk.Project, error) {
	v, _, err := c.ac.RepositoriesApi.GetV5ReposOwnerRepo(context.Background(), org, repo, nil)
	return v, formatErr(err, "get repo")
}

func (c *client) GetRepo(org, repo string) (sdk.Project, error) {
	return c.GetGiteeRepo(org, repo)
}

func (c *client) GetRepos(org string) ([]sdk.Project, error) {
	opt := sdk.GetV5OrgsOrgReposOpts{}
	var r []sdk.Project
	p := int32(1)
	for {
		opt.Page = optional.NewInt32(p)
		ps, _, err := c.ac.RepositoriesApi.GetV5OrgsOrgRepos(context.Background(), org, &opt)
		if err != nil {
			return nil, formatErr(err, "list repos")
		}

		if len(ps) == 0 {
			break
		}
		r = append(r, ps...)
		p++
	}

	return r, nil
}

func (c *client) GetRepoLabels(owner, repo string) ([]sdk.Label, error) {
	labels, _, err := c.ac.LabelsApi.GetV5ReposOwnerRepoLabels(context.Background(), owner, repo, nil)
	return labels, formatErr(err, "get repo labels")
}

func (c *client) AddIssueLabel(org, repo, number, label string) error {
	opt := sdk.PullRequestLabelPostParam{Body: []string{label}}
	_, _, err := c.ac.LabelsApi.PostV5ReposOwnerRepoIssuesNumberLabels(
		context.Background(), org, repo, number, opt)
	return formatErr(err, "add issue label")
}

func (c *client) AddMultiIssueLabel(org, repo, number string, label []string) error {
	opt := sdk.PullRequestLabelPostParam{Body: label}
	_, _, err := c.ac.LabelsApi.PostV5ReposOwnerRepoIssuesNumberLabels(
		context.Background(), org, repo, number, opt)
	return formatErr(err, "add issue label")
}

func (c *client) RemoveIssueLabel(org, repo, number, label string) error {
	label = strings.Replace(label, "/", "%2F", -1)
	_, err := c.ac.LabelsApi.DeleteV5ReposOwnerRepoIssuesNumberLabelsName(
		context.Background(), org, repo, number, label, nil)
	return formatErr(err, "rm issue label")
}

func (c *client) RemoveIssueLabels(org, repo, number string, label []string) error {
	return c.RemoveIssueLabel(org, repo, number, strings.Join(label, ","))
}

func (c *client) ReplacePRAllLabels(owner, repo string, number int32, labels []string) error {
	opt := sdk.PullRequestLabelPostParam{Body: labels}
	_, _, err := c.ac.PullRequestsApi.PutV5ReposOwnerRepoPullsNumberLabels(context.Background(), owner, repo, number, opt)
	return formatErr(err, "replace pr labels")
}

func (c *client) CloseIssue(owner, repo string, number string) error {
	opt := sdk.IssueUpdateParam{Repo: repo, State: sdk.StatusClosed}
	_, err := c.UpdateIssue(owner, number, opt)
	return formatErr(err, "close issue")
}

func (c *client) ReopenIssue(owner, repo string, number string) error {
	opt := sdk.IssueUpdateParam{Repo: repo, State: sdk.StatusOpen}
	_, err := c.UpdateIssue(owner, number, opt)
	return formatErr(err, "reopen issue")
}

func (c *client) UpdateIssue(owner, number string, param sdk.IssueUpdateParam) (sdk.Issue, error) {
	issue, _, err := c.ac.IssuesApi.PatchV5ReposOwnerIssuesNumber(context.Background(), owner, number, param)
	return issue, formatErr(err, "update issue")
}

func (c *client) GetIssueLabels(org, repo, number string) ([]sdk.Label, error) {
	labels, _, err := c.ac.LabelsApi.GetV5ReposOwnerRepoIssuesNumberLabels(context.Background(), org, repo, number, nil)
	return labels, formatErr(err, "get issue labels")
}

func (c *client) UpdateIssueComment(org, repo string, commentID int32, comment string) error {
	opt := sdk.IssueCommentPatchParam{Body: comment}
	_, _, err := c.ac.IssuesApi.PatchV5ReposOwnerRepoIssuesCommentsId(
		context.Background(), org, repo, commentID, opt)
	return formatErr(err, "update comment of issue")
}

func (c *client) GetIssue(org, repo, number string) (sdk.Issue, error) {
	issue, _, err := c.ac.IssuesApi.GetV5ReposOwnerRepoIssuesNumber(context.Background(), org, repo, number, nil)
	return issue, formatErr(err, "get issue")
}

func (c *client) ListIssueComments(org, repo, number string) ([]sdk.Note, error) {
	var r []sdk.Note

	p := int32(1)
	opt := sdk.GetV5ReposOwnerRepoIssuesNumberCommentsOpts{}
	for {
		opt.Page = optional.NewInt32(p)
		cs, _, err := c.ac.IssuesApi.GetV5ReposOwnerRepoIssuesNumberComments(
			context.Background(), org, repo, number, &opt)
		if err != nil {
			return nil, formatErr(err, "list comments of issue")
		}

		if len(cs) == 0 {
			break
		}

		r = append(r, cs...)
		p++
	}

	return r, nil
}

//GetRepoAllBranch get repository all branch
func (c *client) GetRepoAllBranch(org, repo string) ([]sdk.Branch, error) {
	branches, _, err := c.ac.RepositoriesApi.GetV5ReposOwnerRepoBranches(context.Background(), org, repo,
		nil)
	return branches, formatErr(err, "get repo all branch")
}

//GetPathContent Get the content under a specific repository
func (c *client) GetPathContent(org, repo, path, ref string) (sdk.Content, error) {
	op := sdk.GetV5ReposOwnerRepoContentsPathOpts{}
	op.Ref = optional.NewString(ref)

	content, _, err := c.ac.RepositoriesApi.GetV5ReposOwnerRepoContentsPath(
		context.Background(), org, repo, path, &op,
	)
	if err != nil {
		return content, formatErr(err, "get path content")
	}

	if content.DownloadUrl == "" {
		return content, formatErr(fmt.Errorf("file does not exist"), "get path content")
	}

	return content, nil
}

func (c *client) CreateFile(org, repo, branch, path, content, commitMsg string) (sdk.CommitContent, error) {
	opt := sdk.NewFileParam{
		Message: commitMsg,
		Branch:  branch,
		Content: base64.StdEncoding.EncodeToString([]byte(content)),
	}

	v, _, err := c.ac.RepositoriesApi.PostV5ReposOwnerRepoContentsPath(
		context.Background(), org, repo, path, opt,
	)

	return v, formatErr(err, "create file")
}

//GetDirectoryTree Get the directory tree under a specific repository branch or commit sha
func (c *client) GetDirectoryTree(org, repo, sha string, recursive int32) (sdk.Tree, error) {
	op := sdk.GetV5ReposOwnerRepoGitTreesShaOpts{Recursive: optional.NewInt32(recursive)}
	trees, _, err := c.ac.GitDataApi.GetV5ReposOwnerRepoGitTreesSha(context.Background(), org, repo, sha, &op)
	return trees, formatErr(err, "get directory tree")
}

// GetUserPermissionsOfRepo get user permissions in the repository
func (c *client) GetUserPermissionsOfRepo(org, repo, login string) (sdk.ProjectMemberPermission, error) {
	permission, _, err := c.ac.RepositoriesApi.GetV5ReposOwnerRepoCollaboratorsUsernamePermission(
		context.Background(), org, repo, login, nil,
	)

	return permission, formatErr(err, "get user permissions")
}

// CreateRepoLabel create label for the repository
func (c *client) CreateRepoLabel(org, repo, label, color string) error {
	if color == "" {
		color = genrateRGBColor()
	}
	param := sdk.LabelPostParam{
		Name:  label,
		Color: color,
	}

	_, _, err := c.ac.LabelsApi.PostV5ReposOwnerRepoLabels(context.Background(), org, repo, param)

	return formatErr(err, "create a repo label")
}

func (c *client) CreateBranch(org, repo, branch, parentBranch string) error {
	_, _, err := c.ac.RepositoriesApi.PostV5ReposOwnerRepoBranches(
		context.Background(), org, repo,
		sdk.CreateBranchParam{BranchName: branch, Refs: parentBranch},
	)

	return formatErr(err, "create a branch")
}

func (c *client) SetProtectionBranch(org, repo, branch string) error {
	_, _, err := c.ac.RepositoriesApi.PutV5ReposOwnerRepoBranchesBranchProtection(
		context.Background(), org, repo, branch, sdk.BranchProtectionPutParam{},
	)

	return formatErr(err, "set protection branch")
}

func (c *client) CancelProtectionBranch(org, repo, branch string) error {
	_, err := c.ac.RepositoriesApi.DeleteV5ReposOwnerRepoBranchesBranchProtection(
		context.Background(), org, repo, branch,
		&sdk.DeleteV5ReposOwnerRepoBranchesBranchProtectionOpts{},
	)

	return formatErr(err, "cancel protection branch")
}

func (c *client) AddRepoMember(org, repo, login, permission string) error {
	_, _, err := c.ac.RepositoriesApi.PutV5ReposOwnerRepoCollaboratorsUsername(
		context.Background(), org, repo, login,
		sdk.ProjectMemberPutParam{
			Permission: permission,
		},
	)

	return formatErr(err, "add repo member")
}

func (c *client) RemoveRepoMember(org, repo, login string) error {
	v, err := c.ac.RepositoriesApi.DeleteV5ReposOwnerRepoCollaboratorsUsername(
		context.Background(), org, repo, login, nil,
	)

	if err == nil || (v != nil && v.StatusCode == 404) {
		return nil
	}
	return formatErr(err, "remove repo member")
}

func (c *client) CreateRepo(org string, repo sdk.RepositoryPostParam) error {
	_, _, err := c.ac.RepositoriesApi.PostV5OrgsOrgRepos(
		context.Background(), org, repo,
	)

	return formatErr(err, "create repo")
}

func (c *client) SetRepoReviewer(org, repo string, reviewer sdk.SetRepoReviewer) error {
	_, err := c.ac.RepositoriesApi.PutV5ReposOwnerRepoReviewer(
		context.Background(), org, repo, reviewer,
	)

	return formatErr(err, "set repo reviewer")
}

func (c *client) UpdateRepo(org, repo string, info sdk.RepoPatchParam) error {
	_, _, err := c.ac.RepositoriesApi.PatchV5ReposOwnerRepo(
		context.Background(), org, repo, info,
	)

	return formatErr(err, "update repo")
}

func (c *client) GetEnterprisesMember(enterprise, login string) (sdk.EnterpriseMember, error) {
	member, _, err := c.ac.EnterprisesApi.GetV5EnterprisesEnterpriseMembersUsername(
		context.Background(), enterprise, login, nil,
	)

	return member, formatErr(err, "get enterprise")
}

func (c *client) AddProjectLabels(org, repo string, label []string) error {
	opt := sdk.PullRequestLabelPostParam{Body: label}
	_, _, err := c.ac.LabelsApi.PostV5ReposOwnerRepoProjectLabels(
		context.Background(), org, repo, opt)
	return formatErr(err, "add project label")
}

func (c *client) UpdateProjectLabels(org, repo string, label []string) error {
	opt := sdk.PullRequestLabelPostParam{Body: label}
	_, _, err := c.ac.LabelsApi.PutV5ReposOwnerRepoProjectLabels(
		context.Background(), org, repo, opt)
	return formatErr(err, "update project label")
}

func (c *client) CreateIssue(org, repo, title, body string) (sdk.Issue, error) {
	param := sdk.IssueCreateParam{
		Repo:  repo,
		Body:  body,
		Title: title,
	}

	issue, _, err := c.ac.IssuesApi.PostV5ReposOwnerIssues(context.Background(), org, param)

	return issue, formatErr(err, "create issue")
}

func formatErr(err error, doWhat string) error {
	if err == nil {
		return err
	}

	var msg []byte
	if v, ok := err.(sdk.GenericSwaggerError); ok {
		msg = v.Body()
	}

	return fmt.Errorf("failed to %s, err: %s, msg: %q", doWhat, err.Error(), msg)
}
