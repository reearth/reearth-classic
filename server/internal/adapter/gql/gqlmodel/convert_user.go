package gqlmodel

import (
	"log"

	"github.com/reearth/reearthx/account/accountdomain/user"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
)

func ToUser(u *user.User) *User {
	if u == nil {
		return nil
	}

	return &User{
		ID:    IDFrom(u.ID()),
		Name:  u.Name(),
		Email: u.Email(),
		Host:  lo.EmptyableToPtr(u.Host()),
	}
}

func ToUsersFromSimple(users user.SimpleList) []*User {
	return util.Map(users, func(u *user.Simple) *User {
		return ToUserFromSimple(u)
	})
}

func ToUserFromSimple(u *user.Simple) *User {
	if u == nil {
		return nil
	}

	return &User{
		ID:    IDFrom(u.ID),
		Name:  u.Name,
		Email: u.Email,
	}
}

func ToMe(u *user.User) *Me {
	log.Printf("[ToMe] Start ToMe")
	if u == nil {
		log.Printf("[ToMe] user is nil")
		return nil
	}

	me := &Me{
		ID:       IDFrom(u.ID()),
		Name:     u.Name(),
		Email:    u.Email(),
		Lang:     u.Lang(),
		Theme:    Theme(u.Theme()),
		MyTeamID: IDFrom(u.Workspace()),
		Auths: util.Map(u.Auths(), func(a user.Auth) string {
			return a.Provider
		}),
	}
	log.Printf("[ToMe] Created Me object: %+v", me)
	return me
}

func ToTheme(t *Theme) *user.Theme {
	if t == nil {
		return nil
	}

	th := user.ThemeDefault
	switch *t {
	case ThemeDark:
		th = user.ThemeDark
	case ThemeLight:
		th = user.ThemeLight
	}
	return &th
}
