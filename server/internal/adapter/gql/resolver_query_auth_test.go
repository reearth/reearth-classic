package gql

import (
	"context"
	"testing"

	"github.com/reearth/reearth/server/internal/adapter/gql/gqlmodel"
	"github.com/stretchr/testify/assert"
)

// SEC-01: these top-level query resolvers must reject anonymous callers.
// requireUser is the first statement in each, so an operator-less context is
// refused with ErrUnauthorized before any loader is touched (hence no repo or
// Mongo is needed here). Without the guard, an anonymous caller could walk
// layer -> scene -> workspace -> projects and read basicAuthUsername/Password.
func TestQueryResolver_RequireUser_Anonymous(t *testing.T) {
	r := &queryResolver{&Resolver{}}
	ctx := context.Background() // no user attached => anonymous

	t.Run("layer", func(t *testing.T) {
		got, err := r.Layer(ctx, gqlmodel.ID("x"))
		assert.Same(t, ErrUnauthorized, err)
		assert.Nil(t, got)
	})
	t.Run("scene", func(t *testing.T) {
		got, err := r.Scene(ctx, gqlmodel.ID("x"))
		assert.Same(t, ErrUnauthorized, err)
		assert.Nil(t, got)
	})
	t.Run("projects", func(t *testing.T) {
		got, err := r.Projects(ctx, gqlmodel.ID("x"), nil, nil, nil, nil, nil)
		assert.Same(t, ErrUnauthorized, err)
		assert.Nil(t, got)
	})
	t.Run("datasetSchemas", func(t *testing.T) {
		got, err := r.DatasetSchemas(ctx, gqlmodel.ID("x"), nil, nil, nil, nil)
		assert.Same(t, ErrUnauthorized, err)
		assert.Nil(t, got)
	})
	t.Run("datasets", func(t *testing.T) {
		got, err := r.Datasets(ctx, gqlmodel.ID("x"), nil, nil, nil, nil)
		assert.Same(t, ErrUnauthorized, err)
		assert.Nil(t, got)
	})
}
