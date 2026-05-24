package migration

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MigrateTileTypesToNewProviders updates tile_type values to new provider names
// Migration rules:
// - "default" or missing tile_type → "cesium_ion" with cesium_ion_asset_id: 2
// - "default_label" → "cesium_ion" with cesium_ion_asset_id: 3
// - "default_road" → "cesium_ion" with cesium_ion_asset_id: 4
// - "black_marble" → "cesium_ion" with cesium_ion_asset_id: 3812
// - "stamen_toner" → "carto_light"
// - "esri_world_topo" → "open_street_map"
func MigrateTileTypesToNewProviders(ctx context.Context, c DBClient) error {
	col := c.WithCollection("property").Client()

	// Migration 1: default → cesium_ion (asset_id: 2)
	if err := migrateToCesiumIon(ctx, col, "default", 2); err != nil {
		return fmt.Errorf("failed to migrate 'default': %w", err)
	}

	// Migration 2: default_label → cesium_ion (asset_id: 3)
	if err := migrateToCesiumIon(ctx, col, "default_label", 3); err != nil {
		return fmt.Errorf("failed to migrate 'default_label': %w", err)
	}

	// Migration 3: default_road → cesium_ion (asset_id: 4)
	if err := migrateToCesiumIon(ctx, col, "default_road", 4); err != nil {
		return fmt.Errorf("failed to migrate 'default_road': %w", err)
	}

	// Migration 4: black_marble → cesium_ion (asset_id: 3812)
	if err := migrateToCesiumIon(ctx, col, "black_marble", 3812); err != nil {
		return fmt.Errorf("failed to migrate 'black_marble': %w", err)
	}

	// Migration 5: stamen_toner → carto_light (simple rename)
	if err := migrateSimpleRename(ctx, col, "stamen_toner", "carto_light"); err != nil {
		return fmt.Errorf("failed to migrate 'stamen_toner': %w", err)
	}

	// Migration 6: esri_world_topo → open_street_map (simple rename)
	if err := migrateSimpleRename(ctx, col, "esri_world_topo", "open_street_map"); err != nil {
		return fmt.Errorf("failed to migrate 'esri_world_topo': %w", err)
	}

	// Migration 7: Handle missing tile_type → cesium_ion (asset_id: 2)
	if err := migrateMissingTileType(ctx, col); err != nil {
		return fmt.Errorf("failed to migrate missing tile_type: %w", err)
	}

	fmt.Println("[migration] MigrateTileTypesToNewProviders completed successfully")
	return nil
}

// migrateToCesiumIon changes tile_type to cesium_ion and adds cesium_ion_asset_id
func migrateToCesiumIon(ctx context.Context, col *mongo.Collection, oldValue string, assetID int) error {
	filter := bson.M{
		"items.groups.fields.field": "tile_type",
		"items.groups.fields.value": oldValue,
	}

	countCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	n, err := col.CountDocuments(countCtx, filter)
	if err != nil {
		return fmt.Errorf("count failed for '%s': %w", oldValue, err)
	}
	fmt.Printf("[migration] target documents for '%s': %d\n", oldValue, n)
	if n == 0 {
		fmt.Printf("[migration] nothing to do for '%s'\n", oldValue)
		return nil
	}

	// Step 1: Update tile_type value
	update := bson.M{
		"$set": bson.M{
			"items.$[i].groups.$[g].fields.$[f].value": "cesium_ion",
		},
	}
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"i.groups": bson.M{"$type": "array"}},
			bson.M{"g.fields": bson.M{"$type": "array"}},
			bson.M{
				"f.field": "tile_type",
				"f.value": oldValue,
			},
		},
	}
	opts := options.Update().SetArrayFilters(arrayFilters)

	updateCtx, cancel2 := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel2()

	res, err := col.UpdateMany(updateCtx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("update tile_type failed for '%s': %w", oldValue, err)
	}
	fmt.Printf("[migration] '%s' → cesium_ion: matched: %d, modified: %d\n", oldValue, res.MatchedCount, res.ModifiedCount)

	// Step 2: Add cesium_ion_asset_id field
	// First, check if cesium_ion_asset_id already exists to avoid duplicates
	filterWithoutAssetID := bson.M{
		"items.groups.fields": bson.M{
			"$elemMatch": bson.M{
				"field": "tile_type",
				"value": "cesium_ion",
			},
		},
		"items.groups.fields.field": bson.M{"$ne": "cesium_ion_asset_id"},
	}

	updateAssetID := bson.M{
		"$push": bson.M{
			"items.$[i].groups.$[g].fields": bson.M{
				"field": "cesium_ion_asset_id",
				"type":  "number",
				"value": float64(assetID),
			},
		},
	}
	arrayFiltersAssetID := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"i.groups": bson.M{"$type": "array"}},
			bson.M{
				"g.fields": bson.M{
					"$elemMatch": bson.M{
						"field": "tile_type",
						"value": "cesium_ion",
					},
				},
			},
		},
	}
	optsAssetID := options.Update().SetArrayFilters(arrayFiltersAssetID)

	res2, err := col.UpdateMany(updateCtx, filterWithoutAssetID, updateAssetID, optsAssetID)
	if err != nil {
		return fmt.Errorf("add cesium_ion_asset_id failed for '%s': %w", oldValue, err)
	}
	fmt.Printf("[migration] added cesium_ion_asset_id=%d: matched: %d, modified: %d\n", assetID, res2.MatchedCount, res2.ModifiedCount)

	return nil
}

