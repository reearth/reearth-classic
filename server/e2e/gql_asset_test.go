package e2e

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/reearth/reearth/server/internal/usecase/repo"
	"github.com/reearth/reearth/server/pkg/asset"
	"github.com/stretchr/testify/assert"
)

// export REEARTH_DB=mongodb://localhost
// go test -v -run TestGetAssets ./e2e/...

func TestGetAssets(t *testing.T) {
	e, ctx, r, _ := StartGQLServerRepos(t, baseSeeder)

	testDataCount := 100
	SetupTestAssetDatas(t, ctx, r, wID, testDataCount)

	limit := 20

	// ASC order -----------------------------
	checkGetAssetsASC(
		e,
		map[string]any{
			"teamId": wID.String(),
			"pagination": map[string]any{
				"last": limit,
			},
		},
		limit,
		testDataCount,
	)

	checkGetAssetsASC(
		e,
		map[string]any{
			"teamId": wID.String(),
			"pagination": map[string]any{
				"last": limit,
			},
			"sort": "NAME",
		},
		limit,
		testDataCount,
	)
	checkGetAssetsASC(
		e,
		map[string]any{
			"teamId": wID.String(),
			"pagination": map[string]any{
				"last": limit,
			},
			"sort": "SIZE",
		},
		limit,
		testDataCount,
	)
	checkGetAssetsASC(
		e,
		map[string]any{
			"teamId": wID.String(),
			"pagination": map[string]any{
				"last": limit,
			},
			"sort": "DATE",
		},
		limit,
		testDataCount,
	)

	// DESC order -----------------------------
	checkGetAssetsDESC(
		e,
		map[string]any{
			"teamId": wID.String(),
			"pagination": map[string]any{
				"first": limit,
			},
		},
		limit,
		testDataCount,
	)
	checkGetAssetsDESC(
		e,
		map[string]any{
			"teamId": wID.String(),
			"pagination": map[string]any{
				"first": limit,
			},
			"sort": "NAME",
		},
		limit,
		testDataCount,
	)
	checkGetAssetsDESC(
		e,
		map[string]any{
			"teamId": wID.String(),
			"pagination": map[string]any{
				"first": limit,
			},
			"sort": "SIZE",
		},
		limit,
		testDataCount,
	)
	checkGetAssetsDESC(
		e,
		map[string]any{
			"teamId": wID.String(),
			"pagination": map[string]any{
				"first": limit,
			},
			"sort": "DATE",
		},
		limit,
		testDataCount,
	)

}

func TestGetAssetsKeyword(t *testing.T) {

	e, ctx, r, _ := StartGQLServerRepos(t, baseSeeder)

	testDataCount := 100
	SetupTestAssetDatas(t, ctx, r, wID, testDataCount)

	limit := 10

	var endCursor string

	// ASC order -----------------------------
	variables := map[string]any{
		"teamId": wID.String(),
		"pagination": map[string]any{
			"last": limit,
		},
		"keyword": "2", // count 19
	}

	for i := 0; i < 1; i++ {
		if i != 0 {
			variables["pagination"].(map[string]any)["before"] = endCursor
		}
		res := GetAssets(e, uID.String(), variables)
		if i < 1 {
			res.Object().Value("edges").Array().Length().Equal(limit)
		} else {
			res.Object().Value("edges").Array().Length().Equal(9)
		}
		res.Object().Value("pageInfo").Object().ValueEqual("hasNextPage", false)
		res.Object().ValueEqual("totalCount", 19)
		if i < 1 {
			res.Object().Value("pageInfo").Object().ValueEqual("hasPreviousPage", true)
		} else {
			res.Object().Value("pageInfo").Object().ValueEqual("hasPreviousPage", false)
		}
		endCursor = res.Path("$.pageInfo.endCursor").Raw().(string)
	}

	// DESC order -----------------------------
	variables = map[string]any{
		"teamId": wID.String(),
		"pagination": map[string]any{
			"first": limit,
		},
		"keyword": "2", // count 19
	}

	for i := 0; i < 1; i++ {
		if i != 0 {
			variables["pagination"].(map[string]any)["after"] = endCursor
		}

		res := GetAssets(e, uID.String(), variables)
		if i < 1 {
			res.Object().Value("edges").Array().Length().Equal(limit)
		} else {
			res.Object().Value("edges").Array().Length().Equal(9)
		}
		res.Object().Value("pageInfo").Object().ValueEqual("hasPreviousPage", false)
		res.Object().ValueEqual("totalCount", 19)
		if i < 1 {
			res.Object().Value("pageInfo").Object().ValueEqual("hasNextPage", true)
		} else {
			res.Object().Value("pageInfo").Object().ValueEqual("hasNextPage", false)
		}
		endCursor = res.Path("$.pageInfo.endCursor").Raw().(string)
	}
}

