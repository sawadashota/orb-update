package pullrequest

import (
	"context"
	"fmt"

	"github.com/sawadashota/orb-update/internal/orb"

	"github.com/pkg/errors"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

// GitHubPullRequest .
type GitHubPullRequest struct {
	c      Configuration
	client *github.Client
	owner  string
	repo   string
}

// GitRepository .
type GitRepository interface {
	Owner() string
	Name() string
}

// NewGitHubPullRequest .
func NewGitHubPullRequest(ctx context.Context, r Registry, c Configuration) Creator {
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.GithubToken()},
	))

	return &GitHubPullRequest{
		c:      c,
		client: github.NewClient(tc),
		owner:  r.VCSRepository().Owner(),
		repo:   r.VCSRepository().Name(),
	}
}

// Create Pull Request on GitHub
func (g *GitHubPullRequest) Create(ctx context.Context, update *orb.Update, message, baseBranch string) error {
	o := update.After
	_, _, err := g.client.PullRequests.Create(ctx, g.owner, g.repo, &github.NewPullRequest{
		Title: github.String(fmt.Sprintf("orb: Bump %s/%s from %s to %s", o.Namespace(), o.Name(), update.Before.Version(), o.Version())),
		Body:  &message,
		Base:  github.String(g.c.BaseBranch()),
		Head:  github.String(baseBranch),
	})

	if err != nil {
		return errors.Errorf(`failed to create pull request because "%s"`, err)
	}

	return nil
}

// AlreadyCreated Pull Request or not
func (g *GitHubPullRequest) AlreadyCreated(ctx context.Context, branch string) (bool, error) {
	prs, _, err := g.client.PullRequests.List(ctx, g.owner, g.repo, &github.PullRequestListOptions{
		State: "open",
		Head:  branch,
		Base:  g.c.BaseBranch(),
	})
	if err != nil {
		return false, errors.Errorf(`failed to fetch pull request list because "%s"`, err)
	}

	for _, pr := range prs {
		if pr.Head.GetRef() == branch {
			return true, nil
		}
	}
	return false, nil
}
