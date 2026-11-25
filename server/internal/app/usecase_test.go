package app

import (
	"context"
	"testing"

	"github.com/reearth/reearth/server/internal/infrastructure/memory"
	"github.com/reearth/reearth/server/pkg/project"
	"github.com/reearth/reearthx/rerror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReposForOperator_AnonymousDenyAll is the SEC-01 regression guard: an
// anonymous request (nil operator) must get a deny-all filter so workspace/scene
// scoped reads return not-found. Without it, plugin(id) -> scene -> project ->
// basicAuthPassword leaks every project's credentials to anonymous callers.
func TestReposForOperator_AnonymousDenyAll(t *testing.T) {
	ctx := context.Background()

	prj := project.New().NewID().Workspace(project.NewWorkspaceID()).MustBuild()
	base := memory.New() // unfiltered container
	require.NoError(t, base.Project.Save(ctx, prj))

	// sanity: the unfiltered repo returns the project
	got, err := base.Project.FindByID(ctx, prj.ID())
	require.NoError(t, err)
	require.NotNil(t, got)

	// anonymous (nil operator) -> deny-all -> not found
	anon := reposForOperator(base, nil)
	got, err = anon.Project.FindByID(ctx, prj.ID())
	assert.Nil(t, got, "anonymous callers must not be able to read any project")
	assert.ErrorIs(t, err, rerror.ErrNotFound)
}
