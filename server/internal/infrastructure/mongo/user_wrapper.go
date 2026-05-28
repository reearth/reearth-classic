package mongo

import (
	"context"

	"github.com/reearth/reearthx/account/accountinfrastructure/accountmongo"
	"github.com/reearth/reearthx/account/accountusecase/accountrepo"
	"github.com/reearth/reearthx/mongox"
)

// UserWrapper wraps accountmongo.User to prevent index dropping
type UserWrapper struct {
	*accountmongo.User
	client *mongox.Collection
}

func NewUserWrapper(client *mongox.Client) accountrepo.User {
	original := accountmongo.NewUser(client)

	return &UserWrapper{
		User:   original.(*accountmongo.User),
		client: client.WithCollection("user"),
	}
}

func (u *UserWrapper) Init() error {
	ctx := context.Background()

	return createIndexesOnly(ctx, u.client, nil, []string{"id"})
}
