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

func TestChangeEsriAndStamenToDefault(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	// Test case 1: Document with esri_world_topo tile type
	doc1 := bson.M{
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
							bson.M{
								"field": "other_field",
								"value": "some_value",
							},
						},
					},
				},
			},
		},
	}

	// Test case 2: Document with stamen_watercolor tile type
	doc2 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "stamen_watercolor",
							},
						},
					},
				},
			},
		},
	}

	// Test case 3: Document with stamen_toner tile type
	doc3 := bson.M{
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

	// Test case 4: Document with a tile type that should not be changed
	doc4 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "openstreetmap",
							},
						},
					},
				},
			},
		},
	}

	// Test case 5: Document without tile_type field
	doc5 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "other_field",
								"value": "other_value",
							},
						},
					},
				},
			},
		},
	}

	// Insert test documents
	_, err := db.Collection("property").InsertMany(ctx, []any{doc1, doc2, doc3, doc4, doc5})
	require.NoError(t, err)

	// Run migration
	err = ChangeEsriAndStamenToDefault(ctx, c)
	require.NoError(t, err)

	// Verify doc1 was updated
	var updatedDoc1 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc1["_id"]}).Decode(&updatedDoc1)
	require.NoError(t, err)
	value1 := updatedDoc1["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "default", value1)

	// Verify doc2 was updated
	var updatedDoc2 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc2["_id"]}).Decode(&updatedDoc2)
	require.NoError(t, err)
	value2 := updatedDoc2["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "default", value2)

	// Verify doc3 was updated
	var updatedDoc3 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc3["_id"]}).Decode(&updatedDoc3)
	require.NoError(t, err)
	value3 := updatedDoc3["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "default", value3)

	// Verify doc4 was not changed (should remain "openstreetmap")
	var updatedDoc4 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc4["_id"]}).Decode(&updatedDoc4)
	require.NoError(t, err)
	value4 := updatedDoc4["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "openstreetmap", value4)

	// Verify doc5 was not changed
	var updatedDoc5 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc5["_id"]}).Decode(&updatedDoc5)
	require.NoError(t, err)
	fields5 := updatedDoc5["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)
	assert.Len(t, fields5, 1)
	assert.Equal(t, "other_field", fields5[0].(bson.M)["field"])
	assert.Equal(t, "other_value", fields5[0].(bson.M)["value"])
}

func TestChangeEsriAndStamenToDefault_MultipleFieldsInDocument(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	// Document with multiple tile_type fields that need updating
	doc := bson.M{
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
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "stamen_watercolor",
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
								"field": "tile_type",
								"value": "stamen_toner",
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
	err = ChangeEsriAndStamenToDefault(ctx, c)
	require.NoError(t, err)

	// Verify all tile_type fields were updated
	var updatedDoc bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc["_id"]}).Decode(&updatedDoc)
	require.NoError(t, err)

	items := updatedDoc["items"].(primitive.A)

	// Check first item, first group
	value1 := items[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "default", value1)

	// Check first item, second group
	value2 := items[0].(bson.M)["groups"].(primitive.A)[1].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "default", value2)

	// Check second item
	value3 := items[1].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "default", value3)
}

func TestChangeEsriAndStamenToDefault_EmptyCollection(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	// Run migration on empty collection
	err := ChangeEsriAndStamenToDefault(ctx, c)
	require.NoError(t, err)

	// Verify no documents were added
	count, err := db.Collection("property").CountDocuments(ctx, bson.M{})
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestChangeEsriAndStamenToDefault_MixedStructures(t *testing.T) {
	ctx := context.Background()
	db := mongotest.Connect(t)(t)
	c := mongox.NewClientWithDatabase(db)

	// Document 1: Valid structure with tile_type to update
	doc1 := bson.M{
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

	// Document 2: Invalid structure - groups is not an array
	// This document will be skipped by the migration filter
	doc2 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": "not_an_array",
			},
		},
	}

	// Document 3: Invalid structure - fields is missing
	// This document will be skipped by the migration filter
	doc3 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"other": "data",
					},
				},
			},
		},
	}

	// Document 4: Valid structure but different tile type
	doc4 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": bson.A{
					bson.M{
						"fields": bson.A{
							bson.M{
								"field": "tile_type",
								"value": "keep_this",
							},
						},
					},
				},
			},
		},
	}

	// Document 5: Invalid structure - items is nil
	doc5 := bson.M{
		"_id":   primitive.NewObjectID(),
		"items": nil,
	}

	// Document 6: Invalid structure - groups is nil
	doc6 := bson.M{
		"_id": primitive.NewObjectID(),
		"items": bson.A{
			bson.M{
				"groups": nil,
			},
		},
	}

	// Insert test documents
	_, err := db.Collection("property").InsertMany(ctx, []any{doc1, doc2, doc3, doc4, doc5, doc6})
	require.NoError(t, err)

	// Run migration
	err = ChangeEsriAndStamenToDefault(ctx, c)
	require.NoError(t, err)

	// Verify doc1 was updated
	var updatedDoc1 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc1["_id"]}).Decode(&updatedDoc1)
	require.NoError(t, err)
	value1 := updatedDoc1["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "default", value1)

	// Verify doc2 was not changed (invalid structure)
	var updatedDoc2 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc2["_id"]}).Decode(&updatedDoc2)
	require.NoError(t, err)
	assert.Equal(t, "not_an_array", updatedDoc2["items"].(primitive.A)[0].(bson.M)["groups"])

	// Verify doc3 was not changed (missing fields)
	var updatedDoc3 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc3["_id"]}).Decode(&updatedDoc3)
	require.NoError(t, err)
	groups3 := updatedDoc3["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)
	_, hasFields := groups3["fields"]
	assert.False(t, hasFields)

	// Verify doc4 was not changed (different tile type)
	var updatedDoc4 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc4["_id"]}).Decode(&updatedDoc4)
	require.NoError(t, err)
	value4 := updatedDoc4["items"].(primitive.A)[0].(bson.M)["groups"].(primitive.A)[0].(bson.M)["fields"].(primitive.A)[0].(bson.M)["value"]
	assert.Equal(t, "keep_this", value4)

	// Verify doc5 was not changed (nil items)
	var updatedDoc5 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc5["_id"]}).Decode(&updatedDoc5)
	require.NoError(t, err)
	assert.Nil(t, updatedDoc5["items"])

	// Verify doc6 was not changed (nil groups)
	var updatedDoc6 bson.M
	err = db.Collection("property").FindOne(ctx, bson.M{"_id": doc6["_id"]}).Decode(&updatedDoc6)
	require.NoError(t, err)
	assert.Nil(t, updatedDoc6["items"].(primitive.A)[0].(bson.M)["groups"])
}
