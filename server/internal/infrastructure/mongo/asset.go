package mongo

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/reearth/reearth/server/internal/infrastructure/mongo/mongodoc"
	"github.com/reearth/reearth/server/internal/usecase/repo"
	"github.com/reearth/reearth/server/pkg/asset"
	"github.com/reearth/reearth/server/pkg/id"
	"github.com/reearth/reearthx/account/accountdomain"
	"github.com/reearth/reearthx/mongox"
	"github.com/reearth/reearthx/rerror"
	"github.com/reearth/reearthx/usecasex"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	assetIndexes       = []string{"team"}
	assetUniqueIndexes = []string{"id"}
)

type Asset struct {
	client *mongox.ClientCollection
	f      repo.WorkspaceFilter
}

func NewAsset(client *mongox.Client) *Asset {
	return &Asset{client: client.WithCollection("asset")}
}

func (r *Asset) Init(ctx context.Context) error {
	return createIndexes(ctx, r.client, assetIndexes, assetUniqueIndexes)
}

func (r *Asset) Filtered(f repo.WorkspaceFilter) repo.Asset {
	return &Asset{
		client: r.client,
		f:      r.f.Merge(f),
	}
}

func (r *Asset) FindByID(ctx context.Context, id id.AssetID) (*asset.Asset, error) {
	return r.findOne(ctx, bson.M{
		"id": id.String(),
	})
}

func (r *Asset) FindByIDs(ctx context.Context, ids id.AssetIDList) ([]*asset.Asset, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	res, err := r.find(ctx, bson.M{
		"id": bson.M{"$in": ids.Strings()},
	})
	if err != nil {
		return nil, err
	}
	return filterAssets(ids, res), nil
}

func decodeAssetCursor(cursor usecasex.Cursor) (id string, key string) {
	parts := strings.SplitN(string(cursor), ":", 2)
	id = parts[0]
	if len(parts) > 1 {
		key = parts[1]
	}
	return
}

func encodeAssetCursor(a *asset.Asset, sort *asset.SortType) (cursor usecasex.Cursor) {
	var suffix string
	if sort == nil {
		suffix = ":" + a.CreatedAt().Format(time.RFC3339Nano)
	} else {
		switch *sort {
		case asset.SortTypeName:
			suffix = ":" + a.Name()
		case asset.SortTypeSize:
			suffix = ":" + fmt.Sprintf("%d", a.Size())
		case asset.SortTypeCreatedat:
			suffix = ":" + a.CreatedAt().Format(time.RFC3339Nano)
		default:
			suffix = ""
		}
	}

	cursor = usecasex.Cursor(a.ID().String() + suffix)
	return
}

