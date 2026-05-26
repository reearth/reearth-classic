package migration

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RevertTileAndTerrainProviders reverts the changes made by UpdateTileAndTerrainProviders
// This is a recovery migration that should NOT be registered in migrations.go
// Use this only when you need to rollback the UpdateTileAndTerrainProviders migration
//
// Tile reversion rules:
// - "cesium_ion" with asset_id 2 → "default" (and remove asset_id field)
// - "cesium_ion" with asset_id 3 → "default_label" (and remove asset_id field)
// - "cesium_ion" with asset_id 4 → "default_road" (and remove asset_id field)
// - "cesium_ion" with asset_id 3812 → "black_marble" (and remove asset_id field)
// - "open_street_map" → "esri_world_topo" (best-effort: also covers stamen_toner → open_street_map)
// Terrain reversion rules:
// - "reearth_terrain" → "arcgis"
func RevertTileAndTerrainProviders(ctx context.Context, c DBClient) error {
	col := c.WithCollection("property").Client()

	// Tile reversions
	// Reversion 1: cesium_ion (asset_id: 2) → default
	if err := revertCesiumIonToOldType(ctx, col, 2, "default"); err != nil {
		return fmt.Errorf("failed to revert cesium_ion (asset_id: 2) to 'default': %w", err)
	}

	// Reversion 2: cesium_ion (asset_id: 3) → default_label
	if err := revertCesiumIonToOldType(ctx, col, 3, "default_label"); err != nil {
		return fmt.Errorf("failed to revert cesium_ion (asset_id: 3) to 'default_label': %w", err)
	}

	// Reversion 3: cesium_ion (asset_id: 4) → default_road
	if err := revertCesiumIonToOldType(ctx, col, 4, "default_road"); err != nil {
		return fmt.Errorf("failed to revert cesium_ion (asset_id: 4) to 'default_road': %w", err)
	}

	// Reversion 4: cesium_ion (asset_id: 3812) → black_marble
	if err := revertCesiumIonToOldType(ctx, col, 3812, "black_marble"); err != nil {
		return fmt.Errorf("failed to revert cesium_ion (asset_id: 3812) to 'black_marble': %w", err)
	}

	// Reversion 5: open_street_map → esri_world_topo (simple rename)
	// Note: stamen_toner also maps to open_street_map in the forward migration,
	// so this revert is best-effort — stamen_toner tiles will be reverted to esri_world_topo.
	if err := revertTileSimpleRename(ctx, col, "open_street_map", "esri_world_topo"); err != nil {
		return fmt.Errorf("failed to revert tile 'open_street_map': %w", err)
	}

	// Terrain reversions
	// Reversion 7: reearth_terrain → arcgis (simple rename)
	if err := revertTerrainSimpleRename(ctx, col, "reearth_terrain", "arcgis"); err != nil {
		return fmt.Errorf("failed to revert terrain 'reearth_terrain': %w", err)
	}

	fmt.Println("[migration] RevertTileAndTerrainProviders completed successfully")
	return nil
}

