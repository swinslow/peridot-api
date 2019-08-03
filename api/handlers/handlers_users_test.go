// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"net/http"
	"testing"

	"github.com/swinslow/peridot-db/pkg/datastore"
	hu "github.com/swinslow/peridot-api/test/handlerutils"
)

// ===== GET /users =====

func TestCanGetUsersHandlerAsAdmin(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users", "", "admin")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersHandler), "/users")
	hu.ConfirmOKResponse(t, rec)

	// expect full user data b/c we're an admin
	wanted := `{"users": [{"id": 1, "name": "Admin", "github": "admin", "access": "admin"}, {"id": 2, "name": "Operator", "github": "operator", "access": "operator"}, {"id": 3, "name": "Commenter", "github": "commenter", "access": "commenter"}, {"id": 4, "name": "Viewer", "github": "viewer", "access": "viewer"}, {"id": 10, "name": "Disabled", "github": "disabled", "access": "disabled"}]}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCanGetUsersHandlerAsOtherUsers(t *testing.T) {
	// should be same return for all
	wanted := `{"users": [{"id": 1, "github": "admin"}, {"id": 2, "github": "operator"}, {"id": 3, "github": "commenter"}, {"id": 4, "github": "viewer"}, {"id": 10, "github": "disabled"}]}`

	// as operator
	rec, req, env := setupTestEnv(t, "GET", "/users", "", "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersHandler), "/users")
	hu.ConfirmOKResponse(t, rec)
	hu.CheckResponse(t, rec, wanted)

	// as commenter
	rec, req, env = setupTestEnv(t, "GET", "/users", "", "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersHandler), "/users")
	hu.ConfirmOKResponse(t, rec)
	hu.CheckResponse(t, rec, wanted)

	// as viewer
	rec, req, env = setupTestEnv(t, "GET", "/users", "", "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersHandler), "/users")
	hu.ConfirmOKResponse(t, rec)
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetUsersHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users", "", "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersHandler), "/users")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/users", "", "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersHandler), "/users")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== POST /users =====

func TestCanPostUsersHandlerAsAdmin(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/users", `{"name": "Steve", "github": "swinslow", "access": "operator"}`, "admin")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersHandler), "/users")
	hu.ConfirmCreatedResponse(t, rec)

	wanted := `{"id": 11}`
	hu.CheckResponse(t, rec, wanted)

	// and verify state of database now
	users, err := env.db.GetAllUsers()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(users) != 6 {
		t.Errorf("expected %d, got %d", 6, len(users))
	}
	newUser, err := env.db.GetUserByID(11)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	wantedUser := &datastore.User{ID: 11, Name: "Steve", Github: "swinslow", AccessLevel: datastore.AccessOperator}
	if newUser.ID != wantedUser.ID || newUser.Name != wantedUser.Name || newUser.Github != wantedUser.Github || newUser.AccessLevel != wantedUser.AccessLevel {
		t.Errorf("expected %#v, got %#v", wantedUser, newUser)
	}
}

func TestCannotPostUsersHandlerAsOtherUser(t *testing.T) {
	// as operator
	rec, req, env := setupTestEnv(t, "POST", "/users", `{"name": "Steve", "github": "swinslow", "access": "operator"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersHandler), "/users")
	hu.ConfirmAccessDenied(t, rec)

	// as commenter
	rec, req, env = setupTestEnv(t, "POST", "/users", `{"name": "Steve", "github": "swinslow", "access": "commenter"}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersHandler), "/users")
	hu.ConfirmAccessDenied(t, rec)

	// as viewer
	rec, req, env = setupTestEnv(t, "POST", "/users", `{"name": "Steve", "github": "swinslow", "access": "viewer"}`, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersHandler), "/users")
	hu.ConfirmAccessDenied(t, rec)
}

// ===== GET /users/3 =====

func TestCanGetUsersOneHandlerAsAdmin(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users/3", "", "admin")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmOKResponse(t, rec)

	// expect full user data b/c we're an admin
	wanted := `{"user": {"id": 3, "name": "Commenter", "github": "commenter", "access": "commenter"}}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCanGetUsersOneHandlerAsOtherUsers(t *testing.T) {
	// should be same return for all EXCEPT self
	wanted := `{"user": {"id": 3, "github": "commenter"}}`

	// as operator
	rec, req, env := setupTestEnv(t, "GET", "/users/3", "", "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmOKResponse(t, rec)
	hu.CheckResponse(t, rec, wanted)

	// as viewer
	rec, req, env = setupTestEnv(t, "GET", "/users/3", "", "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmOKResponse(t, rec)
	hu.CheckResponse(t, rec, wanted)

	// and for commenter getting somebody else's data
	wanted = `{"user": {"id": 1, "github": "admin"}}`
	rec, req, env = setupTestEnv(t, "GET", "/users/1", "", "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmOKResponse(t, rec)
	hu.CheckResponse(t, rec, wanted)
}

func TestCanGetUsersOneHandlerAsSelf(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users/3", "", "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmOKResponse(t, rec)

	// expect full user data b/c we're getting our own data
	wanted := `{"user": {"id": 3, "name": "Commenter", "github": "commenter", "access": "commenter"}}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetUsersOneHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/users/3", "", "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/users/3", "", "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== PUT /users/3 =====

func TestCanPutUsersOneHandlerAsAdmin(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/users/3", `{"name": "new-name", "github": "new-github", "access": "operator"}`, "admin")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	u, err := env.db.GetUserByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedUser := &datastore.User{ID: 3, Name: "new-name", Github: "new-github", AccessLevel: datastore.AccessOperator}
	if u.ID != wantedUser.ID || u.Name != wantedUser.Name || u.Github != wantedUser.Github || u.AccessLevel != wantedUser.AccessLevel {
		t.Errorf("expected %#v, got %#v", wantedUser, u)
	}
}

func TestCanPutUsersOneHandlerAsAdminWithJustName(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/users/3", `{"name": "new-name"}`, "admin")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	u, err := env.db.GetUserByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedUser := &datastore.User{ID: 3, Name: "new-name", Github: "commenter", AccessLevel: datastore.AccessCommenter}
	if u.ID != wantedUser.ID || u.Name != wantedUser.Name || u.Github != wantedUser.Github || u.AccessLevel != wantedUser.AccessLevel {
		t.Errorf("expected %#v, got %#v", wantedUser, u)
	}
}

func TestCanPutUsersOneHandlerAsAdminWithJustGithub(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/users/3", `{"github": "new-github"}`, "admin")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	u, err := env.db.GetUserByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedUser := &datastore.User{ID: 3, Name: "Commenter", Github: "new-github", AccessLevel: datastore.AccessCommenter}
	if u.ID != wantedUser.ID || u.Name != wantedUser.Name || u.Github != wantedUser.Github || u.AccessLevel != wantedUser.AccessLevel {
		t.Errorf("expected %#v, got %#v", wantedUser, u)
	}
}

func TestCanPutUsersOneHandlerAsAdminWithJustAccessLevel(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/users/3", `{"access": "operator"}`, "admin")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	u, err := env.db.GetUserByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedUser := &datastore.User{ID: 3, Name: "Commenter", Github: "commenter", AccessLevel: datastore.AccessOperator}
	if u.ID != wantedUser.ID || u.Name != wantedUser.Name || u.Github != wantedUser.Github || u.AccessLevel != wantedUser.AccessLevel {
		t.Errorf("expected %#v, got %#v", wantedUser, u)
	}
}

func TestCannotPutUsersOneHandlerAsAdminWithInvalidAccessLevel(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/users/3", `{"access": "oops"}`, "admin")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmBadRequestResponse(t, rec)
}

func TestCanPutUsersOneHandlerAsOperatorSelfWithJustName(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/users/2", `{"name": "new-operator-name"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	u, err := env.db.GetUserByID(2)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedUser := &datastore.User{ID: 2, Name: "new-operator-name", Github: "operator", AccessLevel: datastore.AccessOperator}
	if u.ID != wantedUser.ID || u.Name != wantedUser.Name || u.Github != wantedUser.Github || u.AccessLevel != wantedUser.AccessLevel {
		t.Errorf("expected %#v, got %#v", wantedUser, u)
	}
}

func TestCannotPutUsersOneHandlerAsOperatorForOther(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/users/3", `{"name": "new-operator-name"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)
}

func TestCannotPutUsersOneHandlerAsOperatorForSelfOtherThanName(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/users/2", `{"github": "new-operator-gh"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "PUT", "/users/2", `{"access": "admin"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "PUT", "/users/2", `{"github": "new-operator-gh", "access": "admin"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "PUT", "/users/2", `{"name": "new-operator-name", "github": "new-operator-gh", "access": "admin"}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)
}

func TestCanPutUsersOneHandlerAsCommenterSelfWithJustName(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/users/3", `{"name": "new-commenter-name"}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	u, err := env.db.GetUserByID(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedUser := &datastore.User{ID: 3, Name: "new-commenter-name", Github: "commenter", AccessLevel: datastore.AccessCommenter}
	if u.ID != wantedUser.ID || u.Name != wantedUser.Name || u.Github != wantedUser.Github || u.AccessLevel != wantedUser.AccessLevel {
		t.Errorf("expected %#v, got %#v", wantedUser, u)
	}
}

func TestCannotPutUsersOneHandlerAsCommenterForOther(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/users/2", `{"name": "new-commenter-name"}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)
}

func TestCannotPutUsersOneHandlerAsCommenterForSelfOtherThanName(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/users/3", `{"github": "new-commenter-gh"}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "PUT", "/users/3", `{"access": "admin"}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "PUT", "/users/3", `{"github": "new-commenter-gh", "access": "admin"}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "PUT", "/users/3", `{"name": "new-commenter-name", "github": "new-commenter-gh", "access": "admin"}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)
}

func TestCanPutUsersOneHandlerAsViewerSelfWithJustName(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/users/4", `{"name": "new-viewer-name"}`, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmNoContentResponse(t, rec)

	// and verify state of database now
	u, err := env.db.GetUserByID(4)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantedUser := &datastore.User{ID: 4, Name: "new-viewer-name", Github: "viewer", AccessLevel: datastore.AccessViewer}
	if u.ID != wantedUser.ID || u.Name != wantedUser.Name || u.Github != wantedUser.Github || u.AccessLevel != wantedUser.AccessLevel {
		t.Errorf("expected %#v, got %#v", wantedUser, u)
	}
}

func TestCannotPutUsersOneHandlerAsViewerForOther(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/users/2", `{"name": "new-viewer-name"}`, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)
}

func TestCannotPutUsersOneHandlerAsViewerForSelfOtherThanName(t *testing.T) {
	rec, req, env := setupTestEnv(t, "PUT", "/users/4", `{"github": "new-viewer-gh"}`, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "PUT", "/users/4", `{"access": "admin"}`, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "PUT", "/users/4", `{"github": "new-viewer-gh", "access": "admin"}`, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "PUT", "/users/4", `{"name": "new-viewer-name", "github": "new-viewer-gh", "access": "admin"}`, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)
}

func TestCannotPutUsersOneHandlerAsBadUser(t *testing.T) {
	// disabled, for self
	rec, req, env := setupTestEnv(t, "PUT", "/users/10", `{"name": "new-name"}`, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)

	// disabled, for other
	rec, req, env = setupTestEnv(t, "PUT", "/users/3", `{"name": "new-name"}`, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmAccessDenied(t, rec)

	// invalid
	rec, req, env = setupTestEnv(t, "PUT", "/users/3", `{"name": "new-name"}`, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.usersOneHandler), "/users/{id}")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}
