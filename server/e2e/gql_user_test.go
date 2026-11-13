package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/reearth/reearth/server/internal/app/config"
	"github.com/reearth/reearth/server/internal/usecase/repo"
	"github.com/reearth/reearthx/account/accountdomain"
	"github.com/reearth/reearthx/account/accountdomain/user"
	"github.com/reearth/reearthx/account/accountdomain/workspace"
	"github.com/reearth/reearthx/idx"
	"github.com/reearth/reearthx/rerror"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

var (
	uId1 = accountdomain.NewUserID()
	uId2 = accountdomain.NewUserID()
	uId3 = accountdomain.NewUserID()
	wId1 = accountdomain.NewWorkspaceID()
	wId2 = accountdomain.NewWorkspaceID()
	iId1 = accountdomain.NewIntegrationID()
)

func baseSeederUser(ctx context.Context, r *repo.Container) error {
	auth := user.ReearthSub(uId1.String())
	u := user.New().ID(uId1).
		Name("e2e").
		Email("e2e@e2e.com").
		Auths([]user.Auth{*auth}).
		Theme(user.ThemeDark).
		Lang(language.Japanese).
		Workspace(wId1).
		MustBuild()
	if err := r.User.Save(ctx, u); err != nil {
		return err
	}
	u2 := user.New().ID(uId2).
		Name("e2e2").
		Workspace(wId2).
		Email("e2e2@e2e.com").
		MustBuild()
	if err := r.User.Save(ctx, u2); err != nil {
		return err
	}
	u3 := user.New().ID(uId3).
		Name("e2e3").
		Workspace(wId2).
		Email("e2e3@e2e.com").
		MustBuild()
	if err := r.User.Save(ctx, u3); err != nil {
		return err
	}
	roleOwner := workspace.Member{
		Role:      workspace.RoleOwner,
		InvitedBy: uId1,
	}
	roleReader := workspace.Member{
		Role:      workspace.RoleReader,
		InvitedBy: uId2,
	}

	w := workspace.New().ID(wId1).
		Name("e2e").
		Members(map[idx.ID[accountdomain.User]]workspace.Member{
			uId1: roleOwner,
		}).
		Integrations(map[idx.ID[accountdomain.Integration]]workspace.Member{
			iId1: roleOwner,
		}).
		MustBuild()
	if err := r.Workspace.Save(ctx, w); err != nil {
		return err
	}

	w2 := workspace.New().ID(wId2).
		Name("e2e2").
		Members(map[idx.ID[accountdomain.User]]workspace.Member{
			uId1: roleOwner,
			uId3: roleReader,
		}).
		Integrations(map[idx.ID[accountdomain.Integration]]workspace.Member{
			iId1: roleOwner,
		}).
		MustBuild()
	if err := r.Workspace.Save(ctx, w2); err != nil {
		return err
	}

	return nil
}

// func TestSignUp(t *testing.T) {
// 	e, _ := StartGQLServer(t, &config.Config{
// 		Origins: []string{"https://example.com"},
// 		AuthSrv: config.AuthSrvConfig{
// 			Disabled: true,
// 		},
// 	}, true, baseSeederUser)
// 	query := `mutation { signup(input: {lang: "ja",theme: DEFAULT,secret: "Ajsownndww1"}){ user{ id name email } }}`
// 	request := GraphQLRequest{
// 		Query: query,
// 	}
// 	jsonData, err := json.Marshal(request)
// 	if err != nil {
// 		assert.NoError(t, err)
// 	}
// 	o := e.POST("/api/graphql").
// 		WithHeader("authorization", "Bearer test").
// 		WithHeader("Content-Type", "application/json").
// 		WithHeader("X-Reearth-Debug-User", uId1.String()).
// 		WithBytes(jsonData).Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().Value("signup").Object().Value("user").Object()
// 	o.Value("name").String().Equal("updated")
// 	o.Value("email").String().Equal("hoge@test.com")
// }

func TestUpdateMe(t *testing.T) {
	e, _ := StartGQLServer(t, &config.Config{
		Origins: []string{"https://example.com"},
		AuthSrv: config.AuthSrvConfig{
			Disabled: true,
		},
	}, true, baseSeederUser)
	query := `mutation { updateMe(input: {name: "updated",email:"hoge@test.com",lang: "ja",theme: DEFAULT,password: "Ajsownndww1",passwordConfirmation: "Ajsownndww1"}){ me{ id name email lang theme } }}`
	request := GraphQLRequest{
		Query: query,
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		assert.NoError(t, err)
	}
	o := e.POST("/api/graphql").
		WithHeader("authorization", "Bearer test").
		WithHeader("Content-Type", "application/json").
		WithHeader("X-Reearth-Debug-User", uId1.String()).
		WithBytes(jsonData).Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().Value("updateMe").Object().Value("me").Object()
	o.Value("name").String().Equal("updated")
	o.Value("email").String().Equal("hoge@test.com")
	o.Value("lang").String().Equal("ja")
	o.Value("theme").String().Equal("default")
}

