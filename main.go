package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4/config"
	"os"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"gopkg.in/src-d/go-git.v4/plumbing"

	"gopkg.in/src-d/go-git.v4"
)

func main() {

	_, err := run()
	if err != nil {
		panic(err)
	}
	//if err := cmd.RootCmd().Execute(); err != nil {
	//	_, _ = fmt.Fprintln(os.Stderr, err)
	//	os.Exit(1)
	//}
}

func run() ([]byte, error) {
	// open
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(pwd)
	if err != nil {
		return nil, err
	}

	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	// checkout
	branch := "test-" + time.Now().String()
	ref := plumbing.ReferenceName(branch)
	err = repo.CreateBranch(&config.Branch{
		Name:   branch,
		Merge:  "exp",
	})
	if err != nil {
		return nil, err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: ref,
		//Create: true,
	})
	if err != nil {
		return nil, err
	}

	br, err := repo.Branch(branch)
	if err != nil {
		return nil, err
	}
	fmt.Println(br.Name)

	// create file
	f, err := os.Create(".test.txt")
	if err != nil {
		return nil, err
	}
	f.Write([]byte("hello\n"))
	f.Close()

	// commit
	_, err = w.Add(".")
	if err != nil {
		return nil, err
	}

	hash, _ := w.Commit("update version", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Shota Sawada",
			Email: "xiootas@gmail.com",
			When:  time.Now(),
		},
	})
	err = repo.Storer.SetReference(plumbing.NewReferenceFromStrings(branch, hash.String()))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
