package pullrequest

import (
	"context"

	"github.com/sawadashota/orb-update/internal/extraction"
)

// Creator of Pull Request
type Creator interface {
	AlreadyCreated(ctx context.Context, branch string) (bool, error)
	Create(ctx context.Context, update *extraction.Update, message, baseBranch string) error
}
