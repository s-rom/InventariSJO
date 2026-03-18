package handler_test

import (
"encoding/json"
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

func TestBrands_List_OK(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewBrandsHandler(mq, slog.Default())

expected := []dbsqlc.Brand{{BrandID: 1, Name: "Lenovo"}, {BrandID: 2, Name: "HP"}}
mq.On("ListBrands", mock.Anything).Return(expected, nil)

w, r := newRequest("GET", "/brands", nil, readonlyUser())
h.List(w, r)

assert.Equal(t, http.StatusOK, w.Code)
var got []dbsqlc.Brand
_ = json.Unmarshal(w.Body.Bytes(), &got)
assert.Len(t, got, 2)
mq.AssertExpectations(t)
}

func TestBrands_Create_Admin(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewBrandsHandler(mq, slog.Default())

mq.On("CreateBrand", mock.Anything, "Dell").
Return(dbsqlc.Brand{BrandID: 3, Name: "Dell"}, nil)

w, r := newRequest("POST", "/brands", map[string]any{"name": "Dell"}, adminUser())
h.Create(w, r)

assert.Equal(t, http.StatusCreated, w.Code)
mq.AssertExpectations(t)
}

func TestBrands_Create_Editor(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewBrandsHandler(mq, slog.Default())

mq.On("CreateBrand", mock.Anything, "Asus").
Return(dbsqlc.Brand{BrandID: 4, Name: "Asus"}, nil)

w, r := newRequest("POST", "/brands", map[string]any{"name": "Asus"}, editorUser())
h.Create(w, r)

assert.Equal(t, http.StatusCreated, w.Code)
}

func TestBrands_Create_Readonly_Forbidden(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewBrandsHandler(mq, slog.Default())

w, r := newRequest("POST", "/brands", map[string]any{"name": "Dell"}, readonlyUser())
h.Create(w, r)

assert.Equal(t, http.StatusForbidden, w.Code)
mq.AssertNotCalled(t, "CreateBrand", mock.Anything, mock.Anything)
}

func TestBrands_Create_Tutor_Forbidden(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewBrandsHandler(mq, slog.Default())

w, r := newRequest("POST", "/brands", map[string]any{"name": "Dell"}, tutorUser())
h.Create(w, r)

assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestBrands_Create_EmptyName(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewBrandsHandler(mq, slog.Default())

w, r := newRequest("POST", "/brands", map[string]any{"name": ""}, adminUser())
h.Create(w, r)

assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBrands_Update_NotFound(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewBrandsHandler(mq, slog.Default())

mq.On("UpdateBrand", mock.Anything, mock.AnythingOfType("UpdateBrandParams")).
Return(dbsqlc.Brand{}, pgx.ErrNoRows)

w, r := newRequest("PATCH", "/brands/999", map[string]any{"name": "x"}, adminUser())
r.SetPathValue("id", "999")
h.Update(w, r)

assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBrands_Delete_Admin(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewBrandsHandler(mq, slog.Default())

mq.On("DeleteBrand", mock.Anything, int64(1)).Return(nil)

w, r := newRequest("DELETE", "/brands/1", nil, adminUser())
r.SetPathValue("id", "1")
h.Delete(w, r)

assert.Equal(t, http.StatusNoContent, w.Code)
mq.AssertExpectations(t)
}

func TestBrands_Delete_Readonly_Forbidden(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewBrandsHandler(mq, slog.Default())

w, r := newRequest("DELETE", "/brands/1", nil, readonlyUser())
r.SetPathValue("id", "1")
h.Delete(w, r)

assert.Equal(t, http.StatusForbidden, w.Code)
mq.AssertNotCalled(t, "DeleteBrand", mock.Anything, mock.Anything)
}
