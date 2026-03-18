package handler_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"testing"

	dbsqlc "inventari/api/internal/db/sqlc"
	"inventari/api/internal/handler"
	testmock "inventari/api/internal/testutil/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCPUs_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCPUsHandler(mq, slog.Default())

	expected := []dbsqlc.Cpu{{CpuID: 1, ModelName: "i7-12700"}}
	mq.On("ListCPUs", mock.Anything).Return(expected, nil)

	w, r := newRequest("GET", "/cpus", nil, adminUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.Cpu
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 1)
	assert.Equal(t, "i7-12700", got[0].ModelName)
	mq.AssertExpectations(t)
}

func TestCPUs_Create_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCPUsHandler(mq, slog.Default())

	mq.On("CreateCPU", mock.Anything, mock.AnythingOfType("CreateCPUParams")).
		Return(dbsqlc.Cpu{CpuID: 2, ModelName: "i5-13400"}, nil)

	body := map[string]any{"model_name": "i5-13400"}
	w, r := newRequest("POST", "/cpus", body, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestCPUs_Create_EmptyName(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCPUsHandler(mq, slog.Default())

	body := map[string]any{"model_name": ""}
	w, r := newRequest("POST", "/cpus", body, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mq.AssertNotCalled(t, "CreateCPU", mock.Anything, mock.Anything)
}

func TestCPUs_Update_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCPUsHandler(mq, slog.Default())

	mq.On("UpdateCPU", mock.Anything, mock.AnythingOfType("UpdateCPUParams")).
		Return(dbsqlc.Cpu{CpuID: 1, ModelName: "updated"}, nil)

	body := map[string]any{"model_name": "updated"}
	w, r := newRequest("PATCH", "/cpus/1", body, adminUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestCPUs_Update_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCPUsHandler(mq, slog.Default())

	w, r := newRequest("PATCH", "/cpus/abc", nil, adminUser())
	r.SetPathValue("id", "abc")
	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCPUs_Delete_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCPUsHandler(mq, slog.Default())

	mq.On("DeleteCPU", mock.Anything, int64(5)).Return(nil)

	w, r := newRequest("DELETE", "/cpus/5", nil, adminUser())
	r.SetPathValue("id", "5")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}

func TestCPUs_Delete_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewCPUsHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/cpus/abc", nil, adminUser())
	r.SetPathValue("id", "abc")
	h.Delete(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
