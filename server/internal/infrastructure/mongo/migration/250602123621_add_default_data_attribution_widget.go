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
				widgetID := id.NewWidgetID().String()

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

				newWidget := bson.M{
					"id":        widgetID,
					"plugin":    "reearth",
					"extension": "dataattribution",
					"property":  propertyID,
					"enabled":   true,
					"extended":  false,
				}

				needsAlignInit := false
				if scene["alignsystem"] == nil {
					needsAlignInit = true
				} else if outer, ok := scene["alignsystem"].(bson.M)["outer"]; !ok || outer == nil {
					needsAlignInit = true
				} else if left, ok := outer.(bson.M)["left"]; !ok || left == nil {
					needsAlignInit = true
				} else if bottom, ok := left.(bson.M)["bottom"]; !ok || bottom == nil {
					needsAlignInit = true
				} else if _, ok := bottom.(bson.M)["widgetids"]; !ok {
					needsAlignInit = true
				}

				if needsAlignInit {
					_, err := sceneCol.Client().UpdateByID(ctx, objectID, bson.M{
						"$set": bson.M{
							"alignsystem.outer.left.bottom": bson.M{
								"widgetids":  bson.A{},
								"align":      "start",
								"padding":    nil,
								"gap":        nil,
								"centered":   false,
								"background": nil,
							},
						},
					})
					if err != nil {
						return err
					}
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
