package vcsuser

import (
	"context"

	"github.com/pkg/errors"

	"github.com/google/go-github/github"
	"github.com/sawadashota/orb-update/internal/git"
	"golang.org/x/oauth2"
)

// GithubClient .
type GithubClient struct {
	c      Configuration
	client *github.Client
}

// NewGithubClient .
func NewGithubClient(ctx context.Context, c git.Configuration) *GithubClient {
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.GithubToken()},
	))

	return &GithubClient{
		c:      c,
		client: github.NewClient(tc),
	}
}

// VCSUser .
type VCSUser interface {
	Name() string
	Email() string
}

type user struct {
	name  string
	email string
}

// Name of GitHub user
func (u *user) Name() string {
	return u.name
}

// Email of GitHub user
func (u *user) Email() string {
	return u.email
}

// Fetch user name and email from VCS
// public information in turn override private
func (gc *GithubClient) Fetch(ctx context.Context) (VCSUser, error) {
	gu, _, err := gc.client.Users.Get(ctx, gc.c.GithubUsername())
	if err != nil {
		return nil, errors.Errorf(`failed to fetch GitHub user information because "%s"`, err)
	}

	u := &user{
		name:  gu.GetName(),
		email: gu.GetEmail(),
	}

	if u.name == "" {
		u.name = gc.c.GithubUsername()
	}

	if u.email != "" {
		return u, nil
	}

	emails, _, err := gc.client.Users.ListEmails(ctx, &github.ListOptions{})
	if err != nil {
		return nil, errors.Errorf(`failed to fetch GitHub user's email because "%s"`, err)
	}

	for _, email := range emails {
		if email.GetPrimary() {
			u.email = email.GetEmail()
			break
		}
	}

	if u.email == "" {
		return nil, errors.New("GitHub primary email is not found")
	}

	return u, nil
}
