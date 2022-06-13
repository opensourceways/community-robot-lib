package gitlabclient

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"strings"
)

var _ Client = (*client)(nil)

type client struct {
	ac *gitlab.Client
}

func NewGitlabClient(getToken func() []byte, host string) Client {
	tc := string(getToken())
	opts := gitlab.WithBaseURL(host)

	c, err := gitlab.NewOAuthClient(tc, opts)
	if err != nil {
		return &client{}
	}
	return &client{ac: c}
}

func (cli *client) GetMergeRequest(pid interface{}, mrID int) (gitlab.MergeRequest, error) {
	opts := &gitlab.GetMergeRequestsOptions{}
	r, _, err := cli.ac.MergeRequests.GetMergeRequest(pid, mrID, opts)
	return *r, err
}

func (cli *client) UpdateMergeRequest(pid interface{}, mrID int, options gitlab.UpdateMergeRequestOptions) (gitlab.MergeRequest, error) {
	r, _, err := cli.ac.MergeRequests.UpdateMergeRequest(pid, mrID, &options)
	return *r, err
}

func (cli *client) ListCollaborators(pid interface{}) ([]*gitlab.ProjectMember, error) {
	page := 1
	var r []*gitlab.ProjectMember
	for {
		lp := gitlab.ListOptions{
			Page:    page,
			PerPage: 100,
		}
		op := gitlab.ListProjectMembersOptions{ListOptions: lp}
		members, _, err := cli.ac.ProjectMembers.ListProjectMembers(pid, &op)
		if err != nil {
			return nil, fmt.Errorf(err.Error(), "list members failed")
		}

		if len(members) == 0 {
			break
		}

		r = append(r, members...)
		page++
	}

	return r, nil
}

func (cli *client) IsCollaborator(pid interface{}, userID int) (bool, error) {
	user, _, err := cli.ac.ProjectMembers.GetProjectMember(pid, userID)
	if err != nil || user == nil {
		return false, err
	}
	return true, nil
}

func (cli *client) AddProjectMember(pid interface{}, userID interface{}) error {
	opts := &gitlab.AddProjectMemberOptions{UserID: userID}
	_, _, err := cli.ac.ProjectMembers.AddProjectMember(pid, opts)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) RemoveProjectMember(pid interface{}, userID int) error {
	_, err := cli.ac.ProjectMembers.DeleteProjectMember(pid, userID)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) IsMember(gid interface{}, userID int) (bool, error) {
	gm, _, err := cli.ac.GroupMembers.GetGroupMember(gid, userID)
	if err != nil || gm == nil {
		return false, err
	}

	return true, nil
}

func (cli *client) GetMergeRequestChanges(pid interface{}, mrID int) ([]string, error) {
	mr, _, err := cli.ac.MergeRequests.GetMergeRequestChanges(pid, mrID, &gitlab.GetMergeRequestChangesOptions{})
	if err != nil {
		return []string{}, err
	}
	changedFiles := make([]string, len(mr.Changes))
	for _, c := range mr.Changes {
		if c.DeletedFile {
			changedFiles = append(changedFiles, c.NewPath)
		}
		changedFiles = append(changedFiles, c.NewPath)
	}

	return changedFiles, nil
}

func (cli *client) GetMergeRequestLabels(pid interface{}, mrID int) (gitlab.Labels, error) {
	mr, _, err := cli.ac.MergeRequests.GetMergeRequest(pid, mrID, &gitlab.GetMergeRequestsOptions{})
	if err != nil {
		return gitlab.Labels{}, err
	}

	return mr.Labels, nil
}

func (cli *client) ListMergeRequestComments(pid interface{}, mrID int) ([]*gitlab.Note, error) {
	var notes []*gitlab.Note

	page := 1
	for {
		ls := gitlab.ListOptions{Page: page, PerPage: 100}
		opts := &gitlab.ListMergeRequestNotesOptions{ListOptions: ls}
		comments, _, err := cli.ac.Notes.ListMergeRequestNotes(pid, mrID, opts)

		if err != nil {
			return notes, fmt.Errorf(err.Error(), "get comments for mr failed")
		}

		if len(comments) == 0 {
			break
		}

		notes = append(notes, comments...)
		page++
	}

	return notes, nil
}

func (cli *client) ListIssues(pid interface{}) ([]*gitlab.Issue, error) {
	var issueList []*gitlab.Issue

	page := 1
	for {
		ls := gitlab.ListOptions{
			Page:    page,
			PerPage: 100,
		}
		opts := &gitlab.ListProjectIssuesOptions{ListOptions: ls}
		issues, _, err := cli.ac.Issues.ListProjectIssues(pid, opts)

		if err != nil {
			return issueList, fmt.Errorf(err.Error(), "get issues for project failed")
		}

		if len(issues) == 0 {
			break
		}

		issueList = append(issueList, issues...)
		page++
	}

	return issueList, nil
}

