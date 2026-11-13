package interactor

import (
	"context"
	"strings"
	"testing"

	"github.com/reearth/reearth/server/internal/usecase"
	"github.com/reearth/reearth/server/internal/usecase/interfaces"
	"github.com/reearth/reearthx/account/accountdomain"
	"github.com/reearth/reearthx/account/accountdomain/user"
	"github.com/reearth/reearthx/account/accountinfrastructure/accountmemory"
	"github.com/reearth/reearthx/account/accountusecase"
	"github.com/reearth/reearthx/account/accountusecase/accountrepo"
	"github.com/reearth/reearthx/rerror"
	"github.com/stretchr/testify/assert"
)

func TestUserSearch_FindByEmail(t *testing.T) {
	ctx := context.Background()
	u1 := user.New().NewID().Name("alice").Email("alice@example.com").Workspace(accountdomain.NewWorkspaceID()).MustBuild()
	u2 := user.New().NewID().Name("bob@example.com").Email("robert@example.com").Workspace(accountdomain.NewWorkspaceID()).MustBuild()
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

type caseInsensitiveUserRepo struct {
	*accountmemory.User
	ciCalls int
}

func (r *caseInsensitiveUserRepo) FindByEmailCaseInsensitive(ctx context.Context, email string) (*user.User, error) {
	r.ciCalls++
	all, err := r.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, u := range all {
		if strings.EqualFold(u.Email(), email) {
			return u, nil
		}
	}
	return nil, rerror.ErrNotFound
}

func TestUserSearch_FindByEmail_CaseInsensitiveFallback(t *testing.T) {
	ctx := context.Background()
	stored := user.New().NewID().Name("dave").Email("Dave.Upper@Example.com").Workspace(accountdomain.NewWorkspaceID()).MustBuild()
	lower := user.New().NewID().Name("erin").Email("erin@example.com").Workspace(accountdomain.NewWorkspaceID()).MustBuild()
	tenantStored := user.New().NewID().Name("frank").Email("Frank@Tenant.com").Workspace(accountdomain.NewWorkspaceID()).MustBuild()

	primary := &caseInsensitiveUserRepo{User: accountmemory.NewUserWith(stored, lower)}
	tenant := &caseInsensitiveUserRepo{User: accountmemory.NewUserWith(tenantStored)}
	s := NewUserSearch(&accountrepo.Container{
		User:  primary,
		Users: []accountrepo.User{tenant},
	})

	uid := lower.ID()
	op := &usecase.Operator{AcOperator: &accountusecase.Operator{User: &uid}}

	tests := []struct {
		name   string
		email  string
		want   *user.User
		usesCI bool
	}{
		{"mixed-case stored, lowercase typed", "dave.upper@example.com", stored, true},
		{"mixed-case stored, uppercase typed", "DAVE.UPPER@EXAMPLE.COM", stored, true},
		{"mixed-case stored, exact case typed", "Dave.Upper@Example.com", stored, false},
		{"lowercase stored, uppercase typed", "ERIN@EXAMPLE.COM", lower, false},
		{"mixed-case stored in tenant repo", "frank@tenant.com", tenantStored, true},
		{"partial is still not matched", "dave.upper@", nil, true},
		{"name is still not matched", "dave", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primary.ciCalls, tenant.ciCalls = 0, 0
			got, err := s.FindByEmail(ctx, tt.email, op)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.usesCI, primary.ciCalls+tenant.ciCalls > 0)
		})
	}
}