// revertCesiumIonToOldType changes cesium_ion with specific asset_id back to old tile type and removes asset_id field
func revertCesiumIonToOldType(ctx context.Context, col *mongo.Collection, assetID int, oldType string) error {
	// Find documents with cesium_ion and the specific asset_id
	filter := bson.M{
		"items.groups.fields": bson.M{
			"$elemMatch": bson.M{
				"field": "tile_type",
				"value": "cesium_ion",
			},
		},
		"items.groups.fields.field": "cesium_ion_asset_id",
		"items.groups.fields.value": float64(assetID),
	}

	countCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	n, err := col.CountDocuments(countCtx, filter)
	if err != nil {
		return fmt.Errorf("count failed for cesium_ion (asset_id: %d): %w", assetID, err)
	}
	fmt.Printf("[migration] target documents for cesium_ion (asset_id: %d): %d\n", assetID, n)
	if n == 0 {
		fmt.Printf("[migration] nothing to do for cesium_ion (asset_id: %d)\n", assetID)
		return nil
	}

	updateCtx, cancel2 := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel2()

	// Step 1: Update tile_type value from cesium_ion to old type
	update := bson.M{
		"$set": bson.M{
			"items.$[i].groups.$[g].fields.$[f].value": oldType,
		},
	}
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"i.groups": bson.M{"$type": "array"}},
			bson.M{
				"g.fields": bson.M{
					"$elemMatch": bson.M{
						"field": "cesium_ion_asset_id",
						"value": float64(assetID),
					},
				},
			},
			bson.M{
				"f.field": "tile_type",
				"f.value": "cesium_ion",
			},
		},
	}
	opts := options.Update().SetArrayFilters(arrayFilters)

	res, err := col.UpdateMany(updateCtx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("update tile_type failed for cesium_ion (asset_id: %d): %w", assetID, err)
	}
	fmt.Printf("[migration] cesium_ion (asset_id: %d) → '%s': matched: %d, modified: %d\n", assetID, oldType, res.MatchedCount, res.ModifiedCount)

	// Step 2: Remove cesium_ion_asset_id field
	filterForRemoval := bson.M{
		"$and": bson.A{
			bson.M{
				"items.groups.fields": bson.M{
					"$elemMatch": bson.M{
						"field": "tile_type",
						"value": oldType,
					},
				},
			},
			bson.M{
				"items.groups.fields": bson.M{
					"$elemMatch": bson.M{
						"field": "cesium_ion_asset_id",
						"value": float64(assetID),
					},
				},
			},
		},
	}

	updateRemoval := bson.M{
		"$pull": bson.M{
			"items.$[i].groups.$[g].fields": bson.M{
				"field": "cesium_ion_asset_id",
				"value": float64(assetID),
			},
		},
	}
	arrayFiltersRemoval := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"i.groups": bson.M{"$type": "array"}},
			bson.M{
				"g.fields": bson.M{
					"$elemMatch": bson.M{
						"field": "cesium_ion_asset_id",
						"value": float64(assetID),
					},
				},
			},
		},
	}
	optsRemoval := options.Update().SetArrayFilters(arrayFiltersRemoval)

	res2, err := col.UpdateMany(updateCtx, filterForRemoval, updateRemoval, optsRemoval)
	if err != nil {
		return fmt.Errorf("remove cesium_ion_asset_id failed for asset_id %d: %w", assetID, err)
	}
	fmt.Printf("[migration] removed cesium_ion_asset_id=%d: matched: %d, modified: %d\n", assetID, res2.MatchedCount, res2.ModifiedCount)

	return nil
}

// revertTileSimpleRename changes one tile_type value back to another
func revertTileSimpleRename(ctx context.Context, col *mongo.Collection, currentValue, oldValue string) error {
	filter := bson.M{
		"items.groups.fields.field": "tile_type",
		"items.groups.fields.value": currentValue,
	}

	countCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	n, err := col.CountDocuments(countCtx, filter)
	if err != nil {
		return fmt.Errorf("count failed for tile '%s': %w", currentValue, err)
	}
	fmt.Printf("[migration] target documents for tile '%s': %d\n", currentValue, n)
	if n == 0 {
		fmt.Printf("[migration] nothing to do for tile '%s'\n", currentValue)
		return nil
	}

	update := bson.M{
		"$set": bson.M{
			"items.$[i].groups.$[g].fields.$[f].value": oldValue,
		},
	}
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"i.groups": bson.M{"$type": "array"}},
			bson.M{"g.fields": bson.M{"$type": "array"}},
			bson.M{
				"f.field": "tile_type",
				"f.value": currentValue,
			},
		},
	}
	opts := options.Update().SetArrayFilters(arrayFilters)

	updateCtx, cancel2 := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel2()

	res, err := col.UpdateMany(updateCtx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("update failed for tile '%s': %w", currentValue, err)
	}

	fmt.Printf("[migration] tile '%s' → '%s': matched: %d, modified: %d\n", currentValue, oldValue, res.MatchedCount, res.ModifiedCount)
	return nil
}

// revertTerrainSimpleRename changes one terrainType value back to another
func revertTerrainSimpleRename(ctx context.Context, col *mongo.Collection, currentValue, oldValue string) error {
	filter := bson.M{
		"items.groups.fields.field": "terrainType",
		"items.groups.fields.value": currentValue,
	}

	countCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	n, err := col.CountDocuments(countCtx, filter)
	if err != nil {
		return fmt.Errorf("count failed for terrainType '%s': %w", currentValue, err)
	}
	fmt.Printf("[migration] target documents for terrainType '%s': %d\n", currentValue, n)
	if n == 0 {
		fmt.Printf("[migration] nothing to do for terrainType '%s'\n", currentValue)
		return nil
	}

	update := bson.M{
		"$set": bson.M{
			"items.$[i].groups.$[g].fields.$[f].value": oldValue,
		},
	}
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"i.groups": bson.M{"$type": "array"}},
			bson.M{"g.fields": bson.M{"$type": "array"}},
			bson.M{
				"f.field": "terrainType",
				"f.value": currentValue,
			},
		},
	}
	opts := options.Update().SetArrayFilters(arrayFilters)

	updateCtx, cancel2 := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel2()

	res, err := col.UpdateMany(updateCtx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("update failed for terrainType '%s': %w", currentValue, err)
	}

	fmt.Printf("[migration] terrainType '%s' → '%s': matched: %d, modified: %d\n", currentValue, oldValue, res.MatchedCount, res.ModifiedCount)
	return nil
}
