package mongo

import (
	"context"
	"strings"

	"github.com/reearth/reearth/server/internal/infrastructure/adapter"
	"github.com/reearth/reearth/server/internal/infrastructure/memory"
	"github.com/reearth/reearth/server/internal/infrastructure/mongo/migration"
	"github.com/reearth/reearth/server/internal/usecase/repo"
	"github.com/reearth/reearth/server/pkg/id"
	"github.com/reearth/reearth/server/pkg/plugin"
	"github.com/reearth/reearth/server/pkg/plugin/manifest"
	"github.com/reearth/reearth/server/pkg/property"
	"github.com/reearth/reearth/server/pkg/scene"
	"github.com/reearth/reearthx/account/accountdomain/user"
	"github.com/reearth/reearthx/account/accountinfrastructure/accountmongo"
	"github.com/reearth/reearthx/account/accountusecase/accountrepo"
	"github.com/reearth/reearthx/authserver"
	"github.com/reearth/reearthx/log"
	"github.com/reearth/reearthx/mongox"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(ctx context.Context, db *mongo.Database, account *accountrepo.Container, useTransaction bool) (*repo.Container, error) {
	lock, err := NewLock(db.Collection("locks"))
	if err != nil {
		return nil, err
	}

	client := mongox.NewClientWithDatabase(db)
	if useTransaction {
		client = client.WithTransaction()
	}

	c := &repo.Container{
		Asset:          NewAsset(client),
		AuthRequest:    authserver.NewMongo(client.WithCollection("authRequest")),
		Config:         NewConfig(db.Collection("config"), lock),
		DatasetSchema:  NewDatasetSchema(client),
		Dataset:        NewDataset(client),
		Layer:          NewLayer(client),
		NLSLayer:       NewNLSLayer(client),
		Style:          NewStyle(client),
		Plugin:         NewPlugin(client),
		Project:        NewProject(client),
		PropertySchema: NewPropertySchema(client),
		Property:       NewProperty(client),
		Scene:          NewScene(client),
		Tag:            NewTag(client),
		SceneLock:      NewSceneLock(client),
		Permittable:    accountmongo.NewPermittable(client), // TODO: Delete this once the permission check migration is complete.
		Policy:         NewPolicy(client),
		Role:           accountmongo.NewRole(client), // TODO: Delete this once the permission check migration is complete.
		Storytelling:   NewStorytelling(client),
		Lock:           lock,
		Transaction:    client.Transaction(),
		Workspace:      NewWorkspaceWrapper(client),
		User:           account.User,
	}

	// init
	if err := Init(c); err != nil {
		return nil, err
	}

	// migration
	if err := migration.Do(ctx, client, c.Config); err != nil {
		return nil, err
	}

	return c, nil
}

func NewWithExtensions(ctx context.Context, db *mongo.Database, account *accountrepo.Container, useTransaction bool, src []string) (*repo.Container, error) {
	c, err := New(ctx, db, account, useTransaction)
	if err != nil {
		return nil, err
	}
	if len(src) == 0 {
		return c, nil
	}

	ms, err := manifest.ParseFromUrlList(ctx, src)
	if err != nil {
		return nil, err
	}

	ids := lo.Map(ms, func(m *manifest.Manifest, _ int) id.PluginID {
		return m.Plugin.ID()
	})

	plugins := lo.Map(ms, func(m *manifest.Manifest, _ int) *plugin.Plugin {
		return m.Plugin
	})

	propertySchemas := lo.FlatMap(ms, func(m *manifest.Manifest, _ int) []*property.Schema {
		return m.ExtensionSchema
	})

	c.Extensions = ids
	c.Plugin = adapter.NewPlugin(
		[]repo.Plugin{memory.NewPluginWith(plugins...), c.Plugin},
		c.Plugin,
	)
	c.PropertySchema = adapter.NewPropertySchema(
		[]repo.PropertySchema{memory.NewPropertySchemaWith(propertySchemas...), c.PropertySchema},
		c.PropertySchema,
	)
	return c, nil
}

func Init(r *repo.Container) error {
	if r == nil {
		return nil
	}

	ctx := context.Background()
	return util.Try(
		func() error { return r.Asset.(*Asset).Init(ctx) },
		func() error { return r.AuthRequest.(*authserver.Mongo).Init(ctx) },
		func() error { return r.Dataset.(*Dataset).Init(ctx) },
		func() error { return r.DatasetSchema.(*DatasetSchema).Init(ctx) },
		func() error { return r.Layer.(*Layer).Init(ctx) },
		func() error { return r.Permittable.(*accountmongo.Permittable).Init(ctx) }, // TODO: Delete this once the permission check migration is complete.
		func() error { return r.Plugin.(*Plugin).Init(ctx) },
		func() error { return r.Policy.(*Policy).Init(ctx) },
		func() error { return r.Project.(*Project).Init(ctx) },
		func() error { return r.Property.(*Property).Init(ctx) },
		func() error { return r.PropertySchema.(*PropertySchema).Init(ctx) },
		func() error { return r.Role.(*accountmongo.Role).Init(ctx) }, // TODO: Delete this once the permission check migration is complete.
		func() error { return r.Scene.(*Scene).Init(ctx) },
		func() error { return r.Tag.(*Tag).Init(ctx) },
		func() error { return r.User.(*accountmongo.User).Init() },
		func() error { return r.Workspace.(*WorkspaceWrapper).Init() },
	)
}

func applyWorkspaceFilter(filter interface{}, ids user.WorkspaceIDList) interface{} {
	if ids == nil {
		return filter
	}
	return mongox.And(filter, "",
		bson.M{"$or": []bson.M{
			{"workspace": bson.M{"$in": ids.Strings()}},
			{"team": bson.M{"$in": ids.Strings()}},
		},
		},
	)
}

func applySceneFilter(filter interface{}, ids scene.IDList) interface{} {
	if ids == nil {
		return filter
	}
	return mongox.And(filter, "scene", bson.M{"$in": ids.Strings()})
}

func applyOptionalSceneFilter(filter interface{}, ids scene.IDList) interface{} {
	if ids == nil {
		return filter
	}
	return mongox.And(filter, "", bson.M{"$or": []bson.M{
		{"scene": bson.M{"$in": ids.Strings()}},
		{"scene": nil},
		{"scene": ""},
	}})
}

func createIndexes(ctx context.Context, c *mongox.ClientCollection, keys, uniqueKeys []string) error {
	return createIndexesOnly(ctx, c, keys, uniqueKeys)
}

// createIndexesOnly creates indexes without dropping any existing ones
func createIndexesOnly(ctx context.Context, c *mongox.ClientCollection, keys, uniqueKeys []string) error {
	coll := c.Client()

	// Create regular indexes
	for _, key := range keys {
		indexKeys := bson.D{}
		for _, k := range strings.Split(key, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				indexKeys = append(indexKeys, bson.E{Key: k, Value: 1})
			}
		}
		if len(indexKeys) > 0 {
			indexModel := mongo.IndexModel{
				Keys: indexKeys,
			}
			if _, err := coll.Indexes().CreateOne(ctx, indexModel); err != nil {
				// Ignore error if index already exists
				if !strings.Contains(err.Error(), "already exists") && !strings.Contains(err.Error(), "IndexKeySpecsConflict") {
					log.Errorfc(ctx, "mongo: %s: failed to create index %v: %v\n", c.Client().Name(), indexKeys, err)
				}
			}
		}
	}

	// Create unique indexes
	for _, key := range uniqueKeys {
		indexKeys := bson.D{}
		for _, k := range strings.Split(key, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				indexKeys = append(indexKeys, bson.E{Key: k, Value: 1})
			}
		}
		if len(indexKeys) > 0 {
			indexModel := mongo.IndexModel{
				Keys:    indexKeys,
				Options: options.Index().SetUnique(true),
			}
			if _, err := coll.Indexes().CreateOne(ctx, indexModel); err != nil {
				// Ignore error if index already exists
				if !strings.Contains(err.Error(), "already exists") && !strings.Contains(err.Error(), "IndexKeySpecsConflict") {
					log.Errorfc(ctx, "mongo: %s: failed to create unique index %v: %v\n", c.Client().Name(), indexKeys, err)
				}
			}
		}
	}

	log.Infofc(ctx, "mongo: %s: ensured indexes exist (without dropping any)\n", c.Client().Name())
	return nil
}
