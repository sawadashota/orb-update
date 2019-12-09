package pullrequest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/src-d/go-git.v4/config"

	"github.com/sawadashota/orb-update/driver"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type Git interface {
	BaseBranch() string
	Switch(branch string, create bool) error
	SwitchBack() error
	Commit(message string, branch string) (CommitHash, error)
	Push(ctx context.Context, branch string) error
}

type CommitHash string

func (ch *CommitHash) String() string {
	return string(*ch)
}

func (ch *CommitHash) hash() plumbing.Hash {
	return plumbing.NewHash(ch.String())
}

type DefaultGitClient struct {
	d    driver.Driver
	repo *git.Repository
	base *plumbing.Reference
}

func NewDefaultGitClient(d driver.Driver) (Git, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(pwd)
	if err != nil {
		return nil, err
	}

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	return &DefaultGitClient{
		d:    d,
		repo: repo,
		base: head,
	}, nil
}

func (d *DefaultGitClient) BaseBranch() string {
	return filepath.Base(d.base.Name().String())
}

func (d *DefaultGitClient) Switch(branch string, create bool) error {
	w, err := d.repo.Worktree()
	if err != nil {
		return err
	}

	return w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + branch),
		Create: create,
		Keep:   true,
	})
}

func (d *DefaultGitClient) SwitchBack() error {
	fmt.Printf("SwitchBack: %s\n", d.BaseBranch())
	w, err := d.repo.Worktree()
	if err != nil {
		return err
	}

	return w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + d.BaseBranch()),
	})
}

func (d *DefaultGitClient) Commit(message string, branch string) (CommitHash, error) {
	w, err := d.repo.Worktree()
	if err != nil {
		return "", err
	}

	if _, err := w.Add("."); err != nil {
		return "", err
	}

	h, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  d.d.Configuration().GitAuthorName(),
			Email: d.d.Configuration().GitAuthorEmail(),
			When:  time.Now(),
		},
	})
	if err != nil {
		return "", err
	}

	err = d.repo.Storer.SetReference(plumbing.NewReferenceFromStrings(branch, h.String()))
	if err != nil {
		return "", err
	}

	fmt.Printf("Hash: %s\n", h.String())

	return CommitHash(h.String()), nil
}

func (d *DefaultGitClient) Push(ctx context.Context, branch string) error {
	bs, err := d.repo.Branches()
	if err != nil {
		return err
	}

	for {
		ref, err := bs.Next()
		if err != nil {
			break
		}

		fmt.Println(ref.String())
	}

	ref := fmt.Sprintf("%s:%s", branch, branch)
	return d.repo.PushContext(ctx, &git.PushOptions{
		RemoteName: git.DefaultRemoteName,
		RefSpecs:   []config.RefSpec{config.RefSpec(ref)},
		Auth: &http.BasicAuth{
			Username: d.d.Configuration().GithubUsername(),
			Password: d.d.Configuration().GithubToken(),
		},
		Progress: os.Stdout,
	})
}
