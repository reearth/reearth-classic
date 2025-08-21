package migration

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func ChangeEsriAndStamenToDefault(ctx context.Context, c DBClient) error {
	collection := c.WithCollection("property").Client()

	filter := bson.M{
		"items.groups.fields": bson.M{
			"$elemMatch": bson.M{
				"field": "tile_type",
				"value": bson.M{
					"$in": []string{"esri_world_topo", "stamen_watercolor", "stamen_toner"},
				},
			},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to find property documents: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("failed to close cursor: %v", cerr)
		}
	}()

	var updatedCount int64 = 0

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return fmt.Errorf("failed to decode document: %w", err)
		}

		documentID := doc["_id"]
		updated := false

		if items, ok := doc["items"].(bson.A); ok {
			for _, item := range items {
				if itemMap, ok := item.(bson.M); ok {
					if groups, ok := itemMap["groups"].(bson.A); ok {
						for _, group := range groups {
							if groupMap, ok := group.(bson.M); ok {
								if fields, ok := groupMap["fields"].(bson.A); ok {
									for _, field := range fields {
										if fieldMap, ok := field.(bson.M); ok {
											if fieldMap["field"] == "tile_type" {
												value, ok := fieldMap["value"].(string)
												if ok && (value == "esri_world_topo" || value == "stamen_watercolor" || value == "stamen_toner") {
													fieldMap["value"] = "default"
													updated = true
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}

		if updated {
			updateFilter := bson.M{"_id": documentID}
			updateDoc := bson.M{"$set": bson.M{"items": doc["items"]}}

			result, err := collection.UpdateOne(ctx, updateFilter, updateDoc)
			if err != nil {
				return fmt.Errorf("failed to update document with ID %v: %w", documentID, err)
			}

			if result.ModifiedCount > 0 {
				updatedCount++
				fmt.Printf("Updated document ID: %v\n", documentID)
			}
		}
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("cursor error: %w", err)
	}

	fmt.Printf("Migration completed. Updated %d documents.\n", updatedCount)
	return nil
}
