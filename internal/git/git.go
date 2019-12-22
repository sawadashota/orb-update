package git

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/sawadashota/orb-update/internal/filesystem"
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
	Commit(message string, path []string) (CommitHash, error)
	Push(ctx context.Context, branch string) error
}

// CommitHash .
type CommitHash string

// String .
func (ch *CommitHash) String() string {
	return string(*ch)
}

// DefaultGitClient .
type DefaultGitClient struct {
	c    Configuration
	repo *git.Repository
	base *plumbing.Reference
}

// Clone repository to memory
func Clone(c Configuration, owner, name string) (Git, filesystem.Filesystem, error) {
	fs := memfs.New()

	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: fmt.Sprintf(
			"https://%s:%s@github.com/%s/%s.git",
			c.GithubUsername(),
			c.GithubToken(),
			owner,
			name,
		),
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", c.BaseBranch())),
	})

	if err != nil {
		return nil, nil, errors.Errorf(`failed to git clone because "%s"`, err)
	}

	head, err := repo.Head()
	if err != nil {
		return nil, nil, errors.Errorf(`failed to retrieve git head because "%s"`, err)
	}

	return &DefaultGitClient{
		c:    c,
		repo: repo,
		base: head,
	}, filesystem.NewMemory(fs), nil
}

// OpenCurrentDirectoryRepository opens current directory's repository
func OpenCurrentDirectoryRepository(c Configuration) (Git, filesystem.Filesystem, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, nil, errors.Errorf(`failed to retrieve because "%s"`, err)
	}

	repo, err := git.PlainOpen(pwd)
	if err != nil {
		return nil, nil, errors.Errorf(`failed to open git repository because "%s"`, err)
	}

	head, err := repo.Head()
	if err != nil {
		return nil, nil, errors.Errorf(`failed to retrieve git head because "%s"`, err)
	}

	return &DefaultGitClient{
		c:    c,
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
		return errors.Errorf(`failed to git branch switch because "%s"`, err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + branch),
		Create: create,
		Keep:   true,
	})

	if err != nil {
		return errors.Errorf(`failed to switch branch because "%s"`, err)
	}

	return nil
}

// SwitchBack tp origin branch
func (d *DefaultGitClient) SwitchBack() error {
	w, err := d.repo.Worktree()
	if err != nil {
		return errors.Errorf(`Failed to git branch switch back because "%s"`, err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + d.BaseBranch()),
		Force:  true,
	})

	if err != nil {
		return errors.Errorf(`failed to switch branch because "%s"`, err)
	}

	return nil
}

// Commit .
func (d *DefaultGitClient) Commit(message string, path []string) (CommitHash, error) {
	w, err := d.repo.Worktree()
	if err != nil {
		return "", errors.Errorf(`failed to retrieve git work tree because "%s"`, err)
	}

	for _, p := range path {
		if _, err := w.Add(p); err != nil {
			return "", err
		}
	}

	h, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  d.c.GitAuthorName(),
			Email: d.c.GitAuthorEmail(),
			When:  time.Now(),
		},
	})
	if err != nil {
		return "", errors.Errorf(`failed to commit because "%s"`, err)
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
			Username: d.c.GithubUsername(),
			Password: d.c.GithubToken(),
		},
	})
}
