package mongo

import (
	"context"

	"github.com/reearth/reearthx/account/accountinfrastructure/accountmongo"
	"github.com/reearth/reearthx/account/accountusecase/accountrepo"
	"github.com/reearth/reearthx/mongox"
)

// RoleWrapper wraps accountmongo.Role to prevent index dropping
type RoleWrapper struct {
	*accountmongo.Role
	client *mongox.Collection
}

func NewRoleWrapper(client *mongox.Client) accountrepo.Role {
	original := accountmongo.NewRole(client)

	return &RoleWrapper{
		Role:   original.(*accountmongo.Role),
		client: client.WithCollection("role"),
	}
}

func (r *RoleWrapper) Init(ctx context.Context) error {
	return createIndexesOnly(ctx, r.client, nil, []string{"id"})
}
