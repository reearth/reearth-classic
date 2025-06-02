package migration

import (
	"context"

	"github.com/reearth/reearth/server/pkg/id"
	"github.com/reearth/reearthx/log"
	"github.com/reearth/reearthx/mongox"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddDefaultDataAttributionWidget(ctx context.Context, c DBClient) error {
	sceneCol := c.WithCollection("scene")
	propertyCol := c.WithCollection("property")

	filter := bson.M{
		"widgets": bson.M{
			"$not": bson.M{
				"$elemMatch": bson.M{
					"plugin":    "reearth",
					"extension": "dataattribution",
				},
			},
		},
	}

	return sceneCol.Find(ctx, filter, &mongox.BatchConsumer{
		Size: 1000,
		Callback: func(rows []bson.Raw) error {
			log.Infofc(ctx, "migration: AddDefaultDataAttributionWidget: matched scenes: %d", len(rows))

			for _, row := range rows {

				var scene bson.M
				if err := bson.Unmarshal(row, &scene); err != nil {
					return err
				}

				sceneID, ok := scene["id"].(string)
				if !ok {
					continue
				}
				objectID, ok := scene["_id"].(primitive.ObjectID)
				if !ok {
					continue
				}

				propertyID := id.NewPropertyID().String()
				newProperty := bson.M{
					"id":           propertyID,
					"scene":        sceneID,
					"schemaplugin": "reearth",
					"schemaname":   "dataattribution",
					"items":        bson.A{},
				}

				if _, err := propertyCol.Client().InsertOne(ctx, newProperty); err != nil {
					return err
				}

				widgetID := id.NewWidgetID().String()
				newWidget := bson.M{
					"id":        widgetID,
					"plugin":    "reearth",
					"extension": "dataattribution",
					"property":  propertyID,
					"enabled":   true,
					"extended":  false,
				}

				update := bson.M{
					"$push": bson.M{
						"widgets": newWidget,
						"alignsystem.outer.left.bottom.widgetids": widgetID,
					},
				}
				if _, err := sceneCol.Client().UpdateByID(ctx, objectID, update); err != nil {
					return err
				}

			}
			return nil
		},
	})
}
