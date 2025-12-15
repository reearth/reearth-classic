package e2e

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/reearth/reearth/server/internal/app/config"
	"github.com/reearth/reearth/server/internal/usecase/repo"
	"github.com/reearth/reearth/server/pkg/project"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestCreateProject(t *testing.T) {
	e := StartServer(t, &config.Config{
		Origins: []string{"https://example.com"},
		AuthSrv: config.AuthSrvConfig{
			Disabled: true,
		},
	},
		true, baseSeeder)

	// create project with default coreSupport value (false)
	requestBody := GraphQLRequest{
		OperationName: "CreateProject",
		Query:         "mutation CreateProject($teamId: ID!, $visualizer: Visualizer!, $name: String!, $description: String!, $imageUrl: URL) {\n createProject(\n input: {teamId: $teamId, visualizer: $visualizer, name: $name, description: $description, imageUrl: $imageUrl}\n ) {\n project {\n id\n name\n description\n imageUrl\n coreSupport\n __typename\n }\n __typename\n }\n}",
		Variables: map[string]any{
			"name":        "test",
			"description": "abc",
			"imageUrl":    "",
			"teamId":      wID.String(),
			"visualizer":  "CESIUM",
		},
	}

	e.POST("/api/graphql").
		WithHeader("Origin", "https://example.com").
		WithHeader("authorization", "Bearer test").
		// WithHeader("authorization", "Bearer test").
		WithHeader("X-Reearth-Debug-User", uID.String()).
		WithHeader("Content-Type", "application/json").
		WithJSON(requestBody).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("data").Object().
		Value("createProject").Object().
		Value("project").Object().
		// ValueEqual("id", pId.String()).
		ValueEqual("name", "test").
		ValueEqual("description", "abc").
		ValueEqual("imageUrl", "").
		ValueEqual("coreSupport", false)

	// create coreSupport project
	requestBody = GraphQLRequest{
		OperationName: "CreateProject",
		Query:         "mutation CreateProject($teamId: ID!, $visualizer: Visualizer!, $name: String!, $description: String!, $imageUrl: URL, $coreSupport: Boolean) {\n createProject(\n input: {teamId: $teamId, visualizer: $visualizer, name: $name, description: $description, imageUrl: $imageUrl, coreSupport: $coreSupport}\n ) {\n project {\n id\n name\n description\n imageUrl\n coreSupport\n __typename\n }\n __typename\n }\n}",
		Variables: map[string]any{
			"name":        "test",
			"description": "abc",
			"imageUrl":    "",
			"teamId":      wID.String(),
			"visualizer":  "CESIUM",
			"coreSupport": true,
		},
	}

	e.POST("/api/graphql").
		WithHeader("Origin", "https://example.com").
		// WithHeader("authorization", "Bearer test").
		WithHeader("X-Reearth-Debug-User", uID.String()).
		WithHeader("Content-Type", "application/json").
		WithJSON(requestBody).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("data").Object().
		Value("createProject").Object().
		Value("project").Object().
		// ValueEqual("id", pId.String()).
		ValueEqual("name", "test").
		ValueEqual("description", "abc").
		ValueEqual("imageUrl", "").
		ValueEqual("coreSupport", true)
}

func TestGetProjects(t *testing.T) {
	e, ctx, r, _ := StartGQLServerRepos(t, baseSeeder)

	testDataCount := 99 // +1 base seeder
	SetupTestProjectDatas(t, ctx, r, wID, testDataCount)

	limit := 20

	// ASC order -----------------------------
	checkGetProjectsASC(
		e,
		map[string]any{
			"teamId": wID.String(),
			"last":   limit,
		},
		limit,
		testDataCount+1, // +1 base seeder
	)

	// DESC order -----------------------------
	checkGetProjectsDESC(
		e,
		map[string]any{
			"teamId": wID.String(),
			"first":  limit,
		},
		limit,
		testDataCount+1, // +1 base seeder
	)
}

func checkGetProjectsASC(e *httpexpect.Expect, variables map[string]any, limit int, testDataCount int) {
	var endCursor string
	for i := 0; i < testDataCount/limit; i++ {
		if i != 0 {
			variables["before"] = endCursor // set cursor
		}
		res := GetProjects(e, uID.String(), variables)
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

func checkGetProjectsDESC(e *httpexpect.Expect, variables map[string]any, limit int, testDataCount int) {
	var endCursor string
	for i := 0; i < testDataCount/limit; i++ {
		if i != 0 {
			variables["after"] = endCursor // set cursor
		}

		res := GetProjects(e, uID.String(), variables)
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

func GetProjects(e *httpexpect.Expect, u string, variables map[string]any) *httpexpect.Value {
	requestBody := GraphQLRequest{
		OperationName: "GetProjects",
		Query:         GetProjectsQuery,
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
	return res.Path("$.data.projects")
}

const GetProjectsQuery = `query GetProjects($teamId: ID!, $first: Int, $last: Int, $after: Cursor, $before: Cursor) {
  projects(
    teamId: $teamId
    first: $first
    last: $last
    after: $after
    before: $before
  ) {
    edges {
      node {
        id
        ...ProjectFragment
        scene {
          id
          __typename
        }
        __typename
      }
      __typename
    }
    nodes {
      id
      ...ProjectFragment
      scene {
        id
        __typename
      }
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
}

fragment ProjectFragment on Project {
  id
  name
  description
  imageUrl
  isArchived
  isBasicAuthActive
  basicAuthUsername
  basicAuthPassword
  publicTitle
  publicDescription
  publicImage
  alias
  publishmentStatus
  updatedAt
  coreSupport
  __typename
}`

func SetupTestProjectDatas(t *testing.T, ctx context.Context, r *repo.Container, w project.WorkspaceID, dataCount int) {
	startTime := time.Now()
	startCount := 1000
	for i := startCount; i < startCount+dataCount; i++ {
		name := fmt.Sprintf("%d", i)
		url := fmt.Sprintf("https://example.com/%s", name)
		createdAt := startTime.Add(time.Duration(i-1000) * time.Second)
		AddProject(t, ctx, r, w, name, url, createdAt)
	}
}

func AddProject(t *testing.T, ctx context.Context, r *repo.Container, workspace project.WorkspaceID, name string, u string, updatedAt time.Time) {
	p := project.New().
		NewID().
		Name(name).
		Description(name + " Description").
		ImageURL(lo.Must(url.Parse(u))).
		Workspace(workspace).
		Alias(pAlias).
		UpdatedAt(updatedAt).
		MustBuild()
	err := r.Project.Save(ctx, p)
	assert.Nil(t, err)
}
