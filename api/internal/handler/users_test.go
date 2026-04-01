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

func TestUsers_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewUsersHandler(mq, slog.Default())

	mq.On("ListUsers", mock.Anything).Return([]dbsqlc.ListUsersRow{
		{AppUserID: 1, Username: "admin", RoleID: "admin"},
		{AppUserID: 2, Username: "teacher", RoleID: "tutor"},
	}, nil)

	w, r := newRequest("GET", "/users", nil, adminUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.ListUsersRow
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 2)
	mq.AssertExpectations(t)
}

func TestUsers_List_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewUsersHandler(mq, slog.Default())

	mq.On("ListUsers", mock.Anything).Return([]dbsqlc.ListUsersRow(nil), errors.New("db down"))

	w, r := newRequest("GET", "/users", nil, adminUser())
	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Create ───────────────────────────────────────────────────────────────────

func TestUsers_Create_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewUsersHandler(mq, slog.Default())

	// The handler bcrypt-hashes the password before calling CreateUser.
	// We can't predict the exact hash, so match on AnythingOfType.
	mq.On("CreateUser", mock.Anything, mock.AnythingOfType("CreateUserParams")).
		Return(dbsqlc.CreateUserRow{AppUserID: 10, Username: "newuser", RoleID: "editor"}, nil)

	w, r := newRequest("POST", "/users", map[string]any{
		"username": "newuser",
		"password": "secret123",
		"role_id":  "editor",
	}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestUsers_Create_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewUsersHandler(mq, slog.Default())

	w, r := newRequest("POST", "/users", map[string]any{
		"username": "newuser", "password": "secret", "role_id": "editor",
	}, editorUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mq.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
}

func TestUsers_Create_MissingFields(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewUsersHandler(mq, slog.Default())

	// Missing password.
	w, r := newRequest("POST", "/users", map[string]any{
		"username": "newuser", "role_id": "editor",
	}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUsers_Create_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewUsersHandler(mq, slog.Default())

	mq.On("CreateUser", mock.Anything, mock.AnythingOfType("CreateUserParams")).
		Return(dbsqlc.CreateUserRow{}, errors.New("duplicate key"))

	w, r := newRequest("POST", "/users", map[string]any{
		"username": "dup", "password": "pass", "role_id": "editor",
	}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Update ───────────────────────────────────────────────────────────────────

func TestUsers_Update_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewUsersHandler(mq, slog.Default())

	mq.On("UpdateUser", mock.Anything, mock.AnythingOfType("UpdateUserParams")).
		Return(dbsqlc.UpdateUserRow{AppUserID: 2, Username: "renamed", RoleID: "tutor"}, nil)

	w, r := newRequest("PATCH", "/users/2", map[string]any{"username": "renamed"}, adminUser())
	r.SetPathValue("id", "2")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestUsers_Update_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewUsersHandler(mq, slog.Default())

	w, r := newRequest("PATCH", "/users/2", map[string]any{"username": "x"}, tutorUser())
	r.SetPathValue("id", "2")
	h.Update(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUsers_Update_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewUsersHandler(mq, slog.Default())

	w, r := newRequest("PATCH", "/users/abc", map[string]any{"username": "x"}, adminUser())
	r.SetPathValue("id", "abc")
	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUsers_Update_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewUsersHandler(mq, slog.Default())

	mq.On("UpdateUser", mock.Anything, mock.AnythingOfType("UpdateUserParams")).
		Return(dbsqlc.UpdateUserRow{}, pgx.ErrNoRows)

	w, r := newRequest("PATCH", "/users/99", map[string]any{"username": "ghost"}, adminUser())
	r.SetPathValue("id", "99")
	h.Update(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── Delete ───────────────────────────────────────────────────────────────────

func TestUsers_Delete_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewUsersHandler(mq, slog.Default())

	mq.On("DeleteUser", mock.Anything, int64(5)).Return(nil)

	w, r := newRequest("DELETE", "/users/5", nil, adminUser())
	r.SetPathValue("id", "5")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}

func TestUsers_Delete_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewUsersHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/users/5", nil, editorUser())
	r.SetPathValue("id", "5")
	h.Delete(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUsers_Delete_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewUsersHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/users/abc", nil, adminUser())
	r.SetPathValue("id", "abc")
	h.Delete(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUsers_Delete_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewUsersHandler(mq, slog.Default())

	mq.On("DeleteUser", mock.Anything, int64(5)).Return(errors.New("fk violation"))

	w, r := newRequest("DELETE", "/users/5", nil, adminUser())
	r.SetPathValue("id", "5")
	h.Delete(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
