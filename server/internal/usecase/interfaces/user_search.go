package interfaces

import (
	"context"

	"github.com/reearth/reearth/server/internal/usecase"
	"github.com/reearth/reearthx/account/accountdomain/user"
)

type UserSearch interface {
	// FindByEmail returns the user whose email exactly matches the given address,
	// or nil if there is none. Requires an authenticated operator.
	FindByEmail(ctx context.Context, email string, op *usecase.Operator) (*user.User, error)
}
