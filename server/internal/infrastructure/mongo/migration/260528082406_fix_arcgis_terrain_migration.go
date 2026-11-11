package migration

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FixArcgisTerrain re-runs the arcgis → reearth_terrain rename that was skipped by
// migration 260525000715 due to a wrong document path.
//
// Terrain is a plain group item, so terrainType lives at items[].fields[], not
// items[].groups[].fields[] (which is the grouplist path used by tiles). The
// original migration queried the grouplist path and matched zero documents.
func FixArcgisTerrain(ctx context.Context, c DBClient) error {
	col := c.WithCollection("property").Client()

	if err := fixTerrainSimpleRename(ctx, col, "arcgis", "reearth_terrain"); err != nil {
		return fmt.Errorf("failed to fix terrain 'arcgis': %w", err)
	}

	fmt.Println("[migration] FixArcgisTerrain completed successfully")
	return nil
}

func fixTerrainSimpleRename(ctx context.Context, col *mongo.Collection, oldValue, newValue string) error {
	filter := bson.M{
		"items.fields.field": "terrainType",
		"items.fields.value": oldValue,
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
			"items.$[i].fields.$[f].value": newValue,
		},
	}
	arrayFilters := options.ArrayFilters{
		Filters: []any{
			bson.M{"i.fields": bson.M{"$type": "array"}},
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
