package gql

import (
	"context"
	"errors"
	"strings"

	"github.com/reearth/reearth/server/internal/adapter/gql/gqldataloader"
	"github.com/reearth/reearth/server/internal/adapter/gql/gqlmodel"
	"github.com/reearth/reearthx/account/accountdomain"
	"github.com/reearth/reearthx/account/accountusecase/accountinterfaces"
	"github.com/reearth/reearthx/rerror"
	"github.com/reearth/reearthx/util"
)

type UserLoader struct {
	usecase accountinterfaces.User
}

func NewUserLoader(usecase accountinterfaces.User) *UserLoader {
	return &UserLoader{usecase: usecase}
}

func (c *UserLoader) Fetch(ctx context.Context, ids []gqlmodel.ID) ([]*gqlmodel.User, []error) {
	uids, err := util.TryMap(ids, gqlmodel.ToID[accountdomain.User])
	if err != nil {
		return nil, []error{err}
	}

	res, err := c.usecase.FetchByID(ctx, uids)
	if err != nil {
		return nil, []error{err}
	}

	users := make([]*gqlmodel.User, 0, len(res))
	for _, u := range res {
		users = append(users, gqlmodel.ToUser(u))
	}

	return users, nil
}

// SearchUser looks up a user by exact email address only.
// Partial matches and name lookups are intentionally not supported, so the
// endpoint cannot be used to enumerate users or discover email addresses.
func (c *UserLoader) SearchUser(ctx context.Context, email string) (*gqlmodel.User, error) {
	email = strings.TrimSpace(email)
	if !strings.Contains(email, "@") {
		return nil, nil
	}

	candidates := []string{email}
	if lower := strings.ToLower(email); lower != email {
		candidates = append(candidates, lower)
	}

	for _, e := range candidates {
		// FetchByNameOrEmail matches name OR email exactly; reject name matches.
		res, err := c.usecase.FetchByNameOrEmail(ctx, e)
		if err != nil && !errors.Is(err, rerror.ErrNotFound) {
			return nil, err
		}
		if res != nil && strings.EqualFold(res.Email, email) {
			return gqlmodel.ToUserFromSimple(res), nil
		}
	}
	return nil, nil
}

// data loader

type UserDataLoader interface {
	Load(gqlmodel.ID) (*gqlmodel.User, error)
	LoadAll([]gqlmodel.ID) ([]*gqlmodel.User, []error)
}

func (c *UserLoader) DataLoader(ctx context.Context) UserDataLoader {
	return gqldataloader.NewUserLoader(gqldataloader.UserLoaderConfig{
		Wait:     dataLoaderWait,
		MaxBatch: dataLoaderMaxBatch,
		Fetch: func(keys []gqlmodel.ID) ([]*gqlmodel.User, []error) {
			return c.Fetch(ctx, keys)
		},
	})
}

func (c *UserLoader) OrdinaryDataLoader(ctx context.Context) UserDataLoader {
	return &ordinaryUserLoader{
		fetch: func(keys []gqlmodel.ID) ([]*gqlmodel.User, []error) {
			return c.Fetch(ctx, keys)
		},
	}
}

type ordinaryUserLoader struct {
	fetch func(keys []gqlmodel.ID) ([]*gqlmodel.User, []error)
}

func (l *ordinaryUserLoader) Load(key gqlmodel.ID) (*gqlmodel.User, error) {
	res, errs := l.fetch([]gqlmodel.ID{key})
	if len(errs) > 0 {
		return nil, errs[0]
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return nil, nil
}

func (l *ordinaryUserLoader) LoadAll(keys []gqlmodel.ID) ([]*gqlmodel.User, []error) {
	return l.fetch(keys)
}
