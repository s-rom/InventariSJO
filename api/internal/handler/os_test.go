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

func TestOS_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewOSHandler(mq, slog.Default())

	mq.On("ListOS", mock.Anything).Return([]dbsqlc.O{{OsID: 1, Name: "Ubuntu"}}, nil)

	w, r := newRequest("GET", "/os", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.O
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 1)
	assert.Equal(t, "Ubuntu", got[0].Name)
	mq.AssertExpectations(t)
}

func TestOS_List_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewOSHandler(mq, slog.Default())

	mq.On("ListOS", mock.Anything).Return([]dbsqlc.O(nil), errors.New("db down"))

	w, r := newRequest("GET", "/os", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Create ───────────────────────────────────────────────────────────────────

func TestOS_Create_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewOSHandler(mq, slog.Default())

	mq.On("CreateOS", mock.Anything, "Debian").Return(dbsqlc.O{OsID: 2, Name: "Debian"}, nil)

	w, r := newRequest("POST", "/os", map[string]any{"name": "Debian"}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestOS_Create_EmptyName(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewOSHandler(mq, slog.Default())

	w, r := newRequest("POST", "/os", map[string]any{"name": ""}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mq.AssertNotCalled(t, "CreateOS", mock.Anything, mock.Anything)
}

// ─── Update ───────────────────────────────────────────────────────────────────

func TestOS_Update_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewOSHandler(mq, slog.Default())

	mq.On("UpdateOS", mock.Anything, mock.AnythingOfType("UpdateOSParams")).
		Return(dbsqlc.O{OsID: 1, Name: "Ubuntu 24.04"}, nil)

	w, r := newRequest("PATCH", "/os/1", map[string]any{"name": "Ubuntu 24.04"}, adminUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestOS_Update_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewOSHandler(mq, slog.Default())

	w, r := newRequest("PATCH", "/os/abc", map[string]any{"name": "X"}, adminUser())
	r.SetPathValue("id", "abc")
	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOS_Update_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewOSHandler(mq, slog.Default())

	mq.On("UpdateOS", mock.Anything, mock.AnythingOfType("UpdateOSParams")).
		Return(dbsqlc.O{}, pgx.ErrNoRows)

	w, r := newRequest("PATCH", "/os/99", map[string]any{"name": "X"}, adminUser())
	r.SetPathValue("id", "99")
	h.Update(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── Delete ───────────────────────────────────────────────────────────────────

func TestOS_Delete_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewOSHandler(mq, slog.Default())

	mq.On("DeleteOS", mock.Anything, int64(1)).Return(nil)

	w, r := newRequest("DELETE", "/os/1", nil, adminUser())
	r.SetPathValue("id", "1")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}

func TestOS_Delete_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewOSHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/os/abc", nil, adminUser())
	r.SetPathValue("id", "abc")
	h.Delete(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
