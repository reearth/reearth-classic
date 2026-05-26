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

func TestUpdateTileAndTerrainProviders_TileMigrations(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	// Test case 1: default → cesium_ion with asset_id 2
	doc1 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "default",
							},
						},
					},
				},
			},
		},
	}

	// Test case 2: default_label → cesium_ion with asset_id 3
	doc2 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "default_label",
							},
						},
					},
				},
			},
		},
	}

	// Test case 3: default_road → cesium_ion with asset_id 4
	doc3 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "default_road",
							},
						},
					},
				},
			},
		},
	}

	// Test case 4: black_marble → cesium_ion with asset_id 3812
	doc4 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "black_marble",
							},
						},
					},
				},
			},
		},
	}

	// Test case 5: stamen_toner → carto_light
	doc5 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "stamen_toner",
							},
						},
					},
				},
			},
		},
	}

	// Test case 6: esri_world_topo → open_street_map
	doc6 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "esri_world_topo",
							},
						},
					},
				},
			},
		},
	}

	// Test case 7: tile type that should not be changed
	doc7 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "open_street_map",
							},
						},
					},
				},
			},
		},
	}

	// Insert test documents
	_, err := db.Collection("property").InsertMany(ctx, []any{doc1, doc2, doc3, doc4, doc5, doc6, doc7})
	require.NoError(t, err)

	// Run migration
	err = UpdateTileAndTerrainProviders(ctx, c)
	require.NoError(t, err)

	// Verify doc1: default → cesium_ion with asset_id 2
	var updatedDoc1 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc1["_id"]}).Decode(&updatedDoc1)
	require.NoError(t, err)
	fields1 := updatedDoc1["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Equal(t, "cesium_ion", fields1[0].(bson.M)["value"])
	assert.Equal(t, "cesium_ion_asset_id", fields1[1].(bson.M)["field"])
	assert.Equal(t, float64(2), fields1[1].(bson.M)["value"])

	// Verify doc2: default_label → cesium_ion with asset_id 3
	var updatedDoc2 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc2["_id"]}).Decode(&updatedDoc2)
	require.NoError(t, err)
	fields2 := updatedDoc2["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Equal(t, "cesium_ion", fields2[0].(bson.M)["value"])
	assert.Equal(t, "cesium_ion_asset_id", fields2[1].(bson.M)["field"])
	assert.Equal(t, float64(3), fields2[1].(bson.M)["value"])

	// Verify doc3: default_road → cesium_ion with asset_id 4
	var updatedDoc3 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc3["_id"]}).Decode(&updatedDoc3)
	require.NoError(t, err)
	fields3 := updatedDoc3["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Equal(t, "cesium_ion", fields3[0].(bson.M)["value"])
	assert.Equal(t, "cesium_ion_asset_id", fields3[1].(bson.M)["field"])
	assert.Equal(t, float64(4), fields3[1].(bson.M)["value"])

	// Verify doc4: black_marble → cesium_ion with asset_id 3812
	var updatedDoc4 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc4["_id"]}).Decode(&updatedDoc4)
	require.NoError(t, err)
	fields4 := updatedDoc4["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Equal(t, "cesium_ion", fields4[0].(bson.M)["value"])
	assert.Equal(t, "cesium_ion_asset_id", fields4[1].(bson.M)["field"])
	assert.Equal(t, float64(3812), fields4[1].(bson.M)["value"])

	// Verify doc5: stamen_toner → carto_light
	var updatedDoc5 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc5["_id"]}).Decode(&updatedDoc5)
	require.NoError(t, err)
	value5 := updatedDoc5["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "carto_light", value5)

	// Verify doc6: esri_world_topo → open_street_map
	var updatedDoc6 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc6["_id"]}).Decode(&updatedDoc6)
	require.NoError(t, err)
	value6 := updatedDoc6["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "open_street_map", value6)

	// Verify doc7: unchanged
	var updatedDoc7 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc7["_id"]}).Decode(&updatedDoc7)
	require.NoError(t, err)
	value7 := updatedDoc7["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "open_street_map", value7)
}

