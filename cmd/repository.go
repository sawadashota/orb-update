package cmd

import (
	"strings"

	"github.com/pkg/errors"
)

type repository struct {
	owner string
	name  string
}

func parseRepository(repo string) (*repository, error) {
	s := strings.Split(repo, "/")
	if len(s) != 2 {
		return nil, errors.Errorf("incorrect GitHub repository format: %s", repo)
	}
	return &repository{
		owner: s[0],
		name:  s[1],
	}, nil
}
