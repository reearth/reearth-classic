package migration

import (
	"context"
	"testing"

	"github.com/reearth/reearth/server/pkg/id"
	"github.com/reearth/reearthx/mongox"
	"github.com/reearth/reearthx/mongox/mongotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAddDefaultDataAttributionWidget(t *testing.T) {
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)
	ctx := context.Background()

	sceneID := "scene-001"
	sceneDoc := bson.M{
		"id":      sceneID,
		"widgets": bson.A{},
		"alignsystem": bson.M{
			"inner": nil,
			"outer": bson.M{
				"left": bson.M{
					"top":    nil,
					"middle": nil,
					"bottom": nil,
				},
				"center": nil,
				"right":  nil,
			},
		},
	}

	_, err := db.Collection("scene").InsertOne(ctx, sceneDoc)
	require.NoError(t, err)

	err = AddDefaultDataAttributionWidget(ctx, c)
	require.NoError(t, err)

	var updatedScene bson.M
	err = db.Collection("scene").FindOne(ctx, bson.M{"id": sceneID}).Decode(&updatedScene)
	require.NoError(t, err)

	widgets := updatedScene["widgets"].(primitive.A)
	require.Len(t, widgets, 1)
	widget := widgets[0].(bson.M)
	assert.Equal(t, "reearth", widget["plugin"])
	assert.Equal(t, "dataattribution", widget["extension"])
	assert.NotEmpty(t, widget["property"])

	align := updatedScene["alignsystem"].(bson.M)
	left := align["outer"].(bson.M)["left"].(bson.M)
	bottom := left["bottom"].(bson.M)
	ids := bottom["widgetids"].(primitive.A)
	require.Len(t, ids, 1)
	assert.Equal(t, widget["id"], ids[0])

	count, err := db.Collection("property").CountDocuments(ctx, bson.M{
		"scene":        sceneID,
		"schemaplugin": "reearth",
		"schemaname":   "dataattribution",
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestAddDefaultDataAttributionWidget_WithExistingWidget(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	sceneID := "scene-002"
	existingWidgetID := id.NewWidgetID().String()

	sceneDoc := bson.M{
		"id": sceneID,
		"widgets": bson.A{
			bson.M{
				"id":        existingWidgetID,
				"plugin":    "reearth",
				"extension": "someotherextension",
				"property":  "some-property-id",
				"enabled":   true,
				"extended":  false,
			},
		},
		"alignsystem": bson.M{
			"inner": nil,
			"outer": bson.M{
				"left": bson.M{
					"top":    nil,
					"middle": nil,
					"bottom": bson.M{
						"widgetids":  bson.A{existingWidgetID},
						"align":      "start",
						"padding":    nil,
						"gap":        nil,
						"centered":   false,
						"background": nil,
					},
				},
				"center": nil,
				"right":  nil,
			},
		},
	}
	_, err := db.Collection("scene").InsertOne(ctx, sceneDoc)
	require.NoError(t, err)

	err = AddDefaultDataAttributionWidget(ctx, c)
	require.NoError(t, err)

	var updatedScene bson.M
	err = db.Collection("scene").FindOne(ctx, bson.M{"id": sceneID}).Decode(&updatedScene)
	require.NoError(t, err)

	widgets := updatedScene["widgets"].(primitive.A)
	require.Len(t, widgets, 2)

	var foundDataAttribution bool
	var newWidgetID string
	for _, w := range widgets {
		widget := w.(bson.M)
		if widget["extension"] == "dataattribution" {
			foundDataAttribution = true
			newWidgetID = widget["id"].(string)
			break
		}
	}
	assert.True(t, foundDataAttribution, "dataattribution widget should be added")

	bottom := updatedScene["alignsystem"].(bson.M)["outer"].(bson.M)["left"].(bson.M)["bottom"].(bson.M)
	widgetids := bottom["widgetids"].(primitive.A)
	assert.Contains(t, widgetids, existingWidgetID)
	assert.Contains(t, widgetids, newWidgetID)

	count, err := db.Collection("property").CountDocuments(ctx, bson.M{
		"scene":        sceneID,
		"schemaplugin": "reearth",
		"schemaname":   "dataattribution",
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}
