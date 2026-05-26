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

func TestRevertTileAndTerrainProviders_TileReversions(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	// Test case 1: cesium_ion (asset_id: 2) → default
	doc1 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "cesium_ion",
							},
							bson.M{
								"field": "cesium_ion_asset_id",
								"value": float64(2),
							},
						},
					},
				},
			},
		},
	}

	// Test case 2: cesium_ion (asset_id: 3) → default_label
	doc2 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "cesium_ion",
							},
							bson.M{
								"field": "cesium_ion_asset_id",
								"value": float64(3),
							},
						},
					},
				},
			},
		},
	}

	// Test case 3: cesium_ion (asset_id: 4) → default_road
	doc3 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "cesium_ion",
							},
							bson.M{
								"field": "cesium_ion_asset_id",
								"value": float64(4),
							},
						},
					},
				},
			},
		},
	}

	// Test case 4: cesium_ion (asset_id: 3812) → black_marble
	doc4 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "cesium_ion",
							},
							bson.M{
								"field": "cesium_ion_asset_id",
								"value": float64(3812),
							},
						},
					},
				},
			},
		},
	}

	// Test case 5: open_street_map → esri_world_topo
	doc5 := bson.M{
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

	// Test case 7: cesium_ion with different asset_id should not be changed
	doc6 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "cesium_ion",
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

	// Insert test documents
	_, err := db.Collection("property").InsertMany(ctx, []any{doc1, doc2, doc3, doc4, doc5, doc6})
	require.NoError(t, err)

	// Run revert migration
	err = RevertTileAndTerrainProviders(ctx, c)
	require.NoError(t, err)

	// Verify doc1: cesium_ion (asset_id: 2) → default (and asset_id removed)
	var updatedDoc1 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc1["_id"]}).Decode(&updatedDoc1)
	require.NoError(t, err)
	fields1 := updatedDoc1["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Len(t, fields1, 1)
	assert.Equal(t, "tile_type", fields1[0].(bson.M)["field"])
	assert.Equal(t, "default", fields1[0].(bson.M)["value"])

	// Verify doc2: cesium_ion (asset_id: 3) → default_label (and asset_id removed)
	var updatedDoc2 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc2["_id"]}).Decode(&updatedDoc2)
	require.NoError(t, err)
	fields2 := updatedDoc2["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Len(t, fields2, 1)
	assert.Equal(t, "tile_type", fields2[0].(bson.M)["field"])
	assert.Equal(t, "default_label", fields2[0].(bson.M)["value"])

	// Verify doc3: cesium_ion (asset_id: 4) → default_road (and asset_id removed)
	var updatedDoc3 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc3["_id"]}).Decode(&updatedDoc3)
	require.NoError(t, err)
	fields3 := updatedDoc3["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Len(t, fields3, 1)
	assert.Equal(t, "tile_type", fields3[0].(bson.M)["field"])
	assert.Equal(t, "default_road", fields3[0].(bson.M)["value"])

	// Verify doc4: cesium_ion (asset_id: 3812) → black_marble (and asset_id removed)
	var updatedDoc4 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc4["_id"]}).Decode(&updatedDoc4)
	require.NoError(t, err)
	fields4 := updatedDoc4["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Len(t, fields4, 1)
	assert.Equal(t, "tile_type", fields4[0].(bson.M)["field"])
	assert.Equal(t, "black_marble", fields4[0].(bson.M)["value"])

	// Verify doc5: open_street_map → esri_world_topo
	var updatedDoc5 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc5["_id"]}).Decode(&updatedDoc5)
	require.NoError(t, err)
	value5 := updatedDoc5["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "esri_world_topo", value5)

	// Verify doc6: cesium_ion with asset_id 999 unchanged
	var updatedDoc6 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc6["_id"]}).Decode(&updatedDoc6)
	require.NoError(t, err)
	fields6 := updatedDoc6["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Len(t, fields6, 2)
	assert.Equal(t, "cesium_ion", fields6[0].(bson.M)["value"])
	assert.Equal(t, float64(999), fields6[1].(bson.M)["value"])
}

func TestRevertTileAndTerrainProviders_TerrainReversions(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	// Test case 1: reearth_terrain → arcgis
	doc1 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "terrainType",
								"value": "reearth_terrain",
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

	// Run revert migration
	err = RevertTileAndTerrainProviders(ctx, c)
	require.NoError(t, err)

	// Verify doc1: reearth_terrain → arcgis
	var updatedDoc1 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc1["_id"]}).Decode(&updatedDoc1)
	require.NoError(t, err)
	fields1 := updatedDoc1["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Equal(t, "arcgis", fields1[0].(bson.M)["value"])

	// Verify doc2: cesium unchanged
	var updatedDoc2 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc2["_id"]}).Decode(&updatedDoc2)
	require.NoError(t, err)
	value2 := updatedDoc2["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "cesium", value2)
}