func (cli *client) ListIssueRelatedMergeRequest(pid interface{}, issueID int) ([]*gitlab.MergeRequest, error) {
	var MRList []*gitlab.MergeRequest

	page := 1
	for {
		opts := &gitlab.ListMergeRequestsRelatedToIssueOptions{Page: page, PerPage: 100}
		mrs, _, err := cli.ac.Issues.ListMergeRequestsRelatedToIssue(pid, issueID, opts)

		if err != nil {
			return MRList, fmt.Errorf(err.Error(), "get issues for project failed")
		}

		if len(mrs) == 0 {
			break
		}

		MRList = append(MRList, mrs...)
		page++
	}

	return MRList, nil
}

func (cli *client) DeleteMergeRequestComment(pid interface{}, mrID int, noteID int) error {
	_, err := cli.ac.Notes.DeleteMergeRequestNote(pid, mrID, noteID)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) CreateMergeRequestComment(pid interface{}, mrID int, comment string) error {
	_, _, err := cli.ac.Notes.CreateMergeRequestNote(pid, mrID, &gitlab.CreateMergeRequestNoteOptions{Body: &comment})
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) UpdateMergeRequestComment(pid interface{}, mrID, noteID int, comment string) error {
	_, _, err := cli.ac.Notes.UpdateMergeRequestNote(pid, mrID, noteID, &gitlab.UpdateMergeRequestNoteOptions{Body: &comment})
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) AddMergeRequestLabel(pid interface{}, mrID int, labels gitlab.Labels) error {
	opts := gitlab.UpdateMergeRequestOptions{AddLabels: &labels}
	_, err := cli.UpdateMergeRequest(pid, mrID, opts)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) RemoveMergeRequestLabel(pid interface{}, mrID int, labels gitlab.Labels) error {
	opts := gitlab.UpdateMergeRequestOptions{RemoveLabels: &labels}
	_, err := cli.UpdateMergeRequest(pid, mrID, opts)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) ReplaceMergeRequestAllLabels(pid interface{}, mrID int, labels gitlab.Labels) error {
	opts := gitlab.UpdateMergeRequestOptions{Labels: &labels}
	_, err := cli.UpdateMergeRequest(pid, mrID, opts)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) CloseMergeRequest(pid interface{}, mrID int, state string) error {
	opts := gitlab.UpdateMergeRequestOptions{StateEvent: &state}
	_, err := cli.UpdateMergeRequest(pid, mrID, opts)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) ReopenMergeRequest(pid interface{}, mrID int, state string) error {
	opts := gitlab.UpdateMergeRequestOptions{StateEvent: &state}
	_, err := cli.UpdateMergeRequest(pid, mrID, opts)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) AssignMergeRequest(pid interface{}, mrID int, ids []int) error {
	opts := gitlab.UpdateMergeRequestOptions{AssigneeIDs: &ids}
	_, err := cli.UpdateMergeRequest(pid, mrID, opts)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) UnAssignMergeRequest(pid interface{}, mrID int, ids []int) error {
	opts := gitlab.UpdateMergeRequestOptions{AssigneeIDs: &ids}
	_, err := cli.UpdateMergeRequest(pid, mrID, opts)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) GetMergeRequestCommits(pid interface{}, mrID int) ([]*gitlab.Commit, error) {
	var commits []*gitlab.Commit
	page := 1
	for {
		opts := &gitlab.GetMergeRequestCommitsOptions{Page: page, PerPage: 100}
		cts, _, err := cli.ac.MergeRequests.GetMergeRequestCommits(pid, mrID, opts)

		if err != nil {
			return commits, fmt.Errorf(err.Error(), "get commit of mergerequest failed")
		}

		if len(cts) == 0 {
			break
		}

		commits = append(commits, cts...)
		page++
	}

	return commits, nil
}

func (cli *client) GetSingleRepoCommit(pid interface{}, sha string) (*gitlab.Commit, error) {
	commit, _, err := cli.ac.Commits.GetCommit(pid, sha)
	if err != nil {
		return commit, err
	}

	return commit, nil
}

