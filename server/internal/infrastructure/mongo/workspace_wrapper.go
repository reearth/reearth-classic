package mongo

import (
	"context"

	"github.com/reearth/reearthx/account/accountinfrastructure/accountmongo"
	"github.com/reearth/reearthx/account/accountusecase/accountrepo"
	"github.com/reearth/reearthx/mongox"
)

// WorkspaceWrapper wraps accountmongo.Workspace to prevent index dropping
type WorkspaceWrapper struct {
	*accountmongo.Workspace
	client *mongox.Collection
}

// NewWorkspaceWrapper creates a workspace wrapper that doesn't drop indexes
func NewWorkspaceWrapper(client *mongox.Client) accountrepo.Workspace {
	// Create the original workspace
	original := accountmongo.NewWorkspace(client)

	return &WorkspaceWrapper{
		Workspace: original.(*accountmongo.Workspace),
		client:    client.WithCollection("workspace"),
	}
}

// Init override to prevent index dropping - just create indexes without dropping existing ones
func (w *WorkspaceWrapper) Init() error {
	ctx := context.Background()

	// Only create the "id" unique index that the original workspace expects
	// without calling the original Init() which would drop existing indexes
	return createIndexesOnly(ctx, w.client, nil, []string{"id"})
}
