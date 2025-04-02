package app

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/reearth/reearth/server/internal/adapter"
	"github.com/reearth/reearth/server/internal/usecase"
	"github.com/reearth/reearthx/account/accountdomain"
	"github.com/reearth/reearthx/account/accountdomain/user"
	"github.com/reearth/reearthx/account/accountdomain/workspace"
	"github.com/reearth/reearthx/account/accountusecase"
	"github.com/reearth/reearthx/account/accountusecase/accountinteractor"
	"github.com/reearth/reearthx/log"
	"github.com/reearth/reearthx/rerror"
	"github.com/reearth/reearthx/util"
)

type contextKey string

const (
	debugUserHeader            = "X-Reearth-Debug-User"
	contextUser     contextKey = "reearth_user"
)

// load user from db and attach it to context along with operator
// user id can be from debug header or jwt token
// if its new user, create new user and attach it to context
func attachOpMiddleware(cfg *ServerConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		multiUser := accountinteractor.NewMultiUser(cfg.AccountRepos, cfg.AccountGateways, cfg.Config.SignupSecret, cfg.Config.Host_Web, cfg.AccountRepos.Users)
		return func(c echo.Context) error {
			req := c.Request()
			ctx := req.Context()

			var userID string
			var u *user.User

			// get sub from context
			au := adapter.GetAuthInfo(ctx)
			if u, ok := ctx.Value(contextUser).(string); ok {
				userID = u
			}

			// debug mode
			if cfg.Debug {
				if userID := c.Request().Header.Get(debugUserHeader); userID != "" {
					if uId, err := accountdomain.UserIDFrom(userID); err == nil {
						user2, err := multiUser.FetchByID(ctx, user.IDList{uId})
						if err == nil && len(user2) == 1 {
							u = user2[0]
						}
					}
				}
			}

			if u == nil && userID != "" {
				if userID2, err := accountdomain.UserIDFrom(userID); err == nil {
					u2, err := multiUser.FetchByID(ctx, user.IDList{userID2})
					if err != nil {
						return err
					}
					if len(u2) > 0 {
						u = u2[0]
					}
				} else {
					return err
				}
			}

			if u == nil && au != nil {
				var err error
				// find user
				u, err = multiUser.FetchBySub(ctx, au.Sub)
				if err != nil && err != rerror.ErrNotFound {
					return err
				}
			}

			// save a new sub
			if u != nil && au != nil {
				if err := addAuth0SubToUser(ctx, u, user.AuthFrom(au.Sub), cfg); err != nil {
					return err
				}
			}

			if u != nil {
				ctx = adapter.AttachUser(ctx, u)
				log.Debugfc(ctx, "auth: user: id=%s name=%s email=%s", u.ID(), u.Name(), u.Email())

				log.Debugfc(ctx, "[middleware] generating operator...")
				op, err := generateOperator(ctx, cfg, u)
				if err != nil {
					log.Errorf("[middleware] generateOperator failed: %v", err)
					return err
				}
				log.Debugfc(ctx, "[middleware] operator generated: %+v", op)

				ctx = adapter.AttachOperator(ctx, op)
				log.Debugfc(ctx, "auth: op: %#v", op)
			}

			c.SetRequest(req.WithContext(ctx))
			return next(c)
		}
	}
}

// s
func generateOperator(ctx context.Context, cfg *ServerConfig, u *user.User) (*usecase.Operator, error) {
	log.Debugfc(ctx, "[generateOperator] start: user id=%s", u.ID())

	if u == nil {
		log.Debugfc(ctx, "[generateOperator] user is nil")
		return nil, nil
	}

	uid := u.ID()

	log.Debugfc(ctx, "[generateOperator] WorkspaceRepo type: %T", cfg.Repos.Workspace)
	log.Debugfc(ctx, "[generateOperator] SceneRepo type: %T", cfg.Repos.Scene)

	workspaces, err := cfg.Repos.Workspace.FindByUser(ctx, uid)
	if err != nil {
		log.Errorf("[generateOperator] failed to fetch workspaces: %v", err)
		return nil, err
	}
	if len(workspaces) == 0 {
		log.Debugfc(ctx, "[generateOperator] no workspaces found for user id=%s", uid)
	} else {
		log.Debugfc(ctx, "[generateOperator] fetched %d workspaces", len(workspaces))
	}

	scenes, err := cfg.Repos.Scene.FindByWorkspace(ctx, workspaces.IDs()...)
	if err != nil {
		log.Errorf("[generateOperator] failed to fetch scenes: %v", err)
		return nil, err
	}
	if len(scenes) == 0 {
		log.Debugfc(ctx, "[generateOperator] no scenes found for workspaces: %v", workspaces.IDs())
	} else {
		log.Debugfc(ctx, "[generateOperator] fetched %d scenes", len(scenes))
	}

	readableWorkspaces := workspaces.FilterByUserRole(uid, workspace.RoleReader).IDs()
	writableWorkspaces := workspaces.FilterByUserRole(uid, workspace.RoleWriter).IDs()
	maintainingWorkspaces := workspaces.FilterByUserRole(uid, workspace.RoleMaintainer).IDs()
	owningWorkspaces := workspaces.FilterByUserRole(uid, workspace.RoleOwner).IDs()
	defaultPolicy := util.CloneRef(cfg.Config.Policy.Default)

	log.Debugfc(ctx, "[generateOperator] readableWorkspaces: %v", readableWorkspaces)
	log.Debugfc(ctx, "[generateOperator] writableWorkspaces: %v", writableWorkspaces)
	log.Debugfc(ctx, "[generateOperator] maintainingWorkspaces: %v", maintainingWorkspaces)
	log.Debugfc(ctx, "[generateOperator] owningWorkspaces: %v", owningWorkspaces)

	log.Debugfc(ctx, "[generateOperator] readableScenes: %v", scenes.FilterByWorkspace(readableWorkspaces...).IDs())
	log.Debugfc(ctx, "[generateOperator] writableScenes: %v", scenes.FilterByWorkspace(writableWorkspaces...).IDs())
	log.Debugfc(ctx, "[generateOperator] maintainingScenes: %v", scenes.FilterByWorkspace(maintainingWorkspaces...).IDs())
	log.Debugfc(ctx, "[generateOperator] owningScenes: %v", scenes.FilterByWorkspace(owningWorkspaces...).IDs())

	return &usecase.Operator{
		AcOperator: &accountusecase.Operator{
			User:                   &uid,
			ReadableWorkspaces:     readableWorkspaces,
			WritableWorkspaces:     writableWorkspaces,
			MaintainableWorkspaces: maintainingWorkspaces,
			OwningWorkspaces:       owningWorkspaces,
			DefaultPolicy:          defaultPolicy,
		},
		ReadableScenes:    scenes.FilterByWorkspace(readableWorkspaces...).IDs(),
		WritableScenes:    scenes.FilterByWorkspace(writableWorkspaces...).IDs(),
		MaintainingScenes: scenes.FilterByWorkspace(maintainingWorkspaces...).IDs(),
		OwningScenes:      scenes.FilterByWorkspace(owningWorkspaces...).IDs(),
		DefaultPolicy:     defaultPolicy,
	}, nil
}

func addAuth0SubToUser(ctx context.Context, u *user.User, a user.Auth, cfg *ServerConfig) error {
	if u.AddAuth(a) {
		err := cfg.Repos.User.Save(ctx, u)
		if err != nil {
			return err
		}
	}
	return nil
}

func AuthRequiredMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			if adapter.Operator(ctx) == nil {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}
