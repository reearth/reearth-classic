package migration

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ChangeEsriAndStamenToDefault(ctx context.Context, c DBClient) error {
	col := c.WithCollection("property").Client()

	filter := bson.M{
		"items.groups.fields.field": "tile_type",
		"items.groups.fields.value": bson.M{
			"$in": bson.A{"esri_world_topo", "stamen_watercolor", "stamen_toner"},
		},
	}

	countCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	n, err := col.CountDocuments(countCtx, filter)
	if err != nil {
		return fmt.Errorf("count failed: %w", err)
	}
	fmt.Printf("[migration] target documents: %d\n", n)
	if n == 0 {
		fmt.Println("[migration] nothing to do")
		return nil
	}

	update := bson.M{
		"$set": bson.M{
			"items.$[i].groups.$[g].fields.$[f].value": "default",
		},
	}
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"i.groups": bson.M{"$type": "array"}},
			bson.M{"g.fields": bson.M{"$type": "array"}},
			bson.M{
				"f.field": "tile_type",
				"f.value": bson.M{
					"$in": bson.A{"esri_world_topo", "stamen_watercolor", "stamen_toner"},
				},
			},
		},
	}
	opts := options.Update().SetArrayFilters(arrayFilters)

	updateCtx, cancel2 := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel2()

	res, err := col.UpdateMany(updateCtx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	fmt.Printf("[migration] matched: %d, modified: %d\n", res.MatchedCount, res.ModifiedCount)
	fmt.Printf("[migration] Migration completed. Updated %d documents.\n", res.ModifiedCount)
	return nil
}
