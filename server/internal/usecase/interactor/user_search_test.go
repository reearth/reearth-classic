package interactor

import (
	"context"
	"testing"

	"github.com/reearth/reearth/server/internal/usecase"
	"github.com/reearth/reearth/server/internal/usecase/interfaces"
	"github.com/reearth/reearthx/account/accountdomain"
	"github.com/reearth/reearthx/account/accountdomain/user"
	"github.com/reearth/reearthx/account/accountinfrastructure/accountmemory"
	"github.com/reearth/reearthx/account/accountusecase"
	"github.com/reearth/reearthx/account/accountusecase/accountrepo"
	"github.com/stretchr/testify/assert"
)

func TestUserSearch_FindByEmail(t *testing.T) {
	ctx := context.Background()
	u1 := user.New().NewID().Name("alice").Email("alice@example.com").Workspace(accountdomain.NewWorkspaceID()).MustBuild()
	// a user whose name looks like another user's email must not be matched by name
	u2 := user.New().NewID().Name("bob@example.com").Email("robert@example.com").Workspace(accountdomain.NewWorkspaceID()).MustBuild()
	// a user on an additional (tenant) repo
	u3 := user.New().NewID().Name("carol").Email("carol@example.com").Workspace(accountdomain.NewWorkspaceID()).MustBuild()

	s := NewUserSearch(&accountrepo.Container{
		User:  accountmemory.NewUserWith(u1, u2),
		Users: []accountrepo.User{accountmemory.NewUserWith(u3)},
	})

	uid := u1.ID()
	op := &usecase.Operator{AcOperator: &accountusecase.Operator{User: &uid}}

	tests := []struct {
		name  string
		email string
		want  *user.User
	}{
		{"exact email", "alice@example.com", u1},
		{"trimmed and case-insensitive", "  ALICE@Example.com ", u1},
		{"tenant repo", "carol@example.com", u3},
		{"no name match", "bob@example.com", nil},
		{"no partial email match", "alice@", nil},
		{"no partial domain match", "example.com", nil},
		{"no plain name lookup", "alice", nil},
		{"not found", "nobody@example.com", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.FindByEmail(ctx, tt.email, op)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}

	t.Run("requires an authenticated operator", func(t *testing.T) {
		got, err := s.FindByEmail(ctx, "alice@example.com", nil)
		assert.ErrorIs(t, err, interfaces.ErrOperationDenied)
		assert.Nil(t, got)

		got, err = s.FindByEmail(ctx, "alice@example.com", &usecase.Operator{AcOperator: &accountusecase.Operator{}})
		assert.ErrorIs(t, err, interfaces.ErrOperationDenied)
		assert.Nil(t, got)
	})
}
