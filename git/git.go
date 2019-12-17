package git

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sawadashota/orb-update/driver"
	"github.com/sawadashota/orb-update/filesystem"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// Git .
type Git interface {
	BaseBranch() string
	Switch(branch string, create bool) error
	SwitchBack() error
	Commit(message string, path string) (CommitHash, error)
	Push(ctx context.Context, branch string) error
}

// CommitHash
type CommitHash string

// String .
func (ch *CommitHash) String() string {
	return string(*ch)
}

// DefaultGitClient .
type DefaultGitClient struct {
	d    driver.Driver
	repo *git.Repository
	base *plumbing.Reference
}

// Clone repository to memory
func Clone(d driver.Driver, owner, name string) (Git, filesystem.Filesystem, error) {
	fs := memfs.New()

	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: fmt.Sprintf(
			"https://%s:%s@github.com/%s/%s.git",
			d.Configuration().GithubUsername(),
			d.Configuration().GithubToken(),
			owner,
			name,
		),
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", d.Configuration().BaseBranch())),
	})

	if err != nil {
		return nil, nil, err
	}

	head, err := repo.Head()
	if err != nil {
		return nil, nil, err
	}

	return &DefaultGitClient{
		d:    d,
		repo: repo,
		base: head,
	}, filesystem.NewMemory(fs), nil
}

// OpenCurrentDirectoryRepository opens current directory's repository
func OpenCurrentDirectoryRepository(d driver.Driver) (Git, filesystem.Filesystem, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	repo, err := git.PlainOpen(pwd)
	if err != nil {
		return nil, nil, err
	}

	head, err := repo.Head()
	if err != nil {
		return nil, nil, err
	}

	return &DefaultGitClient{
		d:    d,
		repo: repo,
		base: head,
	}, filesystem.NewOs(), nil
}

// BaseBranch of git
func (d *DefaultGitClient) BaseBranch() string {
	return strings.ReplaceAll(d.base.Name().String(), "refs/heads/", "")
}

// Switch to branch
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

// SwitchBack tp origin branch
func (d *DefaultGitClient) SwitchBack() error {
	w, err := d.repo.Worktree()
	if err != nil {
		return err
	}

	return w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + d.BaseBranch()),
		Force:  true,
	})
}

// Commit .
func (d *DefaultGitClient) Commit(message string, path string) (CommitHash, error) {
	w, err := d.repo.Worktree()
	if err != nil {
		return "", err
	}

	if _, err := w.Add(path); err != nil {
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

	return CommitHash(h.String()), nil
}

// Push to origin
func (d *DefaultGitClient) Push(ctx context.Context, branch string) error {
	ref := fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch)
	return d.repo.PushContext(ctx, &git.PushOptions{
		RemoteName: git.DefaultRemoteName,
		RefSpecs:   []config.RefSpec{config.RefSpec(ref)},
		Auth: &http.BasicAuth{
			Username: d.d.Configuration().GithubUsername(),
			Password: d.d.Configuration().GithubToken(),
		},
	})
}
