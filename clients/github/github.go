package github

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Client struct {
	token  string
	Owner  string
	Repo   string
	Sha    string
	Client *github.Client
}

func NewGithubClient(token string, owner string, repo string, sha string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &Client{
		token:  token,
		Owner:  owner,
		Repo:   repo,
		Sha:    sha,
		Client: github.NewClient(tc),
	}
}
