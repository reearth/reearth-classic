package migration

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateTileAndTerrainProviders updates tile and terrain type values to new provider names
// Tile migration rules:
// - "default" → "cesium_ion" with cesium_ion_asset_id: 2
// - "default_label" → "cesium_ion" with cesium_ion_asset_id: 3
// - "default_road" → "cesium_ion" with cesium_ion_asset_id: 4
// - "black_marble" → "cesium_ion" with cesium_ion_asset_id: 3812
// - "stamen_toner" → "carto_light"
// - "esri_world_topo" → "open_street_map"
// Terrain migration rules:
// - "arcgis" → "reearth_terrain"
func UpdateTileAndTerrainProviders(ctx context.Context, c DBClient) error {
	col := c.WithCollection("property").Client()

	// Tile migrations
	// Migration 1: default → cesium_ion (asset_id: 2)
	if err := migrateTileToCesiumIon(ctx, col, "default", 2); err != nil {
		return fmt.Errorf("failed to migrate tile 'default': %w", err)
	}

	// Migration 2: default_label → cesium_ion (asset_id: 3)
	if err := migrateTileToCesiumIon(ctx, col, "default_label", 3); err != nil {
		return fmt.Errorf("failed to migrate tile 'default_label': %w", err)
	}

	// Migration 3: default_road → cesium_ion (asset_id: 4)
	if err := migrateTileToCesiumIon(ctx, col, "default_road", 4); err != nil {
		return fmt.Errorf("failed to migrate tile 'default_road': %w", err)
	}

	// Migration 4: black_marble → cesium_ion (asset_id: 3812)
	if err := migrateTileToCesiumIon(ctx, col, "black_marble", 3812); err != nil {
		return fmt.Errorf("failed to migrate tile 'black_marble': %w", err)
	}

	// Migration 5: stamen_toner → carto_light (simple rename)
	if err := migrateTileSimpleRename(ctx, col, "stamen_toner", "carto_light"); err != nil {
		return fmt.Errorf("failed to migrate tile 'stamen_toner': %w", err)
	}

	// Migration 6: esri_world_topo → open_street_map (simple rename)
	if err := migrateTileSimpleRename(ctx, col, "esri_world_topo", "open_street_map"); err != nil {
		return fmt.Errorf("failed to migrate tile 'esri_world_topo': %w", err)
	}

	// Terrain migrations
	// Migration 7: arcgis → reearth_terrain (simple rename)
	if err := migrateTerrainSimpleRename(ctx, col, "arcgis", "reearth_terrain"); err != nil {
		return fmt.Errorf("failed to migrate terrain 'arcgis': %w", err)
	}

	fmt.Println("[migration] UpdateTileAndTerrainProviders completed successfully")
	return nil
}

