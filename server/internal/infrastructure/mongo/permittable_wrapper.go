package mongo

import (
	"context"

	"github.com/reearth/reearthx/account/accountinfrastructure/accountmongo"
	"github.com/reearth/reearthx/account/accountusecase/accountrepo"
	"github.com/reearth/reearthx/mongox"
)

// PermittableWrapper wraps accountmongo.Permittable to prevent index dropping
type PermittableWrapper struct {
	*accountmongo.Permittable
	client *mongox.Collection
}

func NewPermittableWrapper(client *mongox.Client) accountrepo.Permittable {
	original := accountmongo.NewPermittable(client)

	return &PermittableWrapper{
		Permittable: original.(*accountmongo.Permittable),
		client:      client.WithCollection("permittable"),
	}
}

func (p *PermittableWrapper) Init(ctx context.Context) error {
	return createIndexesOnly(ctx, p.client, nil, []string{"id"})
}
