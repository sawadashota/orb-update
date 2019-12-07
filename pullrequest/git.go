package pullrequest

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sawadashota/orb-update/driver"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type Git interface {
	BaseBranch() string
	Switch(branch string, create bool) error
	SwitchBack() error
	Commit(message string) (CommitHash, error)
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
	repo git.Repository
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
		repo: *repo,
		base: head,
	}, nil
}

func (d *DefaultGitClient) BaseBranch() string {
	return filepath.Base(d.base.Name().String())
}

func (d *DefaultGitClient) Switch(branch string, create bool) error {
	cmd := exec.Command("git", "switch", "-c", branch)
	cmd.Stdout = os.Stdout
	return cmd.Run()

	//w, err := d.repo.Worktree()
	//if err != nil {
	//	return err
	//}
	//
	////err = d.repo.CreateBranch(&config.Branch{
	////	Name:  branch,
	////	Merge: plumbing.ReferenceName("refs/heads/" + branch),
	////})
	////if err != nil {
	////	return err
	////}
	//
	//return w.Checkout(&git.CheckoutOptions{
	//	Branch: plumbing.ReferenceName(branch),
	//	//Create: true,
	//	Keep: true,
	//})
}

func (d *DefaultGitClient) SwitchBack() error {
	fmt.Println("SwitchBack")
	cmd := exec.Command("git", "switch", "-")
	cmd.Stdout = os.Stdout
	return cmd.Run()

	//w, err := d.repo.Worktree()
	//if err != nil {
	//	return err
	//}
	//
	//return w.Checkout(&git.CheckoutOptions{
	//	Branch: plumbing.ReferenceName(d.base.Name().String()),
	//	Create: false,
	//})
}

func (d *DefaultGitClient) Commit(message string) (CommitHash, error) {
	cmd := exec.Command("git", "commit", "-a", "-m", message)
	cmd.Stdout = os.Stdout
	return "", cmd.Run()

	//w, err := d.repo.Worktree()
	//if err != nil {
	//	return "", err
	//}
	//
	//if _, err := w.Add("."); err != nil {
	//	return "", err
	//}
	//
	//h, err := w.Commit(message, &git.CommitOptions{
	//	Author: &object.Signature{
	//		Name: "orb-update",
	//	},
	//})
	//if err != nil {
	//	return "", err
	//}
	//
	//return CommitHash(h.String()), nil
}

func (d *DefaultGitClient) Push(ctx context.Context, branch string) error {
	cmd := exec.Command("git", "push", "origin", branch)
	cmd.Stdout = os.Stdout
	return cmd.Run()

	//bs, err := d.repo.Branches()
	//if err != nil {
	//	return err
	//}
	//
	//for {
	//	ref, err := bs.Next()
	//	if err != nil {
	//		break
	//	}
	//
	//	fmt.Println(ref.String())
	//}
	//
	//return d.repo.PushContext(ctx, &git.PushOptions{
	//	RemoteName: git.DefaultRemoteName,
	//	RefSpecs:   []config.RefSpec{"refs/heads/*:refs/remotes/origin/*"},
	//	Auth: &http.BasicAuth{
	//		Username: d.d.Configuration().GithubUsername(),
	//		Password: d.d.Configuration().GithubToken(),
	//	},
	//	Progress: os.Stdout,
	//})
}
