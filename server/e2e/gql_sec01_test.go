package e2e

import (
	"context"
	"net/http"
	"testing"

	"github.com/reearth/reearth/server/internal/app/config"
	"github.com/reearth/reearth/server/internal/usecase/repo"
	"github.com/reearth/reearth/server/pkg/project"
)

const sec01Pass = "SEC01_SECRET_PASSWORD_zzqx"

// sec01Seeder seeds a project with basic-auth credentials in the base workspace
// (wID), which user uID belongs to.
func sec01Seeder(ctx context.Context, r *repo.Container) error {
	if err := baseSeeder(ctx, r); err != nil {
		return err
	}
	prj := project.New().NewID().
		Name("sec01").
		Workspace(wID).
		IsBasicAuthActive(true).
		BasicAuthUsername("sec01user").
		BasicAuthPassword(sec01Pass).
		MustBuild()
	return r.Project.Save(ctx, prj)
}

// TestSEC01_AnonymousCannotReadProjectBasicAuth is the SEC-01 regression test: an
// anonymous GraphQL caller must not be able to read a project's basic-auth
// credentials, while a workspace member still can. This exercises the full
// HTTP -> GraphQL -> repo stack; the deny-all filter for anonymous requests (see
// reposForOperator) is what closes the field-traversal bypass
// (plugin -> scene -> project -> basicAuthPassword) as well as the top-level query.
func TestSEC01_AnonymousCannotReadProjectBasicAuth(t *testing.T) {
	e := StartServer(t, &config.Config{
		Origins: []string{"https://example.com"},
		AuthSrv: config.AuthSrvConfig{Disabled: true},
	}, true, sec01Seeder)

	q := GraphQLRequest{
		OperationName: "Projects",
		Query:         "query Projects($teamId: ID!) { projects(teamId: $teamId, first: 50) { nodes { id basicAuthUsername basicAuthPassword } } }",
		Variables:     map[string]any{"teamId": wID.String()},
	}

	// anonymous: must NOT leak the credential.
	e.POST("/api/graphql").
		WithHeader("Origin", "https://example.com").
		WithHeader("Content-Type", "application/json").
		WithJSON(q).
		Expect().
		Status(http.StatusOK).
		Body().NotContains(sec01Pass)

	// authenticated workspace member: can still read it (no regression).
	e.POST("/api/graphql").
		WithHeader("Origin", "https://example.com").
		WithHeader("Content-Type", "application/json").
		WithHeader("X-Reearth-Debug-User", uID.String()).
		WithJSON(q).
		Expect().
		Status(http.StatusOK).
		Body().Contains(sec01Pass)
}
