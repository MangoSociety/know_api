package repository

import (
	"context"
	"github.com/MangoSociety/know_api/internal/migrator/domain"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
	"log"
	"strings"
)

type GitHubRepository interface {
	FetchFiles() ([]domain.GitHubFile, error)
}

type gitHubRepository struct {
	client *github.Client
	owner  string
	repo   string
	ref    string
	dir    string
}

func NewGitHubRepository(token, owner, repo, ref, dir string) GitHubRepository {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &gitHubRepository{
		client: client,
		owner:  owner,
		repo:   repo,
		ref:    ref,
		dir:    dir,
	}
}

func (r *gitHubRepository) FetchFiles() ([]domain.GitHubFile, error) {
	ctx := context.Background()
	tree, _, err := r.client.Git.GetTree(ctx, r.owner, r.repo, r.ref, true)
	if err != nil {
		return nil, err
	}

	var files []domain.GitHubFile
	for _, entry := range tree.Entries {
		if entry.GetType() == "blob" && entry.GetPath() != "" && strings.HasPrefix(entry.GetPath(), r.dir) {
			fileContent, _, _, err := r.client.Repositories.GetContents(ctx, r.owner, r.repo, entry.GetPath(), nil)
			if err != nil {
				log.Println("Error getting content for file:", entry.GetPath(), err)
				continue
			}
			content, err := fileContent.GetContent()
			if err != nil {
				log.Println("Error decoding content for file:", entry.GetPath(), err)
				continue
			}
			files = append(files, domain.GitHubFile{
				Path:    entry.GetPath(),
				Content: content,
			})
		}
	}

	return files, nil
}