func TestRemoveMyAuth(t *testing.T) {
	e, r := StartGQLServer(t, &config.Config{
		Origins: []string{"https://example.com"},
		AuthSrv: config.AuthSrvConfig{
			Disabled: true,
		},
	}, true, baseSeederUser)
	u, err := r.User.FindByID(context.Background(), uId1)
	assert.Nil(t, err)
	assert.Equal(t, &user.Auth{Provider: "reearth", Sub: "reearth|" + uId1.String()}, u.Auths().GetByProvider("reearth"))

	query := `mutation { removeMyAuth(input: {auth: "reearth"}){ me{ id name email lang theme } }}`
	request := GraphQLRequest{
		Query: query,
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		assert.NoError(t, err)
	}
	e.POST("/api/graphql").
		WithHeader("authorization", "Bearer test").
		WithHeader("Content-Type", "application/json").
		WithHeader("X-Reearth-Debug-User", uId1.String()).
		WithBytes(jsonData).Expect().Status(http.StatusOK).JSON().Object()

	u, err = r.User.FindByID(context.Background(), uId1)
	assert.Nil(t, err)
	assert.Nil(t, u.Auths().Get("sub"))
}

func TestDeleteMe(t *testing.T) {
	e, r := StartGQLServer(t, &config.Config{
		Origins: []string{"https://example.com"},
		AuthSrv: config.AuthSrvConfig{
			Disabled: true,
		},
	}, true, baseSeederUser)
	u, err := r.User.FindByID(context.Background(), uId1)
	assert.Nil(t, err)
	assert.NotNil(t, u)

	query := fmt.Sprintf(`mutation { deleteMe(input: {userId: "%s"}){ userId }}`, uId1)
	request := GraphQLRequest{
		Query: query,
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		assert.NoError(t, err)
	}
	e.POST("/api/graphql").
		WithHeader("authorization", "Bearer test").
		WithHeader("Content-Type", "application/json").
		WithHeader("X-Reearth-Debug-User", uId1.String()).
		WithBytes(jsonData).Expect().Status(http.StatusOK).JSON().Object()

	_, err = r.User.FindByID(context.Background(), uId1)
	assert.Equal(t, rerror.ErrNotFound, err)
}

func TestSearchUser(t *testing.T) {
	mixedCaseID := accountdomain.NewUserID()
	e, _ := StartGQLServer(t, &config.Config{
		Origins: []string{"https://example.com"},
		AuthSrv: config.AuthSrvConfig{
			Disabled: true,
		},
	}, true, func(ctx context.Context, r *repo.Container) error {
		if err := baseSeederUser(ctx, r); err != nil {
			return err
		}
		return r.User.Save(ctx, user.New().ID(mixedCaseID).
			Name("mixed").
			Email("Mixed.Case@E2E.com").
			Workspace(accountdomain.NewWorkspaceID()).
			MustBuild())
	})

	search := func(keyword string, asUser string) *httpexpect.Object {
		req := e.POST("/api/graphql").
			WithHeader("Content-Type", "application/json").
			WithJSON(GraphQLRequest{Query: fmt.Sprintf(`{ searchUser(nameOrEmail: %q){ id name email } }`, keyword)})
		if asUser != "" {
			req = req.WithHeader("authorization", "Bearer test").
				WithHeader("X-Reearth-Debug-User", asUser)
		}
		return req.Expect().Status(http.StatusOK).JSON().Object()
	}

	anon := search("e2e2@e2e.com", "")
	anon.Value("data").Object().Value("searchUser").Null()
	anon.Value("errors").Array().Element(0).Object().Value("message").String().Contains("unauthorized")

	o := search("e2e2@e2e.com", uId1.String()).Value("data").Object().Value("searchUser").Object()
	o.Value("id").String().Equal(uId2.String())
	o.Value("name").String().Equal("e2e2")
	o.Value("email").String().Equal("")

	search(" E2E2@e2e.com ", uId1.String()).Value("data").Object().Value("searchUser").Object().
		Value("id").String().Equal(uId2.String())

	self := search("e2e@e2e.com", uId1.String()).Value("data").Object().Value("searchUser").Object()
	self.Value("id").String().Equal(uId1.String())
	self.Value("email").String().Equal("e2e@e2e.com")

	for _, kw := range []string{"mixed.case@e2e.com", "MIXED.CASE@E2E.COM", "Mixed.Case@E2E.com"} {
		m := search(kw, uId1.String()).Value("data").Object().Value("searchUser").Object()
		m.Value("id").String().Equal(mixedCaseID.String())
		m.Value("email").String().Equal("")
	}

	for _, kw := range []string{"e2e2@", "2@e2e", "e2e.com", "E2E2", "e2e2", "e2e", "notfound", "mixed.case@", "ixed.case@e2e.com", "mixed.case@e2e.co", ".*@e2e.com"} {
		search(kw, uId1.String()).Value("data").Object().Value("searchUser").Null()
	}
}

