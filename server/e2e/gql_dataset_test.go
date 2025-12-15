package e2e

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/reearth/reearth/server/internal/usecase/repo"
	"github.com/reearth/reearth/server/pkg/dataset"
	"github.com/reearth/reearth/server/pkg/id"
	"github.com/stretchr/testify/assert"
)

// export REEARTH_DB=mongodb://localhost
// go test -v -run TestGetDatasets ./e2e/...

func TestGetDatasets(t *testing.T) {
	e, ctx, r, _ := StartGQLServerRepos(t, baseSeeder)
	_, _, testSid := createScene(e, pID.String())
	tSid, err := id.SceneIDFrom(testSid)
	assert.Nil(t, err)

	testDataCount := 100
	SetupTestDatasetDatas(t, ctx, r, testDataCount, tSid)

	limit := 20

	// ASC order -----------------------------
	checkGetDatasetsASC(
		e,
		map[string]any{
			"sceneId": testSid,
			"last":    limit,
		},
		limit,
		testDataCount,
	)

	// DESC order -----------------------------
	checkGetDatasetsDESC(
		e,
		map[string]any{
			"sceneId": testSid,
			"first":   limit,
		},
		limit,
		testDataCount,
	)
}

// func ValueDump(val *httpexpect.Value) {
// 	raw := val.Raw()
// 	switch data := raw.(type) {
// 	case map[string]interface{}:
// 		if text, err := json.MarshalIndent(data, "", "  "); err == nil {
// 			fmt.Println(string(text))
// 		}
// 	case []interface{}:
// 		if text, err := json.MarshalIndent(data, "", "  "); err == nil {
// 			fmt.Println(string(text))
// 		}
// 	default:
// 		fmt.Println("Unsupported type:", reflect.TypeOf(raw))
// 	}
// }

func checkGetDatasetsASC(e *httpexpect.Expect, variables map[string]any, limit int, testDataCount int) {
	var endCursor string
	for i := 0; i < testDataCount/limit; i++ {
		if i != 0 {
			variables["before"] = endCursor // set cursor
		}
		res := GetDatasets(e, uID.String(), variables)
		// ValueDump(res)
		res.Object().Value("edges").Array().Length().IsEqual(limit)
		res.Object().Value("pageInfo").Object().HasValue("hasNextPage", false)
		res.Object().HasValue("totalCount", testDataCount)
		if i < testDataCount/limit-1 { // last
			res.Object().Value("pageInfo").Object().HasValue("hasPreviousPage", true)
		} else {
			res.Object().Value("pageInfo").Object().HasValue("hasPreviousPage", false)
		}
		endCursor = res.Path("$.pageInfo.endCursor").Raw().(string)
	}
}

func checkGetDatasetsDESC(e *httpexpect.Expect, variables map[string]any, limit int, testDataCount int) {
	var endCursor string
	for i := 0; i < testDataCount/limit; i++ {
		if i != 0 {
			variables["after"] = endCursor // set cursor
		}

		res := GetDatasets(e, uID.String(), variables)
		res.Object().Value("edges").Array().Length().IsEqual(limit)
		res.Object().Value("pageInfo").Object().HasValue("hasPreviousPage", false)
		res.Object().HasValue("totalCount", testDataCount)
		if i < testDataCount/limit-1 { // last
			res.Object().Value("pageInfo").Object().HasValue("hasNextPage", true)
		} else {
			res.Object().Value("pageInfo").Object().HasValue("hasNextPage", false)
		}
		endCursor = res.Path("$.pageInfo.endCursor").Raw().(string)
	}
}

func GetDatasets(e *httpexpect.Expect, u string, variables map[string]any) *httpexpect.Value {
	requestBody := GraphQLRequest{
		OperationName: "datasetsList",
		Query:         GetDatasetsQuery,
		Variables:     variables,
	}
	res := e.POST("/api/graphql").
		WithHeader("Origin", "https://example.com").
		WithHeader("authorization", "Bearer test").
		WithHeader("X-Reearth-Debug-User", u).
		WithHeader("Content-Type", "application/json").
		WithJSON(requestBody).
		Expect().
		Status(http.StatusOK).
		JSON()

	return res.Path("$.data.datasetSchemas")
}

const GetDatasetsQuery = `query datasetsList($sceneId: ID!, $first: Int, $last: Int, $after: Cursor, $before: Cursor) {
  datasetSchemas(
    sceneId: $sceneId
    first: $first
    last: $last
    after: $after
    before: $before
  ) {
    edges {
      node {
        id
        source
        name
        __typename
      }
      __typename
    }
    nodes {
      id
      source
      name
      __typename
    }
    pageInfo {
      endCursor
      hasNextPage
      hasPreviousPage
      startCursor
      __typename
    }
    totalCount
    __typename
  }
}`

func SetupTestDatasetDatas(t *testing.T, ctx context.Context, r *repo.Container, dataCount int, tSid id.SceneID) {

	startCount := 1000
	for i := startCount; i < startCount+dataCount; i++ {
		name := fmt.Sprintf("%d.csv", i)
		url := fmt.Sprintf("file:///%s", name)
		AddDataset(t, ctx, r, tSid, name, url)
	}
}

func AddDataset(t *testing.T, ctx context.Context, r *repo.Container, tSid id.SceneID, name string, u string) {
	dsfid := dataset.NewFieldID()

	dssf := dataset.NewSchemaField().
		ID(dsfid).
		Name("f1").
		Type(dataset.ValueTypeString).
		MustBuild()

	dss := dataset.NewSchema().
		NewID().
		Name(name).
		Scene(tSid).
		Fields([]*dataset.SchemaField{dssf}).
		Source(u).
		MustBuild()
	err := r.DatasetSchema.Save(ctx, dss)
	assert.Nil(t, err)

	dsf := dataset.NewField(dsfid, dataset.ValueTypeString.ValueFrom("test"), "")
	ds := dataset.New().
		NewID().
		Schema(dss.ID()).
		Scene(tSid).
		Fields([]*dataset.Field{dsf}).
		MustBuild()
	err = r.Dataset.Save(ctx, ds)
	assert.Nil(t, err)

}
