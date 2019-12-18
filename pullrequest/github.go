package pullrequest

import (
	"context"
	"fmt"

	"github.com/google/go-github/v28/github"
	"github.com/sawadashota/orb-update/orb"
	"golang.org/x/oauth2"
)

// GitHubPullRequest .
type GitHubPullRequest struct {
	c          Configuration
	client     *github.Client
	owner      string
	repo       string
	difference orb.Difference
}

// NewGitHubPullRequest .
func NewGitHubPullRequest(ctx context.Context, c Configuration, owner string, repo string, diff *orb.Difference) (Creator, error) {
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.GithubToken()},
	))

	return &GitHubPullRequest{
		c:          c,
		client:     github.NewClient(tc),
		owner:      owner,
		repo:       repo,
		difference: *diff,
	}, nil
}

// Create Pull Request on GitHub
func (g *GitHubPullRequest) Create(ctx context.Context, message, baseBranch string) error {
	o := g.difference.New
	_, _, err := g.client.PullRequests.Create(ctx, g.owner, g.repo, &github.NewPullRequest{
		Title: github.String(fmt.Sprintf("orb: Bump %s/%s from %s to %s", o.Namespace(), o.Name(), g.difference.Old.Version(), o.Version())),
		Body:  &message,
		Base:  github.String(g.c.BaseBranch()),
		Head:  github.String(baseBranch),
	})

	return err
}

// AlreadyCreated Pull Request or not
func (g *GitHubPullRequest) AlreadyCreated(ctx context.Context, branch string) (bool, error) {
	prs, _, err := g.client.PullRequests.List(ctx, g.owner, g.repo, &github.PullRequestListOptions{
		State: "open",
		Head:  branch,
		Base:  g.c.BaseBranch(),
	})
	if err != nil {
		return false, err
	}

	for _, pr := range prs {
		if pr.Head.GetRef() == branch {
			return true, nil
		}
	}
	return false, nil
}
