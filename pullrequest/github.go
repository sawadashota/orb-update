package pullrequest

import (
	"context"
	"fmt"

	"github.com/google/go-github/v28/github"
	"github.com/sawadashota/orb-update/driver"
	"github.com/sawadashota/orb-update/orb"
	"golang.org/x/oauth2"
)

// GitHubRelease .
type GitHubPullRequest struct {
	d          driver.Driver
	client     *github.Client
	owner      string
	repo       string
	difference orb.Difference
}

// NewGitHubPullRequest .
func NewGitHubPullRequest(ctx context.Context, d driver.Driver, owner string, repo string, diff *orb.Difference) (Creator, error) {
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: d.Configuration().GithubToken()},
	))

	return &GitHubPullRequest{
		d:          d,
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
		Base:  github.String(g.d.Configuration().BaseBranch()),
		Head:  github.String(baseBranch),
	})

	return err
}

// AlreadyCreated Pull Request or not
func (g *GitHubPullRequest) AlreadyCreated(ctx context.Context, branch string) (bool, error) {
	prs, _, err := g.client.PullRequests.List(ctx, g.owner, g.repo, &github.PullRequestListOptions{
		State: "open",
		Head:  branch,
		Base:  g.d.Configuration().BaseBranch(),
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
