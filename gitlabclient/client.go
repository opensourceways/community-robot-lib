package gitlabclient

import (
	"github.com/xanzy/go-gitlab"
)

//func NewGitlabClient(getToken func() []byte) (*gitlab.Client, error) {
//	//ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: string(getToken())})
//	//tc := oauth2.NewClient(context.Background(), ts)
//	tc := string(getToken())
//
//	return gitlab.NewClient(tc)
//}

func NewGitlabClient(token, hostURL string) (*gitlab.Client, error) {
	//ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: string(getToken())})
	//tc := oauth2.NewClient(context.Background(), ts)
	tc := token
	opts := gitlab.WithBaseURL(hostURL)

	return gitlab.NewClient(tc, opts)
}
