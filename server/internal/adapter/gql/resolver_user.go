package gql

import (
	"context"

	"github.com/reearth/reearth/server/internal/adapter/gql/gqlmodel"
)

func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

type userResolver struct{ *Resolver }

// Email returns the user's email only to the user themselves.
// For anyone else it returns an empty string, so emails never leak through
// searchUser, node/nodes or workspace member lists.
func (r *userResolver) Email(ctx context.Context, obj *gqlmodel.User) (string, error) {
	if u := getUser(ctx); u != nil && gqlmodel.IDFrom(u.ID()) == obj.ID {
		return obj.Email, nil
	}
	return "", nil
}

func (r *Resolver) Me() MeResolver {
	return &meResolver{r}
}

type meResolver struct{ *Resolver }

func (r *meResolver) MyTeam(ctx context.Context, obj *gqlmodel.Me) (*gqlmodel.Team, error) {
	return dataloaders(ctx).Workspace.Load(obj.MyTeamID)
}

func (r *meResolver) Teams(ctx context.Context, obj *gqlmodel.Me) ([]*gqlmodel.Team, error) {
	return loaders(ctx).Workspace.FindByUser(ctx, obj.ID)
}