func TestUserEmailVisibility(t *testing.T) {
	e, _ := StartGQLServer(t, &config.Config{
		Origins: []string{"https://example.com"},
		AuthSrv: config.AuthSrvConfig{
			Disabled: true,
		},
	}, true, baseSeederUser)

	query := func(q string, asUser string) *httpexpect.Object {
		req := e.POST("/api/graphql").
			WithHeader("Content-Type", "application/json").
			WithJSON(GraphQLRequest{Query: q})
		if asUser != "" {
			req = req.WithHeader("authorization", "Bearer test").
				WithHeader("X-Reearth-Debug-User", asUser)
		}
		return req.Expect().Status(http.StatusOK).JSON().Object()
	}

	nodeQ := func(id string) string {
		return fmt.Sprintf(`{ node(id: %q, type: USER){ id ... on User { name email } } }`, id)
	}
	nodesQ := func(id string) string {
		return fmt.Sprintf(`{ nodes(id: [%q], type: USER){ id ... on User { name email } } }`, id)
	}

	anon := query(nodeQ(uId2.String()), "")
	anon.Value("data").Object().Value("node").Null()
	anon.Value("errors").Array().Element(0).Object().Value("message").String().Contains("unauthorized")
	query(nodesQ(uId2.String()), "").Value("errors").Array().Element(0).Object().
		Value("message").String().Contains("unauthorized")

	o := query(nodeQ(uId2.String()), uId1.String()).Value("data").Object().Value("node").Object()
	o.Value("id").String().Equal(uId2.String())
	o.Value("email").String().Equal("")
	query(nodesQ(uId2.String()), uId1.String()).Value("data").Object().Value("nodes").Array().
		Element(0).Object().Value("email").String().Equal("")

	query(nodeQ(uId1.String()), uId1.String()).Value("data").Object().Value("node").Object().
		Value("email").String().Equal("e2e@e2e.com")

	members := query(`{ me { teams { members { userId user { email } } } } }`, uId1.String()).
		Value("data").Object().Value("me").Object().Value("teams").Array()
	for _, team := range members.Iter() {
		for _, m := range team.Object().Value("members").Array().Iter() {
			mo := m.Object()
			email := mo.Value("user").Object().Value("email").String()
			if mo.Value("userId").String().Raw() == uId1.String() {
				email.Equal("e2e@e2e.com")
			} else {
				email.Equal("")
			}
		}
	}
}

func TestNode(t *testing.T) {
	e, _ := StartGQLServer(t, &config.Config{
		Origins: []string{"https://example.com"},
		AuthSrv: config.AuthSrvConfig{
			Disabled: true,
		},
	}, true, baseSeederUser)
	query := fmt.Sprintf(` { node(id: "%s", type: USER){ id } }`, uId1.String())
	request := GraphQLRequest{
		Query: query,
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		assert.NoError(t, err)
	}
	o := e.POST("/api/graphql").
		WithHeader("authorization", "Bearer test").
		WithHeader("Content-Type", "application/json").
		WithHeader("X-Reearth-Debug-User", uId1.String()).
		WithBytes(jsonData).Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().Value("node").Object()
	o.Value("id").String().Equal(uId1.String())
}

func TestNodes(t *testing.T) {
	e, _ := StartGQLServer(t, &config.Config{
		Origins: []string{"https://example.com"},
		AuthSrv: config.AuthSrvConfig{
			Disabled: true,
		},
	}, true, baseSeederUser)
	query := fmt.Sprintf(` { nodes(id: "%s", type: USER){ id } }`, uId1.String())
	request := GraphQLRequest{
		Query: query,
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		assert.NoError(t, err)
	}
	o := e.POST("/api/graphql").
		WithHeader("authorization", "Bearer test").
		WithHeader("Content-Type", "application/json").
		WithHeader("X-Reearth-Debug-User", uId1.String()).
		WithBytes(jsonData).Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().Value("nodes")
	o.Array().Contains(map[string]string{"id": uId1.String()})
}
