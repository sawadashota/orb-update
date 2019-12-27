package pullrequest

import (
	"context"

	"github.com/sawadashota/orb-update/internal/orb"
)

// Creator of Pull Request
type Creator interface {
	AlreadyCreated(ctx context.Context, branch string) (bool, error)
	Create(ctx context.Context, update *orb.Update, message, baseBranch string) error
}
