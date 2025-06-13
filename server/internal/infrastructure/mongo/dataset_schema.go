package mongo

import (
	"context"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"

	"github.com/reearth/reearth/server/internal/infrastructure/mongo/mongodoc"
	"github.com/reearth/reearth/server/internal/usecase/repo"
	"github.com/reearth/reearth/server/pkg/dataset"
	"github.com/reearth/reearth/server/pkg/id"
	"github.com/reearth/reearthx/mongox"
	"github.com/reearth/reearthx/rerror"
	"github.com/reearth/reearthx/usecasex"
)

var (
	datasetSchemaIndexes       = []string{"scene"}
	datasetSchemaUniqueIndexes = []string{"id"}
)

type DatasetSchema struct {
	client *mongox.ClientCollection
	f      repo.SceneFilter
}

func NewDatasetSchema(client *mongox.Client) *DatasetSchema {
	return &DatasetSchema{
		client: client.WithCollection("datasetSchema"),
	}
}

func (r *DatasetSchema) Init(ctx context.Context) error {
	return createIndexes(ctx, r.client, datasetSchemaIndexes, datasetSchemaUniqueIndexes)
}

func (r *DatasetSchema) Filtered(f repo.SceneFilter) repo.DatasetSchema {
	return &DatasetSchema{
		client: r.client,
		f:      r.f.Merge(f),
	}
}

func (r *DatasetSchema) FindByID(ctx context.Context, id id.DatasetSchemaID) (*dataset.Schema, error) {
	return r.findOne(ctx, bson.M{
		"id": id.String(),
	})
}

func (r *DatasetSchema) FindByIDs(ctx context.Context, ids id.DatasetSchemaIDList) (dataset.SchemaList, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	filter := bson.M{
		"id": bson.M{
			"$in": ids.Strings(),
		},
	}
	res, err := r.find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return filterDatasetSchemas(ids, res), nil
}

func decodeDatasetSchemaCursor(cursor usecasex.Cursor) (id string, key string) {
	parts := strings.SplitN(string(cursor), ":", 2)
	id = parts[0]
	if len(parts) > 1 {
		key = parts[1]
	}
	return
}

func encodeDatasetSchemaCursor(d *dataset.Schema) (cursor usecasex.Cursor) {
	suffix := ":" + d.Scene().String()
	cursor = usecasex.Cursor(d.ID().String() + suffix)
	return
}

func DatasetSchemaFilterWithPagination(absoluteFilter bson.M, pagination *usecasex.Pagination) (bson.M, *options.FindOptions, int64, int) {
	limit := int64(10)
	sortOrder := 1

	if pagination.Cursor.First != nil {
		limit = *pagination.Cursor.First
	} else if pagination.Cursor.Last != nil {
		sortOrder = -1
		limit = *pagination.Cursor.Last
	}

	limit = limit + 1

	sortKey := "scene"
	sort := bson.D{
		{Key: sortKey, Value: sortOrder},
		{Key: "id", Value: sortOrder},
	}

	if pagination.Cursor.After != nil {
		cursor := *pagination.Cursor.After
		afterID, afterKey := decodeDatasetSchemaCursor(cursor)
		var keyValue any = afterKey
		if t, err := time.Parse(time.RFC3339Nano, afterKey); err == nil {
			keyValue = t
		}
		absoluteFilter["$or"] = bson.A{
			bson.M{sortKey: bson.M{"$gt": keyValue}},
			bson.M{
				sortKey: keyValue,
				"id":    bson.M{"$gt": afterID},
			},
		}
	} else if pagination.Cursor.Before != nil {
		cursor := *pagination.Cursor.Before
		beforeID, beforeKey := decodeDatasetSchemaCursor(cursor)
		var keyValue any = beforeKey
		if t, err := time.Parse(time.RFC3339Nano, beforeKey); err == nil {
			keyValue = t
		}
		absoluteFilter["$or"] = bson.A{
			bson.M{sortKey: bson.M{"$lt": keyValue}},
			bson.M{
				sortKey: keyValue,
				"id":    bson.M{"$lt": beforeID},
			},
		}
	}
	findOptions := options.Find().
		SetSort(sort).
		SetLimit(limit)

	return absoluteFilter, findOptions, limit, sortOrder
}

