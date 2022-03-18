package githubclient

import (
	"context"

	"github.com/google/go-github/v36/github"
	"golang.org/x/oauth2"
)

func NewGithubClient(getToken func() []byte) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: string(getToken())})
	tc := oauth2.NewClient(context.Background(), ts)

	return github.NewClient(tc)
}
