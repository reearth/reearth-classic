package adapter

import (
	"context"

	"github.com/reearth/reearth/server/internal/usecase"
	"github.com/reearth/reearth/server/internal/usecase/interfaces"
	"github.com/reearth/reearthx/account/accountdomain/user"
	"github.com/reearth/reearthx/account/accountusecase"
	"github.com/reearth/reearthx/appx"
	"github.com/reearth/reearthx/log"
	"golang.org/x/text/language"
)

type ContextKey string

const (
	contextUser     ContextKey = "user"
	contextOperator ContextKey = "operator"
	ContextAuthInfo ContextKey = "authinfo"
	contextUsecases ContextKey = "usecases"
)

var defaultLang = language.English

type AuthInfo struct {
	Token         string
	Sub           string
	Iss           string
	Name          string
	Email         string
	EmailVerified *bool
}

func AttachUser(ctx context.Context, u *user.User) context.Context {
	return context.WithValue(ctx, contextUser, u)
}

func AttachOperator(ctx context.Context, o *usecase.Operator) context.Context {
	log.Debugfc(ctx, "[AttachOperator] Attaching operator: %+v", o)
	return context.WithValue(ctx, contextOperator, o)
}

func AttachUsecases(ctx context.Context, u *interfaces.Container) context.Context {
	ctx = context.WithValue(ctx, contextUsecases, u)
	return ctx
}

func User(ctx context.Context) *user.User {
	log.Debugfc(ctx, "[User] Start User")
	if v := ctx.Value(contextUser); v != nil {
		log.Debugfc(ctx, "[User] contextUser found in context")
		if u, ok := v.(*user.User); ok {
			log.Debugfc(ctx, "[User] contextUser is of type *user.User: %+v", u)
			return u
		} else {
			log.Debugfc(ctx, "[User] contextUser is not of type *user.User")
		}
	} else {
		log.Debugfc(ctx, "[User] contextUser not found in context")
	}
	return nil
}

func Lang(ctx context.Context, lang *language.Tag) string {
	if lang != nil && !lang.IsRoot() {
		return lang.String()
	}

	u := User(ctx)
	if u == nil {
		return defaultLang.String()
	}

	l := u.Lang()
	if l.IsRoot() {
		return defaultLang.String()
	}

	return l.String()
}

func Operator(ctx context.Context) *usecase.Operator {
	log.Debugfc(ctx, "[Operator] Start Operator")
	if v := ctx.Value(contextOperator); v != nil {
		log.Debugfc(ctx, "[Operator] contextOperator found in context")
		if v2, ok := v.(*usecase.Operator); ok {
			log.Debugfc(ctx, "[Operator] contextOperator is of type *usecase.Operator: %+v", v2)
			return v2
		} else {
			log.Debugfc(ctx, "[Operator] contextOperator is not of type *usecase.Operator")
		}
	} else {
		log.Debugfc(ctx, "[Operator] contextOperator not found in context")
	}
	return nil
}

func AcOperator(ctx context.Context) *accountusecase.Operator {
	if v := ctx.Value(contextOperator); v != nil {
		if v2, ok := v.(*accountusecase.Operator); ok {
			return v2
		}
	}
	return nil
}

func GetAuthInfo(ctx context.Context) *appx.AuthInfo {
	if v := ctx.Value(ContextAuthInfo); v != nil {
		if v2, ok := v.(appx.AuthInfo); ok {
			return &v2
		}
	}
	return nil
}

func Usecases(ctx context.Context) *interfaces.Container {
	return ctx.Value(contextUsecases).(*interfaces.Container)
}
