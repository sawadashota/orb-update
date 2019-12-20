package git

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Repository of VCS
type Repository struct {
	owner string
	name  string
}

// ParseRepository from owner/name format
func ParseRepository(repo string) (*Repository, error) {
	s := strings.Split(repo, "/")
	if len(s) != 2 {
		return nil, errors.Errorf("incorrect GitHub repository format: %s", repo)
	}
	return &Repository{
		owner: s[0],
		name:  s[1],
	}, nil
}

// Owner of repository
func (r *Repository) Owner() string {
	return r.owner
}

// Name of repository
func (r *Repository) Name() string {
	return r.name
}

// String of original format
func (r *Repository) String() string {
	return fmt.Sprintf("%s/%s", r.owner, r.name)
}
