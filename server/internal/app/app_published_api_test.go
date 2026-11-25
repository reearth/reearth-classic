package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/reearth/reearth/server/internal/adapter"
	"github.com/reearth/reearth/server/internal/infrastructure/fs"
	"github.com/reearth/reearth/server/internal/infrastructure/memory"
	"github.com/reearth/reearth/server/internal/usecase/interactor"
	"github.com/reearth/reearth/server/internal/usecase/interfaces"
	"github.com/reearth/reearth/server/pkg/project"
	"github.com/samber/lo"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestPublishedAPIRouting_GatewayToken guards the wiring in publishedAPIRouting
// (used by initEcho). It fails if RequireGatewayToken is dropped from
// /api/published_data, or if /api/published stops stripping the basic-auth
// credentials. SEC-01/02/03.
func TestPublishedAPIRouting_GatewayToken(t *testing.T) {
	const (
		alias    = "alias"
		user     = "u53rn4m3"
		pass     = "p455w0rd-xyz"
		sceneTag = "SCENE_BODY_XYZ"
		token    = "secret-gateway-token"
	)

	prj := project.New().NewID().Workspace(project.NewWorkspaceID()).
		Alias(alias).
		PublishmentStatus(project.PublishmentStatusPublic).
		IsBasicAuthActive(true).
		BasicAuthUsername(user).
		BasicAuthPassword(pass).
		MustBuild()

	ctx := context.Background()
	mfs := afero.NewMemMapFs()
	prjRepo := memory.NewProject()
	storyRepo := memory.NewStorytelling()
	lo.Must0(prjRepo.Save(ctx, prj))
	fileg := lo.Must(fs.NewFile(mfs, ""))
	lo.Must0(fileg.UploadBuiltScene(ctx, strings.NewReader(`{"scene":"`+sceneTag+`"}`), alias))

	newEcho := func(gatewayToken string) *echo.Echo {
		e := echo.New()
		e.Use(ContextMiddleware(func(ctx context.Context) context.Context {
			return adapter.AttachUsecases(ctx, &interfaces.Container{
				Published: interactor.NewPublished(prjRepo, storyRepo, fileg, ""),
			})
		}))
		api := e.Group("/api")
		publishedAPIRouting(api, gatewayToken, "")
		return e
	}

	get := func(e *echo.Echo, path, header string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodGet, path, nil)
		if header != "" {
			r.Header.Set("X-Internal-Auth", header)
		}
		w := httptest.NewRecorder()
		e.ServeHTTP(w, r)
		return w
	}

	t.Run("token configured: /api/published_data is gated", func(t *testing.T) {
		e := newEcho(token)
		assert.Equal(t, http.StatusUnauthorized, get(e, "/api/published_data/"+alias, "").Code, "no token must be rejected")
		ok := get(e, "/api/published_data/"+alias, token)
		assert.Equal(t, http.StatusOK, ok.Code)
		assert.Contains(t, ok.Body.String(), sceneTag, "a valid token must pass through to the data")
	})

	t.Run("token configured: /api/published strips creds unless a valid token", func(t *testing.T) {
		e := newEcho(token)
		noTok := get(e, "/api/published/"+alias, "")
		assert.Equal(t, http.StatusOK, noTok.Code)
		assert.NotContains(t, noTok.Body.String(), pass, "creds must be stripped without a valid token")
		withTok := get(e, "/api/published/"+alias, token)
		assert.Contains(t, withTok.Body.String(), pass, "the trusted gateway must still receive the creds")
	})

	t.Run("no token configured: data fail-opens, metadata still strips creds", func(t *testing.T) {
		e := newEcho("")
		assert.Equal(t, http.StatusOK, get(e, "/api/published_data/"+alias, "").Code, "self-hosted (no token) must stay open")
		meta := get(e, "/api/published/"+alias, "")
		assert.Equal(t, http.StatusOK, meta.Code)
		assert.NotContains(t, meta.Body.String(), pass, "metadata strips creds even with no token configured")
	})
}