func TestUpdateTileAndTerrainProviders_TerrainMigrations(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	// Test case 1: arcgis → reearth_terrain
	doc1 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "terrainType",
								"value": "arcgis",
							},
							bson.M{
								"field": "terrain",
								"value": true,
							},
						},
					},
				},
			},
		},
	}

	// Test case 2: cesium terrain should not be changed
	doc2 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "terrainType",
								"value": "cesium",
							},
						},
					},
				},
			},
		},
	}

	// Insert test documents
	_, err := db.Collection("property").InsertMany(ctx, []any{doc1, doc2})
	require.NoError(t, err)

	// Run migration
	err = UpdateTileAndTerrainProviders(ctx, c)
	require.NoError(t, err)

	// Verify doc1: arcgis → reearth_terrain
	var updatedDoc1 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc1["_id"]}).Decode(&updatedDoc1)
	require.NoError(t, err)
	fields1 := updatedDoc1["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Equal(t, "reearth_terrain", fields1[0].(bson.M)["value"])

	// Verify doc2: cesium unchanged
	var updatedDoc2 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc2["_id"]}).Decode(&updatedDoc2)
	require.NoError(t, err)
	value2 := updatedDoc2["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "cesium", value2)
}

func TestUpdateTileAndTerrainProviders_ExistingAssetId(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	// Document with default tile_type and existing cesium_ion_asset_id
	doc := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "default",
							},
							bson.M{
								"field": "cesium_ion_asset_id",
								"value": float64(999),
							},
						},
					},
				},
			},
		},
	}

	_, err := db.Collection("property").InsertOne(ctx, doc)
	require.NoError(t, err)

	// Run migration
	err = UpdateTileAndTerrainProviders(ctx, c)
	require.NoError(t, err)

	// Verify tile_type was updated but asset_id was preserved
	var updatedDoc bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc["_id"]}).Decode(&updatedDoc)
	require.NoError(t, err)
	fields := updatedDoc["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Equal(t, "cesium_ion", fields[0].(bson.M)["value"])
	assert.Equal(t, float64(999), fields[1].(bson.M)["value"]) // Asset ID should remain 999
}

func TestUpdateTileAndTerrainProviders_MultipleTilesInDocument(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	// Document with multiple tile and terrain fields
	doc := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "default",
							},
						},
					},
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "stamen_toner",
							},
						},
					},
				},
			},
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "terrainType",
								"value": "arcgis",
							},
						},
					},
				},
			},
		},
	}

	_, err := db.Collection("property").InsertOne(ctx, doc)
	require.NoError(t, err)

	// Run migration
	err = UpdateTileAndTerrainProviders(ctx, c)
	require.NoError(t, err)

	// Verify all fields were updated
	var updatedDoc bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc["_id"]}).Decode(&updatedDoc)
	require.NoError(t, err)

	items := updatedDoc["items"].(primitive.A)

	// Check first tile (default → cesium_ion)
	fields1 := items[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Equal(t, "cesium_ion", fields1[0].(bson.M)["value"])
	assert.Equal(t, "cesium_ion_asset_id", fields1[1].(bson.M)["field"])

	// Check second tile (stamen_toner → carto_light)
	fields2 := items[0].(bson.M)["groups"].(primitive.A)[1].(bson.M)["fields"].(primitive.A)
	assert.Equal(t, "carto_light", fields2[0].(bson.M)["value"])

	// Check terrain (arcgis → reearth_terrain)
	fields3 := items[1].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Equal(t, "reearth_terrain", fields3[0].(bson.M)["value"])
}

func TestUpdateTileAndTerrainProviders_EmptyCollection(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	// Run migration on empty collection
	err := UpdateTileAndTerrainProviders(ctx, c)
	require.NoError(t, err)

	// Verify no documents were added
	count, err := db.Collection("property").CountDocuments(ctx, bson.M{})
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestUpdateTileAndTerrainProviders_PreserveOtherFields(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	// Document with tile_type and other fields
	doc := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "default",
							},
							bson.M{
								"field": "other_field",
								"value": "preserve_me",
							},
							bson.M{
								"field": "another_field",
								"value": 123,
							},
						},
					},
				},
			},
		},
	}

	_, err := db.Collection("property").InsertOne(ctx, doc)
	require.NoError(t, err)

	// Run migration
	err = UpdateTileAndTerrainProviders(ctx, c)
	require.NoError(t, err)

	// Verify other fields were preserved
	var updatedDoc bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc["_id"]}).Decode(&updatedDoc)
	require.NoError(t, err)
	fields := updatedDoc["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)

	// Should have 4 fields: tile_type (updated), cesium_ion_asset_id (added), other_field, another_field
	assert.Len(t, fields, 4)
	assert.Equal(t, "cesium_ion", fields[0].(bson.M)["value"])
	assert.Equal(t, "preserve_me", fields[1].(bson.M)["value"])
	assert.Equal(t, int32(123), fields[2].(bson.M)["value"])
	assert.Equal(t, "cesium_ion_asset_id", fields[3].(bson.M)["field"])
}
