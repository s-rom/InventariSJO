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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─── List ─────────────────────────────────────────────────────────────────────

func TestCenters_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCentersHandler(mq, slog.Default())

	mq.On("ListCenters", mock.Anything).Return([]dbsqlc.Center{
		{CenterID: 1, Name: "Institut"},
		{CenterID: 2, Name: "Escola"},
	}, nil)

	w, r := newRequest("GET", "/centers", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.Center
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 2)
	mq.AssertExpectations(t)
}

func TestCenters_List_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCentersHandler(mq, slog.Default())

	mq.On("ListCenters", mock.Anything).Return([]dbsqlc.Center(nil), errors.New("db down"))

	w, r := newRequest("GET", "/centers", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Create ───────────────────────────────────────────────────────────────────

func TestCenters_Create_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCentersHandler(mq, slog.Default())

	mq.On("CreateCenter", mock.Anything, "INS Lleida").
		Return(dbsqlc.Center{CenterID: 3, Name: "INS Lleida"}, nil)

	w, r := newRequest("POST", "/centers", map[string]any{"name": "INS Lleida"}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestCenters_Create_EmptyName(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCentersHandler(mq, slog.Default())

	w, r := newRequest("POST", "/centers", map[string]any{"name": ""}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mq.AssertNotCalled(t, "CreateCenter", mock.Anything, mock.Anything)
}

// ─── Update ───────────────────────────────────────────────────────────────────

func TestCenters_Update_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCentersHandler(mq, slog.Default())

	mq.On("UpdateCenter", mock.Anything, mock.AnythingOfType("UpdateCenterParams")).
		Return(dbsqlc.Center{CenterID: 1, Name: "INS Lleida"}, nil)

	w, r := newRequest("PATCH", "/centers/1", map[string]any{"name": "INS Lleida"}, adminUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestCenters_Update_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCentersHandler(mq, slog.Default())

	w, r := newRequest("PATCH", "/centers/abc", map[string]any{"name": "X"}, adminUser())
	r.SetPathValue("id", "abc")
	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCenters_Update_EmptyName(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCentersHandler(mq, slog.Default())

	w, r := newRequest("PATCH", "/centers/1", map[string]any{"name": ""}, adminUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── Delete ───────────────────────────────────────────────────────────────────

func TestCenters_Delete_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCentersHandler(mq, slog.Default())

	mq.On("DeleteCenter", mock.Anything, int64(2)).Return(nil)

	w, r := newRequest("DELETE", "/centers/2", nil, adminUser())
	r.SetPathValue("id", "2")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}

func TestCenters_Delete_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCentersHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/centers/abc", nil, adminUser())
	r.SetPathValue("id", "abc")
	h.Delete(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCenters_Delete_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCentersHandler(mq, slog.Default())

	mq.On("DeleteCenter", mock.Anything, int64(1)).Return(errors.New("fk violation"))

	w, r := newRequest("DELETE", "/centers/1", nil, adminUser())
	r.SetPathValue("id", "1")
	h.Delete(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