func (r *DatasetSchema) FindByScene(ctx context.Context, sceneID id.SceneID, pagination *usecasex.Pagination) (dataset.SchemaList, *usecasex.PageInfo, error) {

	if !r.f.CanRead(sceneID) {
		return nil, usecasex.EmptyPageInfo(), nil
	}

	absoluteFilter := bson.M{
		"scene": sceneID.String(),
	}

	totalCount, err := r.client.Client().CountDocuments(ctx, absoluteFilter)
	if err != nil {
		return nil, nil, err
	}
	paginationSortilter, findOptions, limit, sortOrder := DatasetSchemaFilterWithPagination(absoluteFilter, pagination)

	cursor, err := r.client.Client().Find(ctx, paginationSortilter, findOptions)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("failed to close cursor: %v", cerr)
		}
	}()

	consumer := mongodoc.NewDatasetSchemaConsumer(r.f.Readable)

	for cursor.Next(ctx) {
		raw := cursor.Current
		if err := consumer.Consume(raw); err != nil {
			return nil, nil, err
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, nil, err
	}

	items := consumer.Result
	resultCount := int64(len(items))

	hasNextPage := false
	hasPreviousPage := false

	if resultCount == limit {
		if sortOrder == 1 {
			hasNextPage = true
		} else if sortOrder == -1 {
			hasPreviousPage = true
		}
		if len(items) > 0 {
			items = items[:len(items)-1]
		}
	}
	var startCursor, endCursor *usecasex.Cursor
	if len(items) > 0 {
		start := encodeDatasetSchemaCursor(items[0])
		end := encodeDatasetSchemaCursor(items[len(items)-1])
		startCursor = &start
		endCursor = &end
	}

	pageInfo := usecasex.NewPageInfo(totalCount, startCursor, endCursor, hasNextPage, hasPreviousPage)

	return items, pageInfo, nil

}

func (r *DatasetSchema) FindBySceneAll(ctx context.Context, sceneID id.SceneID) (dataset.SchemaList, error) {
	if !r.f.CanRead(sceneID) {
		return nil, nil
	}
	return r.find(ctx, bson.M{
		"scene": sceneID.String(),
	})
}

func (r *DatasetSchema) FindBySceneAndSource(ctx context.Context, sceneID id.SceneID, source string) (dataset.SchemaList, error) {
	if !r.f.CanRead(sceneID) {
		return nil, nil
	}
	return r.find(ctx, bson.M{
		"scene":  sceneID.String(),
		"source": string(source),
	})
}

func (r *DatasetSchema) CountByScene(ctx context.Context, id id.SceneID) (int, error) {
	if r.f.Readable != nil && !slices.Contains(r.f.Readable, id) {
		return 0, nil
	}

	res, err := r.client.Count(ctx, bson.M{
		"scene": id.String(),
	})
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

func (r *DatasetSchema) Save(ctx context.Context, datasetSchema *dataset.Schema) error {
	if !r.f.CanWrite(datasetSchema.Scene()) {
		return repo.ErrOperationDenied
	}
	doc, id := mongodoc.NewDatasetSchema(datasetSchema)
	return r.client.SaveOne(ctx, id, doc)
}

func (r *DatasetSchema) SaveAll(ctx context.Context, datasetSchemas dataset.SchemaList) error {
	if len(datasetSchemas) == 0 {
		return nil
	}
	docs, ids := mongodoc.NewDatasetSchemas(datasetSchemas, r.f.Writable)
	return r.client.SaveAll(ctx, ids, docs)
}

func (r *DatasetSchema) Remove(ctx context.Context, id id.DatasetSchemaID) error {
	return r.client.RemoveOne(ctx, r.writeFilter(bson.M{"id": id.String()}))
}

func (r *DatasetSchema) RemoveAll(ctx context.Context, ids id.DatasetSchemaIDList) error {
	if len(ids) == 0 {
		return nil
	}
	return r.client.RemoveAll(ctx, r.writeFilter(bson.M{
		"id": bson.M{"$in": ids.Strings()},
	}))
}

func (r *DatasetSchema) RemoveByScene(ctx context.Context, sceneID id.SceneID) error {
	if !r.f.CanWrite(sceneID) {
		return nil
	}
	if _, err := r.client.Client().DeleteMany(ctx, bson.M{
		"scene": sceneID.String(),
	}); err != nil {
		return rerror.ErrInternalByWithContext(ctx, err)
	}
	return nil
}

func (r *DatasetSchema) find(ctx context.Context, filter any) ([]*dataset.Schema, error) {
	c := mongodoc.NewDatasetSchemaConsumer(r.f.Readable)
	if err := r.client.Find(ctx, filter, c); err != nil {
		return nil, err
	}
	return c.Result, nil
}

func (r *DatasetSchema) findOne(ctx context.Context, filter any) (*dataset.Schema, error) {
	c := mongodoc.NewDatasetSchemaConsumer(r.f.Readable)
	if err := r.client.FindOne(ctx, filter, c); err != nil {
		return nil, err
	}
	return c.Result[0], nil
}

func filterDatasetSchemas(ids []id.DatasetSchemaID, rows []*dataset.Schema) []*dataset.Schema {
	res := make([]*dataset.Schema, 0, len(ids))
	for _, id := range ids {
		var r2 *dataset.Schema
		for _, r := range rows {
			if r.ID() == id {
				r2 = r
				break
			}
		}
		res = append(res, r2)
	}
	return res
}

// func (r *DatasetSchema) readFilter(filter any) any {
// 	return applySceneFilter(filter, r.f.Readable)
// }

func (r *DatasetSchema) writeFilter(filter any) any {
	return applySceneFilter(filter, r.f.Writable)
}