func AssetFilterWithPagination(absoluteFilter bson.M, f repo.AssetFilter) (bson.M, *options.FindOptions, int64, int) {

	limit := int64(10)
	sortOrder := 1

	if f.Pagination.Cursor.First != nil {
		limit = *f.Pagination.Cursor.First
	} else if f.Pagination.Cursor.Last != nil {
		sortOrder = -1
		limit = *f.Pagination.Cursor.Last
	}

	limit = limit + 1

	if f.Sort == nil {
		f.Sort = &asset.SortTypeCreatedat
	}
	sortKey := string(*f.Sort)

	sort := bson.D{
		{Key: sortKey, Value: sortOrder},
		{Key: "id", Value: sortOrder},
	}

	if f.Pagination.Cursor.After != nil {
		cursor := *f.Pagination.Cursor.After
		afterID, afterKey := decodeAssetCursor(cursor)

		var keyValue any = afterKey
		if sortKey == "size" {
			if v, err := strconv.ParseInt(afterKey, 10, 64); err == nil {
				keyValue = v
			}
		} else if sortKey == "createdat" {
			if t, err := time.Parse(time.RFC3339Nano, afterKey); err == nil {
				keyValue = t
			}
		}

		absoluteFilter["$or"] = bson.A{
			bson.M{sortKey: bson.M{"$gt": keyValue}},
			bson.M{
				sortKey: keyValue,
				"id":    bson.M{"$gt": afterID},
			},
		}

	} else if f.Pagination.Cursor.Before != nil {
		cursor := *f.Pagination.Cursor.Before
		beforeID, beforeKey := decodeAssetCursor(cursor)

		var keyValue any = beforeKey
		if sortKey == "size" {
			if v, err := strconv.ParseInt(beforeKey, 10, 64); err == nil {
				keyValue = v
			}
		} else if sortKey == "createdat" {
			if t, err := time.Parse(time.RFC3339Nano, beforeKey); err == nil {
				keyValue = t
			}
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

func (r *Asset) FindByWorkspace(ctx context.Context, id accountdomain.WorkspaceID, f repo.AssetFilter) ([]*asset.Asset, *usecasex.PageInfo, error) {
	if !r.f.CanRead(id) {
		return nil, usecasex.EmptyPageInfo(), nil
	}

	absoluteFilter := bson.M{
		"team": id.String(),
	}

	if f.Keyword != nil {
		keywordRegex := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", regexp.QuoteMeta(*f.Keyword)),
			Options: "i",
		}
		absoluteFilter["name"] = bson.M{"$regex": keywordRegex}
	}

	totalCount, err := r.client.Client().CountDocuments(ctx, absoluteFilter)
	if err != nil {
		return nil, nil, err
	}

	paginationSortilter, findOptions, limit, sortOrder := AssetFilterWithPagination(absoluteFilter, f)

	cursor, err := r.client.Client().Find(ctx, paginationSortilter, findOptions)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("failed to close cursor: %v", cerr)
		}
	}()

	consumer := mongodoc.NewAssetConsumer(r.f.Readable)

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
		start := encodeAssetCursor(items[0], f.Sort)
		end := encodeAssetCursor(items[len(items)-1], f.Sort)
		startCursor = &start
		endCursor = &end
	}

	pageInfo := usecasex.NewPageInfo(totalCount, startCursor, endCursor, hasNextPage, hasPreviousPage)

	return items, pageInfo, nil
}

func (r *Asset) TotalSizeByWorkspace(ctx context.Context, wid accountdomain.WorkspaceID) (int64, error) {
	if !r.f.CanRead(wid) {
		return 0, repo.ErrOperationDenied
	}

	c, err := r.client.Client().Aggregate(ctx, []bson.M{
		{"$match": bson.M{"team": wid.String()}},
		{"$group": bson.M{"_id": nil, "size": bson.M{"$sum": "$size"}}},
	})
	if err != nil {
		return 0, rerror.ErrInternalByWithContext(ctx, err)
	}
	defer func() {
		_ = c.Close(ctx)
	}()

	if !c.Next(ctx) {
		return 0, nil
	}

	type resp struct {
		Size int64
	}
	var res resp
	if err := c.Decode(&res); err != nil {
		return 0, rerror.ErrInternalByWithContext(ctx, err)
	}
	return res.Size, nil
}

func (r *Asset) Save(ctx context.Context, asset *asset.Asset) error {
	if !r.f.CanWrite(asset.Workspace()) {
		return repo.ErrOperationDenied
	}
	doc, id := mongodoc.NewAsset(asset)
	return r.client.SaveOne(ctx, id, doc)
}

func (r *Asset) Remove(ctx context.Context, id id.AssetID) error {
	return r.client.RemoveOne(ctx, r.writeFilter(bson.M{
		"id": id.String(),
	}))
}

func (r *Asset) find(ctx context.Context, filter any) ([]*asset.Asset, error) {
	c := mongodoc.NewAssetConsumer(r.f.Readable)
	if err2 := r.client.Find(ctx, filter, c); err2 != nil {
		return nil, rerror.ErrInternalByWithContext(ctx, err2)
	}
	return c.Result, nil
}

func (r *Asset) findOne(ctx context.Context, filter any) (*asset.Asset, error) {
	c := mongodoc.NewAssetConsumer(r.f.Readable)
	if err := r.client.FindOne(ctx, filter, c); err != nil {
		return nil, err
	}
	return c.Result[0], nil
}

func filterAssets(ids []id.AssetID, rows []*asset.Asset) []*asset.Asset {
	res := make([]*asset.Asset, 0, len(ids))
	for _, id := range ids {
		var r2 *asset.Asset
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

// func (r *Asset) readFilter(filter any) any {
// 	return applyWorkspaceFilter(filter, r.f.Readable)
// }

func (r *Asset) writeFilter(filter any) any {
	return applyWorkspaceFilter(filter, r.f.Writable)
}