// migrateSimpleRename changes one tile_type value to another without adding fields
func migrateSimpleRename(ctx context.Context, col *mongo.Collection, oldValue, newValue string) error {
	filter := bson.M{
		"items.groups.fields.field": "tile_type",
		"items.groups.fields.value": oldValue,
	}

	countCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	n, err := col.CountDocuments(countCtx, filter)
	if err != nil {
		return fmt.Errorf("count failed for '%s': %w", oldValue, err)
	}
	fmt.Printf("[migration] target documents for '%s': %d\n", oldValue, n)
	if n == 0 {
		fmt.Printf("[migration] nothing to do for '%s'\n", oldValue)
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
				"f.field": "tile_type",
				"f.value": oldValue,
			},
		},
	}
	opts := options.Update().SetArrayFilters(arrayFilters)

	updateCtx, cancel2 := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel2()

	res, err := col.UpdateMany(updateCtx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("update failed for '%s': %w", oldValue, err)
	}

	fmt.Printf("[migration] '%s' → '%s': matched: %d, modified: %d\n", oldValue, newValue, res.MatchedCount, res.ModifiedCount)
	return nil
}

// migrateMissingTileType handles properties where tile_type field doesn't exist
func migrateMissingTileType(ctx context.Context, col *mongo.Collection) error {
	// Find properties with groups that have imagery-related fields but no tile_type
	// This is tricky - we need to identify imagery property schemas
	filter := bson.M{
		"$and": bson.A{
			bson.M{"items.groups.fields": bson.M{"$exists": true}},
			bson.M{"items.groups.fields.field": bson.M{"$ne": "tile_type"}},
			// Additional check: only target properties that might be imagery-related
			// Check if schemaplugin is reearth and schemaname might be imagery
			bson.M{
				"$or": bson.A{
					bson.M{"schemaplugin": "reearth", "schemaname": "default"},
					bson.M{"schemaplugin": "reearth", "schemaname": bson.M{"$regex": "^(default|tile)"}},
				},
			},
		},
	}

	countCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	n, err := col.CountDocuments(countCtx, filter)
	if err != nil {
		return fmt.Errorf("count failed for missing tile_type: %w", err)
	}
	fmt.Printf("[migration] properties potentially missing tile_type: %d\n", n)
	if n == 0 {
		fmt.Println("[migration] nothing to do for missing tile_type")
		return nil
	}

	// Add both tile_type and cesium_ion_asset_id fields
	update := bson.M{
		"$push": bson.M{
			"items.$[i].groups.$[g].fields": bson.M{
				"$each": bson.A{
					bson.M{
						"field": "tile_type",
						"type":  "string",
						"value": "cesium_ion",
					},
					bson.M{
						"field": "cesium_ion_asset_id",
						"type":  "number",
						"value": float64(2),
					},
				},
			},
		},
	}
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"i.groups": bson.M{"$type": "array"}},
			bson.M{"g.fields": bson.M{"$type": "array"}},
		},
	}
	opts := options.Update().SetArrayFilters(arrayFilters)

	updateCtx, cancel2 := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel2()

	res, err := col.UpdateMany(updateCtx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("add tile_type to missing properties failed: %w", err)
	}

	fmt.Printf("[migration] added tile_type=cesium_ion and asset_id=2 to missing: matched: %d, modified: %d\n", res.MatchedCount, res.ModifiedCount)
	return nil
}
