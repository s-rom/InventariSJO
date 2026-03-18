package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	dbsqlc "inventari/api/internal/db/sqlc"
	"inventari/api/internal/middleware"
)

func newRequest(method, target string, body any, user dbsqlc.AppUser) (*httptest.ResponseRecorder, *http.Request) {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, target, &buf)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.CtxUser, user)
	req = req.WithContext(ctx)
	return httptest.NewRecorder(), req
}

func adminUser() dbsqlc.AppUser    { return dbsqlc.AppUser{AppUserID: 1, RoleID: "admin"} }
func editorUser() dbsqlc.AppUser   { return dbsqlc.AppUser{AppUserID: 2, RoleID: "editor"} }
func tutorUser() dbsqlc.AppUser    { return dbsqlc.AppUser{AppUserID: 3, RoleID: "tutor"} }
func readonlyUser() dbsqlc.AppUser { return dbsqlc.AppUser{AppUserID: 4, RoleID: "readonly"} }
