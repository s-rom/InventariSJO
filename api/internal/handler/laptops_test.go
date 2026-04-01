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

func TestLaptops_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopsHandler(mq, nil, slog.Default())

	rows := []dbsqlc.ListLaptopsRow{
		{ComputerID: 1, Hostname: "lt-01"},
		{ComputerID: 2, Hostname: "lt-02"},
	}
	mq.On("ListLaptops", mock.Anything).Return(rows, nil)

	w, r := newRequest("GET", "/laptops", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.ListLaptopsRow
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 2)
	mq.AssertExpectations(t)
}

func TestLaptops_List_Empty(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopsHandler(mq, nil, slog.Default())

	mq.On("ListLaptops", mock.Anything).Return([]dbsqlc.ListLaptopsRow{}, nil)

	w, r := newRequest("GET", "/laptops", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLaptops_List_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopsHandler(mq, nil, slog.Default())

	mq.On("ListLaptops", mock.Anything).Return([]dbsqlc.ListLaptopsRow(nil), errors.New("db down"))

	w, r := newRequest("GET", "/laptops", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Get ──────────────────────────────────────────────────────────────────────

func TestLaptops_Get_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopsHandler(mq, nil, slog.Default())

	mq.On("GetLaptop", mock.Anything, int64(7)).
		Return(dbsqlc.GetLaptopRow{ComputerID: 7, Hostname: "lt-07"}, nil)

	w, r := newRequest("GET", "/laptops/7", nil, readonlyUser())
	r.SetPathValue("id", "7")
	h.Get(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestLaptops_Get_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopsHandler(mq, nil, slog.Default())

	mq.On("GetLaptop", mock.Anything, int64(99)).
		Return(dbsqlc.GetLaptopRow{}, pgx.ErrNoRows)

	w, r := newRequest("GET", "/laptops/99", nil, readonlyUser())
	r.SetPathValue("id", "99")
	h.Get(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestLaptops_Get_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopsHandler(mq, nil, slog.Default())

	w, r := newRequest("GET", "/laptops/abc", nil, readonlyUser())
	r.SetPathValue("id", "abc")
	h.Get(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mq.AssertNotCalled(t, "GetLaptop", mock.Anything, mock.Anything)
}

// ─── Create ───────────────────────────────────────────────────────────────────

func TestLaptops_Create_Forbidden_Readonly(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopsHandler(mq, nil, slog.Default())

	w, r := newRequest("POST", "/laptops", map[string]any{"hostname": "lt-10", "laptop_model_id": 1}, readonlyUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestLaptops_Create_Forbidden_Tutor(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopsHandler(mq, nil, slog.Default())

	w, r := newRequest("POST", "/laptops", map[string]any{"hostname": "lt-10", "laptop_model_id": 1}, tutorUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestLaptops_Create_MissingHostname(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopsHandler(mq, nil, slog.Default())

	w, r := newRequest("POST", "/laptops", map[string]any{"laptop_model_id": 1}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLaptops_Create_MissingLaptopModelID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopsHandler(mq, nil, slog.Default())

	w, r := newRequest("POST", "/laptops", map[string]any{"hostname": "lt-10"}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLaptops_Create_BeginFails(t *testing.T) {
	mq := new(testmock.Querier)
	mockDB := new(testmock.MockDB)
	h := handler.NewLaptopsHandler(mq, mockDB, slog.Default())

	mockDB.On("Begin", mock.Anything).Return(nil, errors.New("connection refused"))

	w, r := newRequest("POST", "/laptops", map[string]any{"hostname": "lt-10", "laptop_model_id": 1}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockDB.AssertExpectations(t)
}

// ─── Update ───────────────────────────────────────────────────────────────────

func TestLaptops_Update_Forbidden_Readonly(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopsHandler(mq, nil, slog.Default())

	w, r := newRequest("PATCH", "/laptops/1", map[string]any{"hostname": "new"}, readonlyUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestLaptops_Update_Forbidden_Tutor(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopsHandler(mq, nil, slog.Default())

	w, r := newRequest("PATCH", "/laptops/1", map[string]any{"hostname": "new"}, tutorUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestLaptops_Update_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopsHandler(mq, nil, slog.Default())

	w, r := newRequest("PATCH", "/laptops/abc", map[string]any{"hostname": "new"}, adminUser())
	r.SetPathValue("id", "abc")
	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLaptops_Update_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	mockDB := new(testmock.MockDB)
	h := handler.NewLaptopsHandler(mq, mockDB, slog.Default())

	mockDB.On("Begin", mock.Anything).Return(testmock.NoRowsTx{}, nil)

	w, r := newRequest("PATCH", "/laptops/42", map[string]any{"hostname": "new"}, adminUser())
	r.SetPathValue("id", "42")
	h.Update(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockDB.AssertExpectations(t)
}

func TestLaptops_Update_BeginFails(t *testing.T) {
	mq := new(testmock.Querier)
	mockDB := new(testmock.MockDB)
	h := handler.NewLaptopsHandler(mq, mockDB, slog.Default())

	mockDB.On("Begin", mock.Anything).Return(nil, errors.New("connection refused"))

	w, r := newRequest("PATCH", "/laptops/1", map[string]any{"hostname": "new"}, adminUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockDB.AssertExpectations(t)
}
