package pullrequest

import (
	"context"

	"github.com/sawadashota/orb-update/internal/configfile"
)

// Creator of Pull Request
type Creator interface {
	AlreadyCreated(ctx context.Context, branch string) (bool, error)
	Create(ctx context.Context, update *configfile.Update, message, baseBranch string) error
}
