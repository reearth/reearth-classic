package interactor

import (
	"context"
	"errors"
	"strings"

	"github.com/reearth/reearth/server/internal/usecase"
	"github.com/reearth/reearth/server/internal/usecase/interfaces"
	"github.com/reearth/reearthx/account/accountdomain/user"
	"github.com/reearth/reearthx/account/accountusecase/accountrepo"
	"github.com/reearth/reearthx/rerror"
)

type UserSearch struct {
	repos []accountrepo.User
}

func NewUserSearch(ar *accountrepo.Container) interfaces.UserSearch {
	var repos []accountrepo.User
	if ar != nil {
		if ar.User != nil {
			repos = append(repos, ar.User)
		}
		repos = append(repos, ar.Users...)
	}
	return &UserSearch{repos: repos}
}

func (i *UserSearch) FindByEmail(ctx context.Context, email string, op *usecase.Operator) (*user.User, error) {
	if op == nil || op.AcOperator == nil || op.AcOperator.User == nil {
		return nil, interfaces.ErrOperationDenied
	}

	email = strings.TrimSpace(email)
	if !strings.Contains(email, "@") {
		return nil, nil
	}

	candidates := []string{email}
	if lower := strings.ToLower(email); lower != email {
		candidates = append(candidates, lower)
	}

	for _, e := range candidates {
		for _, r := range i.repos {
			u, err := r.FindByEmail(ctx, e)
			if errors.Is(err, rerror.ErrNotFound) {
				continue
			}
			if err != nil {
				return nil, err
			}
			if u != nil {
				return u, nil
			}
		}
	}
	return nil, nil
}
