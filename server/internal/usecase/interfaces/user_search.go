package interfaces

import (
	"context"

	"github.com/reearth/reearth/server/internal/usecase"
	"github.com/reearth/reearthx/account/accountdomain/user"
)

type UserSearch interface {
	FindByEmail(ctx context.Context, email string, op *usecase.Operator) (*user.User, error)
}
