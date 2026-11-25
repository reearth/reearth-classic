package app

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/reearth/reearth/server/internal/adapter"
	"github.com/reearth/reearth/server/internal/usecase"
	"github.com/reearth/reearth/server/internal/usecase/gateway"
	"github.com/reearth/reearth/server/internal/usecase/interactor"
	"github.com/reearth/reearth/server/internal/usecase/repo"
	"github.com/reearth/reearth/server/pkg/scene"
	"github.com/reearth/reearthx/account/accountdomain"
	"github.com/reearth/reearthx/account/accountusecase/accountgateway"
	"github.com/reearth/reearthx/account/accountusecase/accountrepo"
)

// reposForOperator returns the repo container filtered for the request's operator.
// With an operator, reads are scoped to its readable/writable workspaces and scenes.
// Without one (an anonymous request) it applies a DENY-ALL filter -- empty, non-nil
// lists -- so every workspace/scene-scoped read returns not-found. This is the
// SEC-01 fix: WorkspaceFilter/SceneFilter.CanRead treat a nil Readable as "allow
// all", so leaving anonymous repos unfiltered let an anonymous GraphQL caller read
// any project/scene/plugin (e.g. plugin(id) -> scene -> project -> basicAuthPassword,
// bypassing the per-resolver requireUser guards). Public lookups used by the
// published pages (FindByPublicName, filterByWorkspaces=false) do not consult this
// filter and keep working.
func reposForOperator(r *repo.Container, op *usecase.Operator) *repo.Container {
	if op != nil {
		return r.Filtered(repo.WorkspaceFilterFromOperator(op), repo.SceneFilterFromOperator(op))
	}
	return r.Filtered(
		repo.WorkspaceFilter{Readable: accountdomain.WorkspaceIDList{}, Writable: accountdomain.WorkspaceIDList{}},
		repo.SceneFilter{Readable: scene.IDList{}, Writable: scene.IDList{}},
	)
}

func UsecaseMiddleware(r *repo.Container, g *gateway.Container, ar *accountrepo.Container, ag *accountgateway.Container, config interactor.ContainerConfig) echo.MiddlewareFunc {
	return ContextMiddleware(func(ctx context.Context) context.Context {
		repos := r

		repos = reposForOperator(repos, adapter.Operator(ctx))

		var ar2 *accountrepo.Container
		if op := adapter.AcOperator(ctx); op != nil && ar != nil {
			// Only apply filtering if operator has workspaces
			// This prevents circular dependency during initial workspace fetch (GetMe query)
			// where the operator needs workspaces but can't get them because the filter blocks them
			if len(op.AllReadableWorkspaces()) > 0 {
				ar2 = ar.Filtered(accountrepo.WorkspaceFilterFromOperator(op))
			} else {
				ar2 = ar
			}
		} else {
			ar2 = ar
		}

		uc := interactor.NewContainer(repos, g, ar2, ag, config)
		ctx = adapter.AttachUsecases(ctx, &uc)
		return ctx
	})
}

func ContextMiddleware(fn func(ctx context.Context) context.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			c.SetRequest(req.WithContext(fn(req.Context())))
			return next(c)
		}
	}
}
