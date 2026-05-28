package mongo

import (
	"context"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/reearth/reearth/server/internal/infrastructure/mongo/mongodoc"
	"github.com/reearth/reearth/server/internal/usecase/repo"
	"github.com/reearth/reearth/server/pkg/id"
	"github.com/reearth/reearth/server/pkg/project"
	"github.com/reearth/reearthx/account/accountdomain"
	"github.com/reearth/reearthx/mongox"
	"github.com/reearth/reearthx/rerror"
	"github.com/reearth/reearthx/usecasex"
)

var (
	projectIndexes       = []string{"alias", "alias,publishmentstatus", "team", "workspace"}
	projectUniqueIndexes = []string{"id"}
)

type Project struct {
	client *mongox.ClientCollection
	f      repo.WorkspaceFilter
	s      repo.SceneFilter
}

func NewProject(client *mongox.Client) *Project {
	return &Project{
		client: client.WithCollection("project"),
	}
}

func (r *Project) Init(ctx context.Context) error {
	return createIndexes(ctx, r.client, projectIndexes, projectUniqueIndexes)
}

func (r *Project) Filtered(f repo.WorkspaceFilter) repo.Project {
	return &Project{
		client: r.client,
		f:      r.f.Merge(f),
	}
}

func (r *Project) FindByID(ctx context.Context, id id.ProjectID) (*project.Project, error) {
	return r.findOne(ctx, bson.M{
		"id": id.String(),
	}, true)
}

func (r *Project) FindByScene(ctx context.Context, id id.SceneID) (*project.Project, error) {
	if !r.s.CanRead(id) {
		return nil, nil
	}
	return r.findOne(ctx, bson.M{
		"scene": id.String(),
	}, true)
}

func (r *Project) FindByIDs(ctx context.Context, ids id.ProjectIDList) ([]*project.Project, error) {
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
	return filterProjects(ids, res), nil
}

func decodeProjectCursor(cursor usecasex.Cursor) (id string, key string) {
	parts := strings.SplitN(string(cursor), ":", 2)
	id = parts[0]
	if len(parts) > 1 {
		key = parts[1]
	}
	return
}

func encodeProjectCursor(p *project.Project) (cursor usecasex.Cursor) {
	suffix := ":" + p.UpdatedAt().Format(time.RFC3339Nano)
	cursor = usecasex.Cursor(p.ID().String() + suffix)
	return
}

func ProjectFilterWithPagination(absoluteFilter bson.M, pagination *usecasex.Pagination) (bson.M, *options.FindOptions, int64, int) {
	limit := int64(10)
	sortOrder := 1

	if pagination.Cursor.First != nil {
		limit = *pagination.Cursor.First
	} else if pagination.Cursor.Last != nil {
		sortOrder = -1
		limit = *pagination.Cursor.Last
	}

	limit = limit + 1

	sortKey := "updatedat"
	sort := bson.D{
		{Key: sortKey, Value: sortOrder},
		{Key: "id", Value: sortOrder},
	}

	if pagination.Cursor.After != nil {
		cursor := *pagination.Cursor.After
		afterID, afterKey := decodeProjectCursor(cursor)
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
		beforeID, beforeKey := decodeProjectCursor(cursor)
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

func (r *Project) FindByWorkspace(ctx context.Context, id accountdomain.WorkspaceID, pagination *usecasex.Pagination) ([]*project.Project, *usecasex.PageInfo, error) {
	if !r.f.CanRead(id) {
		return nil, usecasex.EmptyPageInfo(), nil
	}

	absoluteFilter := bson.M{
		"$and": []bson.M{
			{"$or": []bson.M{
				{"workspace": id.String()},
				{"team": id.String()},
			}},
			{"coresupport": bson.M{"$ne": true}},
		},
	}

	totalCount, err := r.client.Client().CountDocuments(ctx, absoluteFilter)
	if err != nil {
		return nil, nil, err
	}
	paginationSortilter, findOptions, limit, sortOrder := ProjectFilterWithPagination(absoluteFilter, pagination)

	cursor, err := r.client.Client().Find(ctx, paginationSortilter, findOptions)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("failed to close cursor: %v", cerr)
		}
	}()

	consumer := mongodoc.NewProjectConsumer(r.f.Readable)

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
		start := encodeProjectCursor(items[0])
		end := encodeProjectCursor(items[len(items)-1])
		startCursor = &start
		endCursor = &end
	}

	pageInfo := usecasex.NewPageInfo(totalCount, startCursor, endCursor, hasNextPage, hasPreviousPage)

	return items, pageInfo, nil

}

func (r *Project) FindByPublicName(ctx context.Context, name string) (*project.Project, error) {
	if name == "" {
		return nil, rerror.ErrNotFound
	}

	f := bson.D{
		{
			Key: "$or",
			Value: []bson.D{
				{{Key: "alias", Value: name}, {Key: "publishmentstatus", Value: bson.D{{Key: "$in", Value: []string{"public", "limited"}}}}},
				{{Key: "domains.domain", Value: name}, {Key: "publishmentstatus", Value: "public"}},
			},
		},
	}

	return r.findOne(ctx, f, false)
}

func (r *Project) CountByWorkspace(ctx context.Context, ws accountdomain.WorkspaceID) (int, error) {
	if !r.f.CanRead(ws) {
		return 0, repo.ErrOperationDenied
	}

	count, err := r.client.Count(ctx, bson.M{
		"team": ws.String(),
	})
	return int(count), err
}

func (r *Project) CountPublicByWorkspace(ctx context.Context, ws accountdomain.WorkspaceID) (int, error) {
	if !r.f.CanRead(ws) {
		return 0, repo.ErrOperationDenied
	}

	count, err := r.client.Count(ctx, bson.M{
		"team": ws.String(),
		"publishmentstatus": bson.M{
			"$in": []string{"public", "limited"},
		},
	})
	return int(count), err
}

func (r *Project) Save(ctx context.Context, project *project.Project) error {
	if !r.f.CanWrite(project.Workspace()) {
		return repo.ErrOperationDenied
	}
	doc, id := mongodoc.NewProject(project)
	return r.client.SaveOne(ctx, id, doc)
}

func (r *Project) Remove(ctx context.Context, id id.ProjectID) error {
	return r.client.RemoveOne(ctx, r.writeFilter(bson.M{"id": id.String()}))
}

func (r *Project) find(ctx context.Context, filter interface{}) ([]*project.Project, error) {
	c := mongodoc.NewProjectConsumer(r.f.Readable)
	if err := r.client.Find(ctx, filter, c); err != nil {
		return nil, err
	}
	return c.Result, nil
}

func (r *Project) findOne(ctx context.Context, filter any, filterByWorkspaces bool) (*project.Project, error) {
	var f []accountdomain.WorkspaceID
	if filterByWorkspaces {
		f = r.f.Readable
	}
	c := mongodoc.NewProjectConsumer(f)
	if err := r.client.FindOne(ctx, filter, c); err != nil {
		return nil, err
	}
	return c.Result[0], nil
}

func filterProjects(ids []id.ProjectID, rows []*project.Project) []*project.Project {
	res := make([]*project.Project, 0, len(ids))
	for _, id := range ids {
		var r2 *project.Project
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

// func (r *Project) readFilter(filter interface{}) interface{} {
// 	return applyWorkspaceFilter(filter, r.f.Readable)
// }

func (r *Project) writeFilter(filter interface{}) interface{} {
	return applyWorkspaceFilter(filter, r.f.Writable)
}
