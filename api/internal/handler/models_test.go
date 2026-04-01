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

// ══════════════════════════════════════════════════════════════════════════════
// Laptop Models
// ══════════════════════════════════════════════════════════════════════════════

func TestLaptopModels_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopModelsHandler(mq, slog.Default())

	mq.On("ListLaptopModels", mock.Anything).Return([]dbsqlc.ListLaptopModelsRow{
		{LaptopModelID: 1, ModelName: "ThinkPad T14"},
	}, nil)

	w, r := newRequest("GET", "/laptop-models", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.ListLaptopModelsRow
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 1)
	mq.AssertExpectations(t)
}

func TestLaptopModels_List_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopModelsHandler(mq, slog.Default())

	mq.On("ListLaptopModels", mock.Anything).Return([]dbsqlc.ListLaptopModelsRow(nil), errors.New("db down"))

	w, r := newRequest("GET", "/laptop-models", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestLaptopModels_Get_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopModelsHandler(mq, slog.Default())

	mq.On("GetLaptopModel", mock.Anything, int64(1)).Return(dbsqlc.GetLaptopModelRow{LaptopModelID: 1, ModelName: "ThinkPad T14"}, nil)

	w, r := newRequest("GET", "/laptop-models/1", nil, readonlyUser())
	r.SetPathValue("id", "1")
	h.Get(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestLaptopModels_Get_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopModelsHandler(mq, slog.Default())

	w, r := newRequest("GET", "/laptop-models/abc", nil, readonlyUser())
	r.SetPathValue("id", "abc")
	h.Get(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLaptopModels_Get_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopModelsHandler(mq, slog.Default())

	mq.On("GetLaptopModel", mock.Anything, int64(99)).Return(dbsqlc.GetLaptopModelRow{}, pgx.ErrNoRows)

	w, r := newRequest("GET", "/laptop-models/99", nil, readonlyUser())
	r.SetPathValue("id", "99")
	h.Get(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestLaptopModels_Create_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopModelsHandler(mq, slog.Default())

	mq.On("CreateLaptopModel", mock.Anything, mock.AnythingOfType("CreateLaptopModelParams")).
		Return(dbsqlc.LaptopModel{LaptopModelID: 5, ModelName: "XPS 13"}, nil)

	w, r := newRequest("POST", "/laptop-models", map[string]any{
		"brand_id":         1,
		"model_name":       "XPS 13",
		"base_ram_gb":      16,
		"base_ram_type":    "DDR4",
		"base_storage_gb":  512,
		"base_storage_type": "SSD",
	}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestLaptopModels_Create_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopModelsHandler(mq, slog.Default())

	w, r := newRequest("POST", "/laptop-models", map[string]any{
		"brand_id": 1, "model_name": "X",
	}, tutorUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestLaptopModels_Create_MissingFields(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopModelsHandler(mq, slog.Default())

	// Missing brand_id (zero value == 0 → rejected)
	w, r := newRequest("POST", "/laptop-models", map[string]any{"model_name": "X"}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLaptopModels_Update_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopModelsHandler(mq, slog.Default())

	mq.On("UpdateLaptopModel", mock.Anything, mock.AnythingOfType("UpdateLaptopModelParams")).
		Return(dbsqlc.LaptopModel{LaptopModelID: 1, ModelName: "XPS 15"}, nil)

	w, r := newRequest("PATCH", "/laptop-models/1", map[string]any{"model_name": "XPS 15"}, adminUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestLaptopModels_Update_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopModelsHandler(mq, slog.Default())

	w, r := newRequest("PATCH", "/laptop-models/abc", map[string]any{"model_name": "X"}, adminUser())
	r.SetPathValue("id", "abc")
	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLaptopModels_Update_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopModelsHandler(mq, slog.Default())

	mq.On("UpdateLaptopModel", mock.Anything, mock.AnythingOfType("UpdateLaptopModelParams")).
		Return(dbsqlc.LaptopModel{}, pgx.ErrNoRows)

	w, r := newRequest("PATCH", "/laptop-models/99", map[string]any{"model_name": "Ghost"}, adminUser())
	r.SetPathValue("id", "99")
	h.Update(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestLaptopModels_Delete_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopModelsHandler(mq, slog.Default())

	mq.On("DeleteLaptopModel", mock.Anything, int64(3)).Return(nil)

	w, r := newRequest("DELETE", "/laptop-models/3", nil, adminUser())
	r.SetPathValue("id", "3")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}

func TestLaptopModels_Delete_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopModelsHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/laptop-models/3", nil, tutorUser())
	r.SetPathValue("id", "3")
	h.Delete(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestLaptopModels_Delete_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewLaptopModelsHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/laptop-models/abc", nil, adminUser())
	r.SetPathValue("id", "abc")
	h.Delete(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ══════════════════════════════════════════════════════════════════════════════
// Desktop Models
// ══════════════════════════════════════════════════════════════════════════════

func TestDesktopModels_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopModelsHandler(mq, slog.Default())

	mq.On("ListDesktopModels", mock.Anything).Return([]dbsqlc.ListDesktopModelsRow{
		{DesktopModelID: 1, ModelName: "OptiPlex 7090"},
	}, nil)

	w, r := newRequest("GET", "/desktop-models", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.ListDesktopModelsRow
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 1)
	mq.AssertExpectations(t)
}

func TestDesktopModels_List_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopModelsHandler(mq, slog.Default())

	mq.On("ListDesktopModels", mock.Anything).Return([]dbsqlc.ListDesktopModelsRow(nil), errors.New("db down"))

	w, r := newRequest("GET", "/desktop-models", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDesktopModels_Get_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopModelsHandler(mq, slog.Default())

	mq.On("GetDesktopModel", mock.Anything, int64(1)).Return(dbsqlc.GetDesktopModelRow{DesktopModelID: 1, ModelName: "OptiPlex"}, nil)

	w, r := newRequest("GET", "/desktop-models/1", nil, readonlyUser())
	r.SetPathValue("id", "1")
	h.Get(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestDesktopModels_Get_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopModelsHandler(mq, slog.Default())

	mq.On("GetDesktopModel", mock.Anything, int64(99)).Return(dbsqlc.GetDesktopModelRow{}, pgx.ErrNoRows)

	w, r := newRequest("GET", "/desktop-models/99", nil, readonlyUser())
	r.SetPathValue("id", "99")
	h.Get(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDesktopModels_Create_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopModelsHandler(mq, slog.Default())

	mq.On("CreateDesktopModel", mock.Anything, mock.AnythingOfType("CreateDesktopModelParams")).
		Return(dbsqlc.DesktopModel{DesktopModelID: 2, ModelName: "OptiPlex 7090"}, nil)

	w, r := newRequest("POST", "/desktop-models", map[string]any{
		"brand_id":          1,
		"model_name":        "OptiPlex 7090",
		"base_ram_gb":       16,
		"base_ram_type":     "DDR4",
		"base_storage_gb":   512,
		"base_storage_type": "SSD",
	}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestDesktopModels_Create_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopModelsHandler(mq, slog.Default())

	w, r := newRequest("POST", "/desktop-models", map[string]any{
		"brand_id": 1, "model_name": "X",
	}, readonlyUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDesktopModels_Update_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopModelsHandler(mq, slog.Default())

	mq.On("UpdateDesktopModel", mock.Anything, mock.AnythingOfType("UpdateDesktopModelParams")).
		Return(dbsqlc.DesktopModel{DesktopModelID: 1, ModelName: "OptiPlex 7091"}, nil)

	w, r := newRequest("PATCH", "/desktop-models/1", map[string]any{"model_name": "OptiPlex 7091"}, editorUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestDesktopModels_Update_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopModelsHandler(mq, slog.Default())

	mq.On("UpdateDesktopModel", mock.Anything, mock.AnythingOfType("UpdateDesktopModelParams")).
		Return(dbsqlc.DesktopModel{}, pgx.ErrNoRows)

	w, r := newRequest("PATCH", "/desktop-models/99", map[string]any{"model_name": "Ghost"}, adminUser())
	r.SetPathValue("id", "99")
	h.Update(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDesktopModels_Delete_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopModelsHandler(mq, slog.Default())

	mq.On("DeleteDesktopModel", mock.Anything, int64(4)).Return(nil)

	w, r := newRequest("DELETE", "/desktop-models/4", nil, adminUser())
	r.SetPathValue("id", "4")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}

func TestDesktopModels_Delete_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopModelsHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/desktop-models/4", nil, readonlyUser())
	r.SetPathValue("id", "4")
	h.Delete(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDesktopModels_Delete_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewDesktopModelsHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/desktop-models/xyz", nil, adminUser())
	r.SetPathValue("id", "xyz")
	h.Delete(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
