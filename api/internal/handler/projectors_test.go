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

// ═══════════════════════════════════════════════════════════════
// ProjectorModels
// ═══════════════════════════════════════════════════════════════

func TestProjectorModels_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorModelsHandler(mq, slog.Default())

	rows := []dbsqlc.ListProjectorModelsRow{{ProjectorModelID: 1, ModelName: "EX3260"}}
	mq.On("ListProjectorModels", mock.Anything).Return(rows, nil)

	w, r := newRequest("GET", "/projector-models", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.ListProjectorModelsRow
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 1)
	mq.AssertExpectations(t)
}

func TestProjectorModels_List_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorModelsHandler(mq, slog.Default())

	mq.On("ListProjectorModels", mock.Anything).Return([]dbsqlc.ListProjectorModelsRow(nil), errors.New("db down"))

	w, r := newRequest("GET", "/projector-models", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProjectorModels_Get_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorModelsHandler(mq, slog.Default())

	mq.On("GetProjectorModel", mock.Anything, int64(2)).
		Return(dbsqlc.GetProjectorModelRow{ProjectorModelID: 2, ModelName: "EX3260"}, nil)

	w, r := newRequest("GET", "/projector-models/2", nil, readonlyUser())
	r.SetPathValue("id", "2")
	h.Get(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestProjectorModels_Get_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorModelsHandler(mq, slog.Default())

	mq.On("GetProjectorModel", mock.Anything, int64(99)).
		Return(dbsqlc.GetProjectorModelRow{}, pgx.ErrNoRows)

	w, r := newRequest("GET", "/projector-models/99", nil, readonlyUser())
	r.SetPathValue("id", "99")
	h.Get(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProjectorModels_Get_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorModelsHandler(mq, slog.Default())

	w, r := newRequest("GET", "/projector-models/abc", nil, readonlyUser())
	r.SetPathValue("id", "abc")
	h.Get(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProjectorModels_Create_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorModelsHandler(mq, slog.Default())

	created := dbsqlc.ProjectorModel{ProjectorModelID: 1, BrandID: 3, ModelName: "EX3260"}
	mq.On("CreateProjectorModel", mock.Anything, dbsqlc.CreateProjectorModelParams{
		BrandID: 3, ModelName: "EX3260",
	}).Return(created, nil)

	body := map[string]any{"brand_id": 3, "model_name": "EX3260"}
	w, r := newRequest("POST", "/projector-models", body, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestProjectorModels_Create_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorModelsHandler(mq, slog.Default())

	w, r := newRequest("POST", "/projector-models", map[string]any{"brand_id": 1, "model_name": "x"}, readonlyUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestProjectorModels_Create_MissingFields(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorModelsHandler(mq, slog.Default())

	// missing model_name
	w, r := newRequest("POST", "/projector-models", map[string]any{"brand_id": 1}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProjectorModels_Update_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorModelsHandler(mq, slog.Default())

	updated := dbsqlc.ProjectorModel{ProjectorModelID: 1, ModelName: "NewName"}
	mq.On("UpdateProjectorModel", mock.Anything, mock.Anything).Return(updated, nil)

	body := map[string]any{"model_name": "NewName"}
	w, r := newRequest("PATCH", "/projector-models/1", body, adminUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestProjectorModels_Update_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorModelsHandler(mq, slog.Default())

	mq.On("UpdateProjectorModel", mock.Anything, mock.Anything).Return(dbsqlc.ProjectorModel{}, pgx.ErrNoRows)

	w, r := newRequest("PATCH", "/projector-models/99", map[string]any{"model_name": "x"}, adminUser())
	r.SetPathValue("id", "99")
	h.Update(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProjectorModels_Update_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorModelsHandler(mq, slog.Default())

	w, r := newRequest("PATCH", "/projector-models/1", map[string]any{"model_name": "x"}, readonlyUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestProjectorModels_Delete_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorModelsHandler(mq, slog.Default())

	mq.On("DeleteProjectorModel", mock.Anything, int64(1)).Return(nil)

	w, r := newRequest("DELETE", "/projector-models/1", nil, adminUser())
	r.SetPathValue("id", "1")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}

func TestProjectorModels_Delete_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorModelsHandler(mq, slog.Default())

	mq.On("DeleteProjectorModel", mock.Anything, int64(99)).Return(pgx.ErrNoRows)

	w, r := newRequest("DELETE", "/projector-models/99", nil, adminUser())
	r.SetPathValue("id", "99")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProjectorModels_Delete_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorModelsHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/projector-models/1", nil, editorUser())
	r.SetPathValue("id", "1")
	h.Delete(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ═══════════════════════════════════════════════════════════════
// Projectors
// ═══════════════════════════════════════════════════════════════

func TestProjectors_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorsHandler(mq, nil, slog.Default())

	rows := []dbsqlc.ListProjectorsRow{
		{ProjectorID: 1, ProjectorModelID: 2},
		{ProjectorID: 2, ProjectorModelID: 2},
	}
	mq.On("ListProjectors", mock.Anything).Return(rows, nil)

	w, r := newRequest("GET", "/projectors", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.ListProjectorsRow
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 2)
	mq.AssertExpectations(t)
}

func TestProjectors_List_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorsHandler(mq, nil, slog.Default())

	mq.On("ListProjectors", mock.Anything).Return([]dbsqlc.ListProjectorsRow(nil), errors.New("db down"))

	w, r := newRequest("GET", "/projectors", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProjectors_Get_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorsHandler(mq, nil, slog.Default())

	mq.On("GetProjector", mock.Anything, int64(4)).
		Return(dbsqlc.GetProjectorRow{ProjectorID: 4}, nil)

	w, r := newRequest("GET", "/projectors/4", nil, readonlyUser())
	r.SetPathValue("id", "4")
	h.Get(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestProjectors_Get_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorsHandler(mq, nil, slog.Default())

	mq.On("GetProjector", mock.Anything, int64(99)).
		Return(dbsqlc.GetProjectorRow{}, pgx.ErrNoRows)

	w, r := newRequest("GET", "/projectors/99", nil, readonlyUser())
	r.SetPathValue("id", "99")
	h.Get(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProjectors_Get_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorsHandler(mq, nil, slog.Default())

	w, r := newRequest("GET", "/projectors/abc", nil, readonlyUser())
	r.SetPathValue("id", "abc")
	h.Get(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProjectors_Create_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorsHandler(mq, nil, slog.Default())

	body := map[string]any{"projector_model_id": 1}
	w, r := newRequest("POST", "/projectors", body, readonlyUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestProjectors_Create_MissingModel(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorsHandler(mq, nil, slog.Default())

	w, r := newRequest("POST", "/projectors", map[string]any{}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProjectors_Create_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorsHandler(mq, nil, slog.Default())

	created := dbsqlc.Projector{ProjectorID: 10, ProjectorModelID: 3}
	mq.On("CreateProjector", mock.Anything, mock.MatchedBy(func(p dbsqlc.CreateProjectorParams) bool {
		return p.ProjectorModelID == 3
	})).Return(created, nil)
	mq.On("InsertAuditLog", mock.Anything, mock.Anything).Return(nil)
	mq.On("GetProjector", mock.Anything, int64(10)).Return(dbsqlc.GetProjectorRow{ProjectorID: 10}, nil)

	body := map[string]any{"projector_model_id": 3, "status": "actiu"}
	w, r := newRequest("POST", "/projectors", body, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestProjectors_Update_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorsHandler(mq, nil, slog.Default())

	w, r := newRequest("PATCH", "/projectors/1", map[string]any{}, readonlyUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestProjectors_Update_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorsHandler(mq, nil, slog.Default())

	mq.On("GetProjector", mock.Anything, int64(99)).Return(dbsqlc.GetProjectorRow{}, pgx.ErrNoRows)

	w, r := newRequest("PATCH", "/projectors/99", map[string]any{}, adminUser())
	r.SetPathValue("id", "99")
	h.Update(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mq.AssertExpectations(t)
}

func TestProjectors_Update_ClearsObservations(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorsHandler(mq, nil, slog.Default())

	old := dbsqlc.GetProjectorRow{ProjectorID: 1}
	updated := dbsqlc.Projector{ProjectorID: 1}
	final := dbsqlc.GetProjectorRow{ProjectorID: 1}

	mq.On("GetProjector", mock.Anything, int64(1)).Return(old, nil).Once()
	mq.On("UpdateProjector", mock.Anything, mock.MatchedBy(func(p dbsqlc.UpdateProjectorParams) bool {
		// observations must be NULL (not Valid) when client sends null
		return p.ProjectorID == 1 && !p.Observations.Valid
	})).Return(updated, nil)
	mq.On("InsertAuditLog", mock.Anything, mock.Anything).Return(nil)
	mq.On("GetProjector", mock.Anything, int64(1)).Return(final, nil).Once()

	body := map[string]any{"observations": nil}
	w, r := newRequest("PATCH", "/projectors/1", body, adminUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestProjectors_Delete_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorsHandler(mq, nil, slog.Default())

	w, r := newRequest("DELETE", "/projectors/1", nil, editorUser())
	r.SetPathValue("id", "1")
	h.Delete(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestProjectors_Delete_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorsHandler(mq, nil, slog.Default())

	mq.On("GetProjector", mock.Anything, int64(99)).Return(dbsqlc.GetProjectorRow{}, pgx.ErrNoRows)

	w, r := newRequest("DELETE", "/projectors/99", nil, adminUser())
	r.SetPathValue("id", "99")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProjectors_Delete_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewProjectorsHandler(mq, nil, slog.Default())

	old := dbsqlc.GetProjectorRow{ProjectorID: 1}
	mq.On("GetProjector", mock.Anything, int64(1)).Return(old, nil)
	mq.On("DeleteProjector", mock.Anything, int64(1)).Return(nil)
	mq.On("InsertAuditLog", mock.Anything, mock.Anything).Return(nil)

	w, r := newRequest("DELETE", "/projectors/1", nil, adminUser())
	r.SetPathValue("id", "1")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}
