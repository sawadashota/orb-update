package pullrequest

import (
	"context"
	"fmt"

	"github.com/sawadashota/orb-update/driver"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v28/github"
	"github.com/sawadashota/orb-update/orb"
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

func (g *GitHubPullRequest) Create(ctx context.Context, message string) error {
	gc, err := NewDefaultGitClient(g.d)
	if err != nil {
		return err
	}

	fmt.Println("Switch")
	if err := gc.Switch(g.branch(), true); err != nil {
		return err
	}

	defer gc.SwitchBack()

	fmt.Println("Commit")
	if _, err := gc.Commit(message, g.branch()); err != nil {
		return err
	}

	fmt.Println("Push")
	if err := gc.Push(ctx, g.branch()); err != nil {
		return err
	}

	fmt.Println("PullRequests")
	o := g.difference.New
	_, _, err = g.client.PullRequests.Create(ctx, g.owner, g.repo, &github.NewPullRequest{
		Title: github.String(fmt.Sprintf("Bump %s/%s to %s", o.Namespace(), o.Name(), o.Version())),
		Body:  &message,
		Base:  github.String(gc.BaseBranch()),
		Head:  github.String(g.branch()),
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *GitHubPullRequest) branch() string {
	o := g.difference.New
	return fmt.Sprintf("orb-update/%s/%s-%s-aa", o.Namespace(), o.Name(), o.Version())
}
