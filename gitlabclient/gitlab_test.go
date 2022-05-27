package gitlabclient

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"testing"
)

func TestNewGitlabClient(t *testing.T) {
	cl, err := NewGitlabClient("xxx", "https://source.openeuler.sh/api/v4")
	if err != nil {
		fmt.Println("err:         ", err)
	}

	// create merge request
	//title := "test mr"
	//labels := gitlab.Labels{"wh-mr", "test-mr"}
	//desc := "test create mr"
	//src := "wh"
	//tar := "master"
	//opt := &gitlab.CreateMergeRequestOptions{
	//	Title:        &title,
	//	Labels:       &labels,
	//	Description:  &desc,
	//	SourceBranch: &src,
	//	TargetBranch: &tar,
	//}
	//_, _, err = cl.MergeRequests.CreateMergeRequest("19114", opt)
	//if err != nil {
	//	fmt.Println("create mr failed : ", err)
	//}

	// add comments for mr
	//body := "test comment"
	//_, _, err = cl.Notes.CreateMergeRequestNote("19114", 2, &gitlab.CreateMergeRequestNoteOptions{
	//	Body: &body,
	//})

	// get mr
	mr, _, err := cl.MergeRequests.GetMergeRequest("19114", 2, &gitlab.GetMergeRequestsOptions{})
	fmt.Println(mr)
}