func TestGetAssetsLmit30(t *testing.T) {
	e, ctx, r, _ := StartGQLServerRepos(t, baseSeeder)

	testDataCount := 100
	SetupTestAssetDatas(t, ctx, r, wID, testDataCount)

	limit := 30

	var endCursor string

	// ASC order -----------------------------
	variables := map[string]any{
		"teamId": wID.String(),
		"pagination": map[string]any{
			"last": limit,
		},
	}
	for i := 0; i < 4; i++ {
		if i != 0 {
			variables["pagination"].(map[string]any)["before"] = endCursor
		}
		res := GetAssets(e, uID.String(), variables)
		if i < 3 {
			res.Object().Value("edges").Array().Length().Equal(limit)
		} else {
			res.Object().Value("edges").Array().Length().Equal(10)
		}
		res.Object().Value("pageInfo").Object().ValueEqual("hasNextPage", false)
		res.Object().ValueEqual("totalCount", testDataCount)
		if i < 3 {
			res.Object().Value("pageInfo").Object().ValueEqual("hasPreviousPage", true)
		} else {
			res.Object().Value("pageInfo").Object().ValueEqual("hasPreviousPage", false)
		}
		endCursor = res.Path("$.pageInfo.endCursor").Raw().(string)
	}

	// DESC order -----------------------------
	variables = map[string]any{
		"teamId": wID.String(),
		"pagination": map[string]any{
			"first": limit,
		},
	}
	for i := 0; i < 4; i++ {
		if i != 0 {
			variables["pagination"].(map[string]any)["after"] = endCursor
		}

		res := GetAssets(e, uID.String(), variables)
		if i < 3 {
			res.Object().Value("edges").Array().Length().Equal(limit)
		} else {
			res.Object().Value("edges").Array().Length().Equal(10)
		}
		res.Object().Value("pageInfo").Object().ValueEqual("hasPreviousPage", false)
		res.Object().ValueEqual("totalCount", testDataCount)
		if i < 3 {
			res.Object().Value("pageInfo").Object().ValueEqual("hasNextPage", true)
		} else {
			res.Object().Value("pageInfo").Object().ValueEqual("hasNextPage", false)
		}
		endCursor = res.Path("$.pageInfo.endCursor").Raw().(string)
	}
}

func checkGetAssetsASC(e *httpexpect.Expect, variables map[string]any, limit int, testDataCount int) {
	var endCursor string
	for i := 0; i < testDataCount/limit; i++ {
		if i != 0 {
			variables["pagination"].(map[string]any)["before"] = endCursor // set cursor
		}
		res := GetAssets(e, uID.String(), variables)
		res.Object().Value("edges").Array().Length().Equal(limit)
		res.Object().Value("pageInfo").Object().ValueEqual("hasNextPage", false)
		res.Object().ValueEqual("totalCount", testDataCount)
		if i < testDataCount/limit-1 { // last
			res.Object().Value("pageInfo").Object().ValueEqual("hasPreviousPage", true)
		} else {
			res.Object().Value("pageInfo").Object().ValueEqual("hasPreviousPage", false)
		}
		endCursor = res.Path("$.pageInfo.endCursor").Raw().(string)
	}
}

func checkGetAssetsDESC(e *httpexpect.Expect, variables map[string]any, limit int, testDataCount int) {
	var endCursor string
	for i := 0; i < testDataCount/limit; i++ {
		if i != 0 {
			variables["pagination"].(map[string]any)["after"] = endCursor // set cursor
		}

		res := GetAssets(e, uID.String(), variables)
		res.Object().Value("edges").Array().Length().Equal(limit)
		res.Object().Value("pageInfo").Object().ValueEqual("hasPreviousPage", false)
		res.Object().ValueEqual("totalCount", testDataCount)
		if i < testDataCount/limit-1 { // last
			res.Object().Value("pageInfo").Object().ValueEqual("hasNextPage", true)
		} else {
			res.Object().Value("pageInfo").Object().ValueEqual("hasNextPage", false)
		}
		endCursor = res.Path("$.pageInfo.endCursor").Raw().(string)
	}
}

func GetAssets(e *httpexpect.Expect, u string, variables map[string]any) *httpexpect.Value {
	requestBody := GraphQLRequest{
		OperationName: "GetAssets",
		Query:         GetAssetsQuery,
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
	return res.Path("$.data.assets")
}

const GetAssetsQuery = `query GetAssets(
  $teamId: ID!
  $sort: AssetSortType
  $keyword: String
  $pagination: Pagination
) {
  assets(
    teamId: $teamId
    keyword: $keyword
    sort: $sort
    pagination: $pagination
  ) {
    edges {
      cursor
      node {
        id
        teamId
        name
        size
        url
        contentType
        __typename
      }
      __typename
    }
    nodes {
      id
      teamId
      name
      size
      url
      contentType
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

func SetupTestAssetDatas(t *testing.T, ctx context.Context, r *repo.Container, w asset.WorkspaceID, dataCount int) {
	startTime := time.Now()
	startCount := 1000
	for i := startCount; i < startCount+dataCount; i++ {
		name := fmt.Sprintf("%d", i)
		size := int64(i)
		url := fmt.Sprintf("https://example.com/%s", name)
		createdAt := startTime.Add(time.Duration(i-1000) * time.Second)
		AddAsset(t, ctx, r, w, name, size, url, createdAt)
	}
}

func AddAsset(t *testing.T, ctx context.Context, r *repo.Container, workspace asset.WorkspaceID, name string, size int64, url string, createdAt time.Time) {
	a, err := asset.New().
		NewID().
		Workspace(workspace).
		Name(name).
		Size(size).
		URL(url).
		CreatedAt(createdAt).
		Build()
	assert.Nil(t, err)
	err = r.Asset.Save(ctx, a)
	assert.Nil(t, err)
}
