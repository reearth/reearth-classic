package mongo

import (
	"context"
	"testing"

	"github.com/reearth/reearthx/account/accountdomain"
	"github.com/reearth/reearthx/account/accountdomain/user"
	"github.com/reearth/reearthx/mongox"
	"github.com/reearth/reearthx/mongox/mongotest"
	"github.com/reearth/reearthx/rerror"
	"github.com/stretchr/testify/assert"
)

func TestUserWrapper_FindByEmailCaseInsensitive(t *testing.T) {
	c := mongotest.Connect(t)(t)
	ctx := context.Background()

	r := NewUserWrapper(mongox.NewClientWithDatabase(c))
	u1 := user.New().NewID().Name("dave").Email("Dave.Upper@Example.com").Workspace(accountdomain.NewWorkspaceID()).MustBuild()
	u2 := user.New().NewID().Name("ab").Email("ab@example.com").Workspace(accountdomain.NewWorkspaceID()).MustBuild()
	u3 := user.New().NewID().Name("axb").Email("axb@example.com").Workspace(accountdomain.NewWorkspaceID()).MustBuild()
	for _, u := range []*user.User{u1, u2, u3} {
		assert.NoError(t, r.Save(ctx, u))
	}

	f := r.(*UserWrapper)

	tests := []struct {
		name    string
		email   string
		want    *user.User
		wantErr error
	}{
		{"lowercase input, mixed-case stored", "dave.upper@example.com", u1, nil},
		{"uppercase input, mixed-case stored", "DAVE.UPPER@EXAMPLE.COM", u1, nil},
		{"exact case", "Dave.Upper@Example.com", u1, nil},
		{"anchored at start", "b@example.com", nil, rerror.ErrNotFound},
		{"anchored at end", "ab@example.co", nil, rerror.ErrNotFound},
		{"partial", "dave.upper@", nil, rerror.ErrNotFound},
		{"regex metacharacters are escaped", "a.b@example.com", nil, rerror.ErrNotFound},
		{"regex wildcard is escaped", ".*", nil, rerror.ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.FindByEmailCaseInsensitive(ctx, tt.email)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.ID(), got.ID())
			assert.Equal(t, tt.want.Email(), got.Email())
		})
	}

	got, err := f.FindByEmailCaseInsensitive(ctx, "")
	assert.NoError(t, err)
	assert.Nil(t, got)
}