// migrateTileToCesiumIon changes tile_type to cesium_ion and adds cesium_ion_asset_id.
// Uses two sequential operations. Step 2 uses a group-level filter to target only the
// groups that were just converted, so documents with multiple different legacy tile types
// are handled correctly (each type gets its own asset_id without cross-contamination).
func migrateTileToCesiumIon(ctx context.Context, col *mongo.Collection, oldValue string, assetID int) error {
	filter := bson.M{
		"items.groups.fields.field": "tile_type",
		"items.groups.fields.value": oldValue,
	}

	countCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Collect matching document IDs for auditing before modifying
	cursor, err := col.Find(countCtx, filter, &options.FindOptions{
		Projection: bson.M{"id": 1},
	})
	if err != nil {
		return fmt.Errorf("find failed for tile '%s': %w", oldValue, err)
	}
	var matchedDocs []struct {
		ID string `bson:"id"`
	}
	if err := cursor.All(countCtx, &matchedDocs); err != nil {
		return fmt.Errorf("cursor failed for tile '%s': %w", oldValue, err)
	}
	n := int64(len(matchedDocs))

	fmt.Printf("[migration] target documents for tile '%s': %d\n", oldValue, n)
	if n == 0 {
		fmt.Printf("[migration] nothing to do for tile '%s'\n", oldValue)
		return nil
	}
	for _, doc := range matchedDocs {
		fmt.Printf("[migration]   property id=%s will be migrated: %s -> cesium_ion (asset_id=%d)\n", doc.ID, oldValue, assetID)
	}

	updateCtx, cancel2 := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel2()

	// Step 1: Update tile_type value from oldValue to cesium_ion
	step1 := bson.M{
		"$set": bson.M{
			"items.$[i].groups.$[g].fields.$[f].value": "cesium_ion",
		},
	}
	step1ArrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"i.groups": bson.M{"$type": "array"}},
			bson.M{"g.fields": bson.M{"$type": "array"}},
			bson.M{"f.field": "tile_type", "f.value": oldValue},
		},
	}
	res1, err := col.UpdateMany(updateCtx, filter, step1, options.Update().SetArrayFilters(step1ArrayFilters))
	if err != nil {
		return fmt.Errorf("update tile_type failed for '%s': %w", oldValue, err)
	}
	fmt.Printf("[migration] tile '%s' → cesium_ion: matched=%d modified=%d\n", oldValue, res1.MatchedCount, res1.ModifiedCount)

	// Step 2: Add cesium_ion_asset_id to groups that now have cesium_ion but no asset_id yet.
	// The group-level $[g] filter targets only groups matching BOTH conditions, so groups
	// that already had cesium_ion with their own asset_id are not affected.
	filterForAssetID := bson.M{
		"items.groups": bson.M{
			"$elemMatch": bson.M{
				"$and": []interface{}{
					bson.M{"fields": bson.M{"$elemMatch": bson.M{"field": "tile_type", "value": "cesium_ion"}}},
					bson.M{"fields": bson.M{"$not": bson.M{"$elemMatch": bson.M{"field": "cesium_ion_asset_id"}}}},
				},
			},
		},
	}
	step2 := bson.M{
		"$push": bson.M{
			"items.$[i].groups.$[g].fields": bson.M{
				"field": "cesium_ion_asset_id",
				"type":  "number",
				"value": float64(assetID),
			},
		},
	}
	step2ArrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"i.groups": bson.M{"$type": "array"}},
			bson.M{
				"$and": []interface{}{
					bson.M{"g.fields": bson.M{"$elemMatch": bson.M{"field": "tile_type", "value": "cesium_ion"}}},
					bson.M{"g.fields": bson.M{"$not": bson.M{"$elemMatch": bson.M{"field": "cesium_ion_asset_id"}}}},
				},
			},
		},
	}
	res2, err := col.UpdateMany(updateCtx, filterForAssetID, step2, options.Update().SetArrayFilters(step2ArrayFilters))
	if err != nil {
		return fmt.Errorf("add cesium_ion_asset_id failed for '%s': %w", oldValue, err)
	}
	fmt.Printf("[migration] added cesium_ion_asset_id=%d: matched=%d modified=%d\n", assetID, res2.MatchedCount, res2.ModifiedCount)

	return nil
}

// migrateTileSimpleRename changes one tile_type value to another without adding fields
func migrateTileSimpleRename(ctx context.Context, col *mongo.Collection, oldValue, newValue string) error {
	filter := bson.M{
		"items.groups.fields.field": "tile_type",
		"items.groups.fields.value": oldValue,
	}

	countCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cursor, err := col.Find(countCtx, filter, &options.FindOptions{Projection: bson.M{"id": 1}})
	if err != nil {
		return fmt.Errorf("find failed for tile '%s': %w", oldValue, err)
	}
	var matchedDocs []struct {
		ID string `bson:"id"`
	}
	if err := cursor.All(countCtx, &matchedDocs); err != nil {
		return fmt.Errorf("cursor failed for tile '%s': %w", oldValue, err)
	}
	n := int64(len(matchedDocs))

	fmt.Printf("[migration] target documents for tile '%s': %d\n", oldValue, n)
	if n == 0 {
		fmt.Printf("[migration] nothing to do for tile '%s'\n", oldValue)
		return nil
	}
	for _, doc := range matchedDocs {
		fmt.Printf("[migration]   property id=%s will be migrated: %s -> %s\n", doc.ID, oldValue, newValue)
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
		return fmt.Errorf("update failed for tile '%s': %w", oldValue, err)
	}

	fmt.Printf("[migration] tile '%s' → '%s': matched: %d, modified: %d\n", oldValue, newValue, res.MatchedCount, res.ModifiedCount)
	return nil
}

// migrateTerrainSimpleRename changes one terrainType value to another
func migrateTerrainSimpleRename(ctx context.Context, col *mongo.Collection, oldValue, newValue string) error {
	filter := bson.M{
		"items.groups.fields.field": "terrainType",
		"items.groups.fields.value": oldValue,
	}

	countCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cursor, err := col.Find(countCtx, filter, &options.FindOptions{Projection: bson.M{"id": 1}})
	if err != nil {
		return fmt.Errorf("find failed for terrainType '%s': %w", oldValue, err)
	}
	var matchedDocs []struct {
		ID string `bson:"id"`
	}
	if err := cursor.All(countCtx, &matchedDocs); err != nil {
		return fmt.Errorf("cursor failed for terrainType '%s': %w", oldValue, err)
	}
	n := int64(len(matchedDocs))

	fmt.Printf("[migration] target documents for terrainType '%s': %d\n", oldValue, n)
	if n == 0 {
		fmt.Printf("[migration] nothing to do for terrainType '%s'\n", oldValue)
		return nil
	}
	for _, doc := range matchedDocs {
		fmt.Printf("[migration]   property id=%s will be migrated: terrainType %s -> %s\n", doc.ID, oldValue, newValue)
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
