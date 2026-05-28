package migration

import (
	"context"
	"testing"

	"github.com/reearth/reearthx/mongox"
	"github.com/reearth/reearthx/mongox/mongotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFixArcgisTerrain(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	// Realistic property document: a tiles grouplist item followed by a terrain group item,
	// matching the structure seen in production (confirmed via mongodoc.PropertyDocument).
	// Terrain (group) stores fields at items[].fields[], not items[].groups[].fields[].
	docArcgis := bson.M{
		"_id": primitive.NewObjectID(),
		"id":  "test-arcgis-prop-001",
		"items": bson.A{
			bson.M{
				"type":        "grouplist",
				"schemagroup": "tiles",
				"groups": bson.A{
					bson.M{
						"type":   "group",
						"fields": bson.A{bson.M{"field": "tile_type", "type": "STRING", "value": "open_street_map"}},
					},
				},
			},
			bson.M{
				"type":        "group",
				"schemagroup": "terrain",
				"fields": bson.A{
					bson.M{"field": "terrain", "type": "BOOL", "value": true},
					bson.M{"field": "terrainType", "type": "STRING", "value": "arcgis"},
				},
			},
		},
	}

	// A doc with cesium terrain — must not be touched.
	docCesium := bson.M{
		"_id": primitive.NewObjectID(),
		"id":  "test-cesium-prop-001",
		"items": bson.A{
			bson.M{
				"type":        "group",
				"schemagroup": "terrain",
				"fields": bson.A{
					bson.M{"field": "terrain", "type": "BOOL", "value": true},
					bson.M{"field": "terrainType", "type": "STRING", "value": "cesium"},
				},
			},
		},
	}

	// A doc with no terrain at all — must not be touched.
	docNoTerrain := bson.M{
		"_id": primitive.NewObjectID(),
		"id":  "test-no-terrain-prop-001",
		"items": bson.A{
			bson.M{
				"type":        "grouplist",
				"schemagroup": "tiles",
				"groups": bson.A{
					bson.M{
						"type":   "group",
						"fields": bson.A{bson.M{"field": "tile_type", "type": "STRING", "value": "open_street_map"}},
					},
				},
			},
		},
	}

	_, err := db.Collection("property").InsertMany(ctx, []any{docArcgis, docCesium, docNoTerrain})
	require.NoError(t, err)

	err = FixArcgisTerrain(ctx, c)
	require.NoError(t, err)

	// arcgis → reearth_terrain, terrain: true preserved
	var result1 bson.M
	require.NoError(t, db.Collection("property").FindOne(ctx, bson.M{"_id": docArcgis["_id"]}).Decode(&result1))
	terrainItem := result1["items"].(primitive.A)[1].(bson.M)
	fields1 := terrainItem["fields"].(primitive.A)
	assert.Equal(t, true, fields1[0].(bson.M)["value"], "terrain: true should be preserved")
	assert.Equal(t, "reearth_terrain", fields1[1].(bson.M)["value"], "arcgis should become reearth_terrain")

	// cesium unchanged
	var result2 bson.M
	require.NoError(t, db.Collection("property").FindOne(ctx, bson.M{"_id": docCesium["_id"]}).Decode(&result2))
	fields2 := result2["items"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Equal(t, "cesium", fields2[1].(bson.M)["value"], "cesium terrain should be unchanged")

	// tile-only doc unchanged
	var result3 bson.M
	require.NoError(t, db.Collection("property").FindOne(ctx, bson.M{"_id": docNoTerrain["_id"]}).Decode(&result3))
	tileGroups := result3["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)
	assert.Equal(t, "open_street_map", tileGroups[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"])
}
