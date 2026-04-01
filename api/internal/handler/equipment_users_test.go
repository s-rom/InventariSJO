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

func TestEquipmentUsers_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewEquipmentUsersHandler(mq, slog.Default())

	mq.On("ListEquipmentUsers", mock.Anything).Return([]dbsqlc.EquipmentUser{
		{EquipmentUserID: 1, Name: "Laboratori"},
		{EquipmentUserID: 2, Name: "Biblioteca"},
	}, nil)

	w, r := newRequest("GET", "/equipment-users", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.EquipmentUser
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 2)
	mq.AssertExpectations(t)
}

func TestEquipmentUsers_List_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewEquipmentUsersHandler(mq, slog.Default())

	mq.On("ListEquipmentUsers", mock.Anything).Return([]dbsqlc.EquipmentUser(nil), errors.New("db down"))

	w, r := newRequest("GET", "/equipment-users", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Create ───────────────────────────────────────────────────────────────────

func TestEquipmentUsers_Create_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewEquipmentUsersHandler(mq, slog.Default())

	mq.On("CreateEquipmentUser", mock.Anything, "Secretaria").
		Return(dbsqlc.EquipmentUser{EquipmentUserID: 3, Name: "Secretaria"}, nil)

	w, r := newRequest("POST", "/equipment-users", map[string]any{"name": "Secretaria"}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestEquipmentUsers_Create_EmptyName(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewEquipmentUsersHandler(mq, slog.Default())

	w, r := newRequest("POST", "/equipment-users", map[string]any{"name": ""}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mq.AssertNotCalled(t, "CreateEquipmentUser", mock.Anything, mock.Anything)
}

func TestEquipmentUsers_Create_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewEquipmentUsersHandler(mq, slog.Default())

	mq.On("CreateEquipmentUser", mock.Anything, "Dup").Return(dbsqlc.EquipmentUser{}, errors.New("duplicate"))

	w, r := newRequest("POST", "/equipment-users", map[string]any{"name": "Dup"}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Update ───────────────────────────────────────────────────────────────────

func TestEquipmentUsers_Update_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewEquipmentUsersHandler(mq, slog.Default())

	mq.On("UpdateEquipmentUser", mock.Anything, dbsqlc.UpdateEquipmentUserParams{EquipmentUserID: 2, Name: "Sala d'actes"}).
		Return(dbsqlc.EquipmentUser{EquipmentUserID: 2, Name: "Sala d'actes"}, nil)

	w, r := newRequest("PATCH", "/equipment-users/2", map[string]any{"name": "Sala d'actes"}, adminUser())
	r.SetPathValue("id", "2")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestEquipmentUsers_Update_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewEquipmentUsersHandler(mq, slog.Default())

	w, r := newRequest("PATCH", "/equipment-users/abc", map[string]any{"name": "X"}, adminUser())
	r.SetPathValue("id", "abc")
	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEquipmentUsers_Update_EmptyName(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewEquipmentUsersHandler(mq, slog.Default())

	w, r := newRequest("PATCH", "/equipment-users/1", map[string]any{"name": ""}, adminUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEquipmentUsers_Update_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewEquipmentUsersHandler(mq, slog.Default())

	mq.On("UpdateEquipmentUser", mock.Anything, mock.AnythingOfType("UpdateEquipmentUserParams")).
		Return(dbsqlc.EquipmentUser{}, errors.New("db error"))

	w, r := newRequest("PATCH", "/equipment-users/1", map[string]any{"name": "X"}, adminUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Delete ───────────────────────────────────────────────────────────────────

func TestEquipmentUsers_Delete_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewEquipmentUsersHandler(mq, slog.Default())

	mq.On("DeleteEquipmentUser", mock.Anything, int64(2)).Return(nil)

	w, r := newRequest("DELETE", "/equipment-users/2", nil, adminUser())
	r.SetPathValue("id", "2")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}

func TestEquipmentUsers_Delete_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewEquipmentUsersHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/equipment-users/abc", nil, adminUser())
	r.SetPathValue("id", "abc")
	h.Delete(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEquipmentUsers_Delete_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewEquipmentUsersHandler(mq, slog.Default())

	mq.On("DeleteEquipmentUser", mock.Anything, int64(1)).Return(errors.New("fk violation"))

	w, r := newRequest("DELETE", "/equipment-users/1", nil, adminUser())
	r.SetPathValue("id", "1")
	h.Delete(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
