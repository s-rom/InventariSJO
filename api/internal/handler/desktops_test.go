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

func TestDesktops_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopsHandler(mq, nil, slog.Default())

	rows := []dbsqlc.ListDesktopsRow{
		{ComputerID: 1, Hostname: "pc-01"},
		{ComputerID: 2, Hostname: "pc-02"},
	}
	mq.On("ListDesktops", mock.Anything).Return(rows, nil)

	w, r := newRequest("GET", "/desktops", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.ListDesktopsRow
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 2)
	mq.AssertExpectations(t)
}

func TestDesktops_List_Empty(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopsHandler(mq, nil, slog.Default())

	mq.On("ListDesktops", mock.Anything).Return([]dbsqlc.ListDesktopsRow{}, nil)

	w, r := newRequest("GET", "/desktops", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.ListDesktopsRow
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Empty(t, got)
}

func TestDesktops_List_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopsHandler(mq, nil, slog.Default())

	mq.On("ListDesktops", mock.Anything).Return([]dbsqlc.ListDesktopsRow(nil), errors.New("db down"))

	w, r := newRequest("GET", "/desktops", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Get ──────────────────────────────────────────────────────────────────────

func TestDesktops_Get_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopsHandler(mq, nil, slog.Default())

	mq.On("GetDesktop", mock.Anything, int64(5)).
		Return(dbsqlc.GetDesktopRow{ComputerID: 5, Hostname: "pc-05"}, nil)

	w, r := newRequest("GET", "/desktops/5", nil, readonlyUser())
	r.SetPathValue("id", "5")
	h.Get(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestDesktops_Get_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopsHandler(mq, nil, slog.Default())

	mq.On("GetDesktop", mock.Anything, int64(99)).
		Return(dbsqlc.GetDesktopRow{}, pgx.ErrNoRows)

	w, r := newRequest("GET", "/desktops/99", nil, readonlyUser())
	r.SetPathValue("id", "99")
	h.Get(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDesktops_Get_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopsHandler(mq, nil, slog.Default())

	w, r := newRequest("GET", "/desktops/abc", nil, readonlyUser())
	r.SetPathValue("id", "abc")
	h.Get(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mq.AssertNotCalled(t, "GetDesktop", mock.Anything, mock.Anything)
}

// ─── Create ───────────────────────────────────────────────────────────────────

func TestDesktops_Create_Forbidden_Readonly(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopsHandler(mq, nil, slog.Default())

	w, r := newRequest("POST", "/desktops", map[string]any{"hostname": "pc-10"}, readonlyUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDesktops_Create_Forbidden_Tutor(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopsHandler(mq, nil, slog.Default())

	w, r := newRequest("POST", "/desktops", map[string]any{"hostname": "pc-10"}, tutorUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDesktops_Create_MissingHostname(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopsHandler(mq, nil, slog.Default())

	w, r := newRequest("POST", "/desktops", map[string]any{"observations": "missing hostname"}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDesktops_Create_BeginFails(t *testing.T) {
	mq := new(testmock.Querier)
	mockDB := new(testmock.MockDB)
	h := handler.NewDesktopsHandler(mq, mockDB, slog.Default())

	mockDB.On("Begin", mock.Anything).Return(nil, errors.New("connection refused"))

	w, r := newRequest("POST", "/desktops", map[string]any{"hostname": "pc-10"}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockDB.AssertExpectations(t)
}

// ─── Update ───────────────────────────────────────────────────────────────────

func TestDesktops_Update_Forbidden_Readonly(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopsHandler(mq, nil, slog.Default())

	w, r := newRequest("PATCH", "/desktops/1", map[string]any{"hostname": "new"}, readonlyUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDesktops_Update_Forbidden_Tutor(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopsHandler(mq, nil, slog.Default())

	w, r := newRequest("PATCH", "/desktops/1", map[string]any{"hostname": "new"}, tutorUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDesktops_Update_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopsHandler(mq, nil, slog.Default())

	w, r := newRequest("PATCH", "/desktops/abc", map[string]any{"hostname": "new"}, adminUser())
	r.SetPathValue("id", "abc")
	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDesktops_Update_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	mockDB := new(testmock.MockDB)
	h := handler.NewDesktopsHandler(mq, mockDB, slog.Default())

	mockDB.On("Begin", mock.Anything).Return(testmock.NoRowsTx{}, nil)

	w, r := newRequest("PATCH", "/desktops/42", map[string]any{"hostname": "new"}, adminUser())
	r.SetPathValue("id", "42")
	h.Update(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockDB.AssertExpectations(t)
}

func TestDesktops_Update_BeginFails(t *testing.T) {
	mq := new(testmock.Querier)
	mockDB := new(testmock.MockDB)
	h := handler.NewDesktopsHandler(mq, mockDB, slog.Default())

	mockDB.On("Begin", mock.Anything).Return(nil, errors.New("connection refused"))

	w, r := newRequest("PATCH", "/desktops/1", map[string]any{"hostname": "new"}, adminUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockDB.AssertExpectations(t)
}
