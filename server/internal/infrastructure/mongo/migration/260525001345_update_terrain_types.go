package migration

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateTerrainTypes updates terrain type values
// Migration rules:
// - If terrain is enabled but no terrainType field exists → set terrainType to "cesium"
// - If terrainType is "arcgis" → change to "reearth_terrain"
func UpdateTerrainTypes(ctx context.Context, c DBClient) error {
	col := c.WithCollection("property").Client()

	// Migration 1: arcgis → reearth_terrain (simple rename)
	if err := migrateTerrainType(ctx, col, "arcgis", "reearth_terrain"); err != nil {
		return fmt.Errorf("failed to migrate 'arcgis' terrain: %w", err)
	}

	// Migration 2: Add terrainType = "cesium" for enabled terrain without type
	if err := addMissingTerrainType(ctx, col); err != nil {
		return fmt.Errorf("failed to add missing terrainType: %w", err)
	}

	fmt.Println("[migration] UpdateTerrainTypes completed successfully")
	return nil
}

// migrateTerrainType changes one terrainType value to another
func migrateTerrainType(ctx context.Context, col *mongo.Collection, oldValue, newValue string) error {
	filter := bson.M{
		"items.groups.fields.field": "terrainType",
		"items.groups.fields.value": oldValue,
	}

	countCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	n, err := col.CountDocuments(countCtx, filter)
	if err != nil {
		return fmt.Errorf("count failed for terrainType '%s': %w", oldValue, err)
	}
	fmt.Printf("[migration] target documents for terrainType '%s': %d\n", oldValue, n)
	if n == 0 {
		fmt.Printf("[migration] nothing to do for terrainType '%s'\n", oldValue)
		return nil
	}

	update := bson.M{
		"$set": bson.M{
			"items.$[i].groups.$[g].fields.$[f].value": newValue,
		},
	}
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"i.groups": bson.M{"$type": "array"}},
			bson.M{"g.fields": bson.M{"$type": "array"}},
			bson.M{
				"f.field": "terrainType",
				"f.value": oldValue,
			},
		},
	}
	opts := options.Update().SetArrayFilters(arrayFilters)

	updateCtx, cancel2 := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel2()

	res, err := col.UpdateMany(updateCtx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("update failed for terrainType '%s': %w", oldValue, err)
	}

	fmt.Printf("[migration] terrainType '%s' → '%s': matched: %d, modified: %d\n", oldValue, newValue, res.MatchedCount, res.ModifiedCount)
	return nil
}

// addMissingTerrainType adds terrainType = "cesium" for properties where terrain is enabled but no terrainType exists
func addMissingTerrainType(ctx context.Context, col *mongo.Collection) error {
	// Find properties where:
	// 1. terrain field exists and is true
	// 2. terrainType field does NOT exist in the same group
	filter := bson.M{
		"items.groups.fields": bson.M{
			"$elemMatch": bson.M{
				"field": "terrain",
				"value": true,
			},
		},
		// Check that terrainType field doesn't exist
		"items.groups.fields.field": bson.M{"$ne": "terrainType"},
	}

	countCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	n, err := col.CountDocuments(countCtx, filter)
	if err != nil {
		return fmt.Errorf("count failed for missing terrainType: %w", err)
	}
	fmt.Printf("[migration] properties with terrain enabled but missing terrainType: %d\n", n)
	if n == 0 {
		fmt.Println("[migration] nothing to do for missing terrainType")
		return nil
	}

	// Add terrainType field to groups that have terrain=true
	update := bson.M{
		"$push": bson.M{
			"items.$[i].groups.$[g].fields": bson.M{
				"field": "terrainType",
				"type":  "string",
				"value": "cesium",
			},
		},
	}
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"i.groups": bson.M{"$type": "array"}},
			bson.M{
				"g.fields": bson.M{
					"$elemMatch": bson.M{
						"field": "terrain",
						"value": true,
					},
				},
			},
		},
	}
	opts := options.Update().SetArrayFilters(arrayFilters)

	updateCtx, cancel2 := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel2()

	res, err := col.UpdateMany(updateCtx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("add terrainType failed: %w", err)
	}

	fmt.Printf("[migration] added terrainType=cesium for enabled terrain: matched: %d, modified: %d\n", res.MatchedCount, res.ModifiedCount)
	return nil
}
