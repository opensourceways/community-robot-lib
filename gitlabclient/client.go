package gitlabclient

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
)

var _ Client = (*client)(nil)

type client struct {
	ac *gitlab.Client
}

//func NewGitlabClient(getToken func() []byte, host string) (*gitlab.Client, error) {
//	tc := string(getToken())
//	opts := gitlab.WithBaseURL(host)
//
//	return gitlab.NewOAuthClient(tc, opts)
//}

func NewGitlabClient(getToken func() []byte, host string) Client {
	tc := string(getToken())
	opts := gitlab.WithBaseURL(host)

	c, err := gitlab.NewOAuthClient(tc, opts)
	if err != nil {
		return &client{}
	}
	return &client{ac: c}
}

//func NewGitlabClient(token, hostURL string) (*gitlab.Client, error) {
//	//ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: string(getToken())})
//	//tc := oauth2.NewClient(context.Background(), ts)
//	tc := token
//	opts := gitlab.WithBaseURL(hostURL)
//
//	return gitlab.NewClient(tc, opts)
//}

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
