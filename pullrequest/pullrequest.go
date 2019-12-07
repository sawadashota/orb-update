package pullrequest

import "context"

type Creator interface {
	Create(ctx context.Context, message string) error
}