func TestRevertTileAndTerrainProviders_RoundTrip(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	// Original documents before any migration
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

	doc2 := bson.M{
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
						},
					},
				},
			},
		},
	}

	_, err := db.Collection("property").InsertMany(ctx, []any{doc1, doc2})
	require.NoError(t, err)

	// Step 1: Apply forward migration
	err = UpdateTileAndTerrainProviders(ctx, c)
	require.NoError(t, err)

	// Verify forward migration worked
	var afterForward1 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc1["_id"]}).Decode(&afterForward1)
	require.NoError(t, err)
	fields1 := afterForward1["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Equal(t, "cesium_ion", fields1[0].(bson.M)["value"])
	assert.Equal(t, "cesium_ion_asset_id", fields1[1].(bson.M)["field"])

	var afterForward2 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc2["_id"]}).Decode(&afterForward2)
	require.NoError(t, err)
	value2 := afterForward2["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "reearth_terrain", value2)

	// Step 2: Apply revert migration
	err = RevertTileAndTerrainProviders(ctx, c)
	require.NoError(t, err)

	// Verify we're back to original state
	var afterRevert1 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc1["_id"]}).Decode(&afterRevert1)
	require.NoError(t, err)
	fieldsReverted1 := afterRevert1["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Len(t, fieldsReverted1, 1)
	assert.Equal(t, "tile_type", fieldsReverted1[0].(bson.M)["field"])
	assert.Equal(t, "default", fieldsReverted1[0].(bson.M)["value"])

	var afterRevert2 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc2["_id"]}).Decode(&afterRevert2)
	require.NoError(t, err)
	valueReverted2 := afterRevert2["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "arcgis", valueReverted2)
}

func TestRevertTileAndTerrainProviders_MultipleTilesInDocument(t *testing.T) {
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
								"value": "cesium_ion",
							},
							bson.M{
								"field": "cesium_ion_asset_id",
								"value": float64(2),
							},
						},
					},
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
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "terrainType",
								"value": "reearth_terrain",
							},
						},
					},
				},
			},
		},
	}

	_, err := db.Collection("property").InsertOne(ctx, doc)
	require.NoError(t, err)

	// Run revert migration
	err = RevertTileAndTerrainProviders(ctx, c)
	require.NoError(t, err)

	// Verify all fields were reverted
	var updatedDoc bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc["_id"]}).Decode(&updatedDoc)
	require.NoError(t, err)

	items := updatedDoc["items"].(primitive.A)

	// Check first tile (cesium_ion → default, asset_id removed)
	fields1 := items[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Len(t, fields1, 1)
	assert.Equal(t, "default", fields1[0].(bson.M)["value"])

	// Check second tile (open_street_map → esri_world_topo)
	fields2 := items[0].(bson.M)["groups"].(primitive.A)[1].(bson.M)["fields"].(primitive.A)
	assert.Equal(t, "esri_world_topo", fields2[0].(bson.M)["value"])

	// Check terrain (reearth_terrain → arcgis)
	fields3 := items[1].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Equal(t, "arcgis", fields3[0].(bson.M)["value"])
}

func TestRevertTileAndTerrainProviders_EmptyCollection(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	// Run revert migration on empty collection
	err := RevertTileAndTerrainProviders(ctx, c)
	require.NoError(t, err)

	// Verify no documents were added
	count, err := db.Collection("property").CountDocuments(ctx, bson.M{})
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestRevertTileAndTerrainProviders_PreserveOtherFields(t *testing.T) {
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
								"value": "cesium_ion",
							},
							bson.M{
								"field": "cesium_ion_asset_id",
								"value": float64(2),
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

	// Run revert migration
	err = RevertTileAndTerrainProviders(ctx, c)
	require.NoError(t, err)

	// Verify other fields were preserved and asset_id was removed
	var updatedDoc bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc["_id"]}).Decode(&updatedDoc)
	require.NoError(t, err)
	fields := updatedDoc["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)

	// Should have 3 fields: tile_type (reverted to default), other_field, another_field
	// cesium_ion_asset_id should be removed
	assert.Len(t, fields, 3)
	assert.Equal(t, "default", fields[0].(bson.M)["value"])
	assert.Equal(t, "preserve_me", fields[1].(bson.M)["value"])
	assert.Equal(t, int32(123), fields[2].(bson.M)["value"])
}