func (cli *client) MergeMergeRequest(pid interface{}, mrID int) error {
	_, _, err := cli.ac.MergeRequests.AcceptMergeRequest(pid, mrID, &gitlab.AcceptMergeRequestOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) GetGroups() ([]*gitlab.Group, error) {
	var groups []*gitlab.Group

	page := 1
	for {
		ls := gitlab.ListOptions{
			Page:    page,
			PerPage: 100,
		}
		opts := gitlab.ListGroupsOptions{ListOptions: ls}
		grps, _, err := cli.ac.Groups.ListGroups(&opts)

		if err != nil {
			return groups, err
		}

		if len(grps) == 0 {
			break
		}

		groups = append(groups, grps...)
		page++
	}

	return groups, nil
}

func (cli *client) GetProjects(gid interface{}) ([]*gitlab.Project, error) {
	var projects []*gitlab.Project

	page := 1
	for {
		ls := gitlab.ListOptions{
			Page:    page,
			PerPage: 100,
		}
		opts := gitlab.ListGroupProjectsOptions{ListOptions: ls}
		prjs, _, err := cli.ac.Groups.ListGroupProjects(gid, &opts)

		if err != nil {
			return projects, err
		}

		if len(prjs) == 0 {
			break
		}

		projects = append(projects, prjs...)
		page++
	}

	return projects, nil
}

func (cli *client) GetProject(pid interface{}) (*gitlab.Project, error) {
	prj, _, err := cli.ac.Projects.GetProject(pid, &gitlab.GetProjectOptions{})

	if err != nil {
		return prj, err
	}

	return prj, nil
}

func (cli *client) CreateProject(opts gitlab.CreateProjectOptions) error {
	_, _, err := cli.ac.Projects.CreateProject(&opts)

	if err != nil {
		return err
	}

	return nil
}

func (cli *client) UpdateProject(pid interface{}, opts gitlab.EditProjectOptions) error {
	_, _, err := cli.ac.Projects.EditProject(pid, &opts)

	if err != nil {
		return err
	}

	return nil
}

func (cli *client) AddProjectLabel(pid interface{}, label, color string) error {
	_, _, err := cli.ac.Labels.CreateLabel(pid, &gitlab.CreateLabelOptions{Name: &label, Color: &color})
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) UpdateProjectLabel(pid interface{}, oldLabel, label, color string) error {
	_, _, err := cli.ac.Labels.UpdateLabel(pid, &gitlab.UpdateLabelOptions{Name: &oldLabel, NewName: &label, Color: &color})
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) GetProjectLabels(pid interface{}) ([]*gitlab.Label, error) {
	var labels []*gitlab.Label
	page := 1
	for {
		ls := gitlab.ListOptions{
			Page:    page,
			PerPage: 100,
		}
		opts := gitlab.ListLabelsOptions{ListOptions: ls}
		lbs, _, err := cli.ac.Labels.ListLabels(pid, &opts)
		if err != nil {
			return labels, err
		}

		if len(lbs) == 0 {
			break
		}

		labels = append(labels, lbs...)
		page++
	}

	return labels, nil
}

func (cli *client) AssignIssue(pid interface{}, issueID int, assignees []int) error {
	opts := gitlab.UpdateIssueOptions{AssigneeIDs: &assignees}
	_, _, err := cli.ac.Issues.UpdateIssue(pid, issueID, &opts)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) UpdateIssue(pid interface{}, issueID int, opts gitlab.UpdateIssueOptions) error {
	_, _, err := cli.ac.Issues.UpdateIssue(pid, issueID, &opts)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) CreateIssue(pid interface{}, opts gitlab.CreateIssueOptions) error {
	_, _, err := cli.ac.Issues.CreateIssue(pid, &opts)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) UnAssignIssue(pid interface{}, issueID int, assignees []int) error {
	opts := gitlab.UpdateIssueOptions{AssigneeIDs: &assignees}
	err := cli.UpdateIssue(pid, issueID, opts)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) RemoveAssignIssue(pid interface{}, issueID int) error {
	var assignIDs []int
	opts := gitlab.UpdateIssueOptions{AssigneeIDs: &assignIDs}
	err := cli.UpdateIssue(pid, issueID, opts)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) CreateIssueComment(pid interface{}, issueID int, comment string) error {
	opts := &gitlab.CreateIssueNoteOptions{Body: &comment}
	_, _, err := cli.ac.Notes.CreateIssueNote(pid, issueID, opts)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) ListIssueComments(pid interface{}, issueID int) ([]*gitlab.Note, error) {
	var comments []*gitlab.Note
	page := 1
	for {
		ls := gitlab.ListOptions{
			Page:    page,
			PerPage: 100,
		}
		opts := gitlab.ListIssueNotesOptions{ListOptions: ls}
		cmts, _, err := cli.ac.Notes.ListIssueNotes(pid, issueID, &opts)
		if err != nil {
			return comments, err
		}

		if len(cmts) == 0 {
			break
		}

		comments = append(comments, cmts...)
		page++
	}

	return comments, nil
}

func (cli *client) UpdateIssueComment(pid interface{}, issueID, noteID int, comment string) error {
	opts := &gitlab.UpdateIssueNoteOptions{Body: &comment}
	_, _, err := cli.ac.Notes.UpdateIssueNote(pid, issueID, noteID, opts)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) RemoveIssueComment(pid interface{}, issueID, noteID int) error {
	_, err := cli.ac.Notes.DeleteIssueNote(pid, issueID, noteID)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) GetIssueLabels(pid interface{}, issueID int) ([]string, error) {
	issue, _, err := cli.ac.Issues.GetIssue(pid, issueID)
	if err != nil {
		return nil, err
	}

	return issue.Labels, nil
}

func (cli *client) RemoveIssueLabels(pid interface{}, issueID int, labels gitlab.Labels) error {
	err := cli.UpdateIssue(pid, issueID, gitlab.UpdateIssueOptions{RemoveLabels: &labels})
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) AddIssueLabels(pid interface{}, issueID int, labels gitlab.Labels) error {
	err := cli.UpdateIssue(pid, issueID, gitlab.UpdateIssueOptions{AddLabels: &labels})
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) CloseIssue(pid interface{}, issueID int) error {
	action := "close"
	err := cli.UpdateIssue(pid, issueID, gitlab.UpdateIssueOptions{StateEvent: &action})
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) ReopenIssue(pid interface{}, issueID int) error {
	action := "reopen"
	err := cli.UpdateIssue(pid, issueID, gitlab.UpdateIssueOptions{StateEvent: &action})
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) GetSingleIssue(pid interface{}, issueID int) (*gitlab.Issue, error) {
	issue, _, err := cli.ac.Issues.GetIssue(pid, issueID)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func (cli *client) CreateBranch(pid interface{}, branch, ref string) error {
	_, _, err := cli.ac.Branches.CreateBranch(pid, &gitlab.CreateBranchOptions{Branch: &branch, Ref: &ref})
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) GetProjectAllBranches(pid interface{}) ([]*gitlab.Branch, error) {
	var branches []*gitlab.Branch

	page := 1
	for {
		ls := gitlab.ListOptions{Page: page, PerPage: 100}
		brs, _, err := cli.ac.Branches.ListBranches(pid, &gitlab.ListBranchesOptions{ListOptions: ls})

		if err != nil {
			return branches, err
		}

		if len(brs) == 0 {
			break
		}

		branches = append(branches, brs...)

		page++
	}

	return branches, nil
}

func (cli *client) SetProtectionBranch(pid interface{}, branch string) error {
	forcePush := false
	_, _, err := cli.ac.ProtectedBranches.ProtectRepositoryBranches(pid,
		&gitlab.ProtectRepositoryBranchesOptions{Name: &branch, AllowForcePush: &forcePush})

	if err != nil {
		return err
	}

	return nil
}

func (cli *client) UnProtectBranch(pid interface{}, branch string) error {
	_, err := cli.ac.ProtectedBranches.UnprotectRepositoryBranches(pid, branch)
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) CreateFile(pid interface{}, file string, opts gitlab.CreateFileOptions) error {
	_, _, err := cli.ac.RepositoryFiles.CreateFile(pid, file, &opts)

	if err != nil {
		return err
	}

	return nil
}

func (cli *client) GetPathContent(pid interface{}, file, filePath string) (*gitlab.File, error) {
	fileInfo, _, err := cli.ac.RepositoryFiles.GetFile(pid, file, &gitlab.GetFileOptions{Ref: &filePath})
	if err != nil {
		return fileInfo, err
	}

	return fileInfo, nil
}

func (cli *client) GetDirectoryTree(pid interface{}, opts gitlab.ListTreeOptions) ([]*gitlab.TreeNode, error) {
	var t []*gitlab.TreeNode

	page := 1
	for {
		ls := gitlab.ListOptions{Page: page, PerPage: 100}
		opts.ListOptions = ls
		trees, _, err := cli.ac.Repositories.ListTree(pid, &opts)

		if err != nil {
			return t, err
		}

		if len(trees) == 0 {
			break
		}

		t = append(t, trees...)
		page++
	}

	return t, nil
}

func (cli *client) GetUserPermissionOfProject(pid interface{}, userID int) (bool, error) {
	members, err := cli.ListCollaborators(pid)
	if err != nil {
		return false, err
	}

	for _, m := range members {
		if m.ID != userID {
			return false, nil
		}

		if m.ID == userID {
			if m.AccessLevel == 30 || m.AccessLevel == 40 || m.AccessLevel == 50 {
				return true, nil
			}
		}
	}

	return false, nil
}

func (cli *client) CreateProjectLabel(pid interface{}, label, color string) error {
	_, _, err := cli.ac.Labels.CreateLabel(pid, &gitlab.CreateLabelOptions{Name: &label, Color: &color})
	if err != nil {
		return err
	}

	return nil
}

func (cli *client) GetMergeRequestLabelChanges(pid interface{}, mrID int) ([]*gitlab.LabelEvent, error) {
	var labelEvents []*gitlab.LabelEvent

	page := 1
	for {
		ls := gitlab.ListOptions{Page: page, PerPage: 100}
		es, _, err := cli.ac.ResourceLabelEvents.ListMergeRequestsLabelEvents(pid, mrID,
			&gitlab.ListLabelEventsOptions{ListOptions: ls})

		if err != nil {
			return labelEvents, err
		}

		if len(es) == 0 {
			break
		}

		labelEvents = append(labelEvents, es...)
		page++
	}

	return labelEvents, nil
}

func GetMROrgAndRepo(e *gitlab.MergeEvent) (org, repo string) {
	org = strings.Split(e.Project.PathWithNamespace, "/")[0]
	repo = strings.Split(e.Project.PathWithNamespace, "/")[1]
	return org, repo
}

func GetMRAuthor(e *gitlab.MergeEvent) (author string) {
	return e.User.Username
}

func GetMRNumber(e *gitlab.MergeEvent) (mrID int) {
	return e.ObjectAttributes.IID
}

func GetIssueOrgAndRepo(e *gitlab.IssueEvent) (org, repo string) {
	org = strings.Split(e.Project.PathWithNamespace, "/")[0]
	repo = strings.Split(e.Project.PathWithNamespace, "/")[1]
	return org, repo
}

func GetIssueAuthor(e *gitlab.IssueEvent) (author string) {
	return e.User.Username
}

func GetIssueNumber(e *gitlab.IssueEvent) (issueID int) {
	return e.ObjectAttributes.IID
}

func CheckSourceBranchChanged(e *gitlab.MergeEvent) (changed bool) {
	if e == nil {
		return false
	}
	if e.ObjectAttributes.OldRev != "" && e.ObjectAttributes.OldRev != e.ObjectAttributes.LastCommit.ID {
		return true
	}

	return false
}

func CheckLabelUpdate(e *gitlab.MergeEvent) (labelUpdated bool) {
	if e == nil {
		return false
	}
	pre := make(map[string]string, len(e.Changes.Labels.Previous))
	cur := make(map[string]string, len(e.Changes.Labels.Current))
	for _, i := range e.Changes.Labels.Previous {
		pre[i.Name] = i.Name
	}
	for _, j := range e.Changes.Labels.Current {
		cur[j.Name] = j.Name
	}

	for v := range cur {
		if _, ok := pre[v]; ok {
			continue
		} else {
			return true
		}
	}

	return false
}

func GetMRCommentOrgAndRepo(e *gitlab.MergeCommentEvent) (org, repo string) {
	org = strings.Split(e.Project.PathWithNamespace, "/")[0]
	repo = strings.Split(e.Project.PathWithNamespace, "/")[1]
	return org, repo
}

func GetMRCommentAuthor(e *gitlab.MergeCommentEvent) (author string) {
	return e.User.Username
}

func GetMRCommentAuthorID(e *gitlab.MergeCommentEvent) (authorID int) {
	return e.User.ID
}

func GetMRCommentBody(e *gitlab.MergeCommentEvent) (comment string) {
	return e.ObjectAttributes.Note
}

func GetIssueCommentOrgAndRepo(e *gitlab.IssueCommentEvent) (org, repo string) {
	org = strings.Split(e.Project.PathWithNamespace, "/")[0]
	repo = strings.Split(e.Project.PathWithNamespace, "/")[1]
	return org, repo
}

func GetIssueCommentAuthor(e *gitlab.IssueCommentEvent) (author string) {
	return e.User.Username
}

func GetIssueCommentAuthorID(e *gitlab.IssueCommentEvent) (authorID int) {
	return e.User.ID
}

func GetIssueCommentBody(e *gitlab.IssueCommentEvent) (comment string) {
	return e.ObjectAttributes.Note
}
