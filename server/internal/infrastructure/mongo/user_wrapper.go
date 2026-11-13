package mongo

import (
	"context"
	"regexp"

	"github.com/reearth/reearthx/account/accountdomain/user"
	"github.com/reearth/reearthx/account/accountinfrastructure/accountmongo"
	"github.com/reearth/reearthx/account/accountinfrastructure/accountmongo/mongodoc"
	"github.com/reearth/reearthx/account/accountusecase/accountrepo"
	"github.com/reearth/reearthx/mongox"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserWrapper wraps accountmongo.User to prevent index dropping
type UserWrapper struct {
	*accountmongo.User
	client *mongox.Collection
}

func NewUserWrapper(client *mongox.Client) accountrepo.User {
	original := accountmongo.NewUser(client)

	return &UserWrapper{
		User:   original.(*accountmongo.User),
		client: client.WithCollection("user"),
	}
}

func (u *UserWrapper) Init() error {
	ctx := context.Background()

	return createIndexesOnly(ctx, u.client, nil, []string{"id"})
}

func (u *UserWrapper) FindByEmailCaseInsensitive(ctx context.Context, email string) (*user.User, error) {
	if email == "" {
		return nil, nil
	}
	c := mongodoc.NewUserConsumer("")
	filter := bson.M{"email": bson.M{"$regex": primitive.Regex{Pattern: "^" + regexp.QuoteMeta(email) + "$", Options: "i"}}}
	if err := u.client.FindOne(ctx, filter, c); err != nil {
		return nil, err
	}
	return c.Result[0], nil
}
