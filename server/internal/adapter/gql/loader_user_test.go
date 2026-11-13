package gql

import (
	"context"
	"testing"

	"github.com/reearth/reearth/server/internal/adapter/gql/gqlmodel"
	"github.com/reearth/reearth/server/internal/usecase"
	"github.com/reearth/reearthx/account/accountdomain"
	"github.com/reearth/reearthx/account/accountdomain/user"
	"github.com/stretchr/testify/assert"
)

type fakeUserSearch struct {
	u *user.User
}

func (f *fakeUserSearch) FindByEmail(_ context.Context, _ string, _ *usecase.Operator) (*user.User, error) {
	return f.u, nil
}

func TestUserLoader_SearchUser(t *testing.T) {
	ctx := context.Background()
	u := user.New().NewID().Name("bob").Email("bob@example.com").Workspace(accountdomain.NewWorkspaceID()).MustBuild().WithHost("tenant.example.com")

	got, err := NewUserLoader(nil, &fakeUserSearch{u: u}).SearchUser(ctx, "bob@example.com")
	assert.NoError(t, err)
	assert.Equal(t, gqlmodel.IDFrom(u.ID()), got.ID)
	assert.Equal(t, "bob", got.Name)
	assert.Nil(t, got.Host)

	got, err = NewUserLoader(nil, &fakeUserSearch{}).SearchUser(ctx, "nobody@example.com")
	assert.NoError(t, err)
	assert.Nil(t, got)
}
