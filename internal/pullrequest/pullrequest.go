package pullrequest

import (
	"context"

	"github.com/sawadashota/orb-update/internal/orb"
)

// Creator of Pull Request
type Creator interface {
	AlreadyCreated(ctx context.Context, branch string) (bool, error)
	Create(ctx context.Context, diff *orb.Difference, message, baseBranch string) error
}
