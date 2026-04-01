package handler_test

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"testing"

	dbsqlc "inventari/api/internal/db/sqlc"
	"inventari/api/internal/handler"
	testmock "inventari/api/internal/testutil/mock"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─── List ─────────────────────────────────────────────────────────────────────

func TestRoles_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRolesHandler(mq, slog.Default())

	mq.On("ListRoles", mock.Anything).Return([]dbsqlc.Role{
		{RoleID: "admin"},
		{RoleID: "editor"},
	}, nil)

	w, r := newRequest("GET", "/roles", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.Role
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 2)
	mq.AssertExpectations(t)
}

func TestRoles_List_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRolesHandler(mq, slog.Default())

	mq.On("ListRoles", mock.Anything).Return([]dbsqlc.Role(nil), errors.New("db down"))

	w, r := newRequest("GET", "/roles", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Create ───────────────────────────────────────────────────────────────────

func TestRoles_Create_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRolesHandler(mq, slog.Default())

	mq.On("CreateRole", mock.Anything, mock.AnythingOfType("CreateRoleParams")).
		Return(dbsqlc.Role{RoleID: "supervisor"}, nil)

	w, r := newRequest("POST", "/roles", map[string]any{"role_id": "supervisor"}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestRoles_Create_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRolesHandler(mq, slog.Default())

	w, r := newRequest("POST", "/roles", map[string]any{"role_id": "supervisor"}, editorUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mq.AssertNotCalled(t, "CreateRole", mock.Anything, mock.Anything)
}

func TestRoles_Create_EmptyRoleID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRolesHandler(mq, slog.Default())

	w, r := newRequest("POST", "/roles", map[string]any{"role_id": ""}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── Delete ───────────────────────────────────────────────────────────────────

func TestRoles_Delete_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRolesHandler(mq, slog.Default())

	mq.On("DeleteRole", mock.Anything, "custom").Return(nil)

	w, r := newRequest("DELETE", "/roles/custom", nil, adminUser())
	r.SetPathValue("id", "custom")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}

func TestRoles_Delete_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRolesHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/roles/custom", nil, tutorUser())
	r.SetPathValue("id", "custom")
	h.Delete(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mq.AssertNotCalled(t, "DeleteRole", mock.Anything, mock.Anything)
}

func TestRoles_Delete_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRolesHandler(mq, slog.Default())

	mq.On("DeleteRole", mock.Anything, "ghost").Return(pgx.ErrNoRows)

	w, r := newRequest("DELETE", "/roles/ghost", nil, adminUser())
	r.SetPathValue("id", "ghost")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRoles_Delete_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRolesHandler(mq, slog.Default())

	mq.On("DeleteRole", mock.Anything, "custom").Return(errors.New("db error"))

	w, r := newRequest("DELETE", "/roles/custom", nil, adminUser())
	r.SetPathValue("id", "custom")
	h.Delete(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
