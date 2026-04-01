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

// ─── ListByCenter ─────────────────────────────────────────────────────────────

func TestRooms_ListByCenter_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRoomsHandler(mq, slog.Default())

	mq.On("ListRoomsByCenter", mock.Anything, int64(10)).Return([]dbsqlc.Room{
		{RoomID: 1, CenterID: 10, Name: "Aula 1"},
		{RoomID: 2, CenterID: 10, Name: "Aula 2"},
	}, nil)

	w, r := newRequest("GET", "/centers/10/rooms", nil, readonlyUser())
	r.SetPathValue("centerId", "10")
	h.ListByCenter(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.Room
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 2)
	mq.AssertExpectations(t)
}

func TestRooms_ListByCenter_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRoomsHandler(mq, slog.Default())

	w, r := newRequest("GET", "/centers/abc/rooms", nil, readonlyUser())
	r.SetPathValue("centerId", "abc")
	h.ListByCenter(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRooms_ListByCenter_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRoomsHandler(mq, slog.Default())

	mq.On("ListRoomsByCenter", mock.Anything, int64(1)).Return([]dbsqlc.Room(nil), errors.New("db down"))

	w, r := newRequest("GET", "/centers/1/rooms", nil, readonlyUser())
	r.SetPathValue("centerId", "1")
	h.ListByCenter(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Create ───────────────────────────────────────────────────────────────────

func TestRooms_Create_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRoomsHandler(mq, slog.Default())

	mq.On("CreateRoom", mock.Anything, dbsqlc.CreateRoomParams{CenterID: 5, Name: "Sala TIC"}).
		Return(dbsqlc.Room{RoomID: 1, CenterID: 5, Name: "Sala TIC"}, nil)

	w, r := newRequest("POST", "/centers/5/rooms", map[string]any{"name": "Sala TIC"}, adminUser())
	r.SetPathValue("centerId", "5")
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestRooms_Create_BadCenterID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRoomsHandler(mq, slog.Default())

	w, r := newRequest("POST", "/centers/abc/rooms", map[string]any{"name": "X"}, adminUser())
	r.SetPathValue("centerId", "abc")
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRooms_Create_EmptyName(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRoomsHandler(mq, slog.Default())

	w, r := newRequest("POST", "/centers/1/rooms", map[string]any{"name": ""}, adminUser())
	r.SetPathValue("centerId", "1")
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── Update ───────────────────────────────────────────────────────────────────

func TestRooms_Update_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRoomsHandler(mq, slog.Default())

	mq.On("UpdateRoom", mock.Anything, dbsqlc.UpdateRoomParams{RoomID: 3, Name: "Laboratori"}).
		Return(dbsqlc.Room{RoomID: 3, CenterID: 1, Name: "Laboratori"}, nil)

	w, r := newRequest("PATCH", "/rooms/3", map[string]any{"name": "Laboratori"}, adminUser())
	r.SetPathValue("id", "3")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestRooms_Update_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRoomsHandler(mq, slog.Default())

	w, r := newRequest("PATCH", "/rooms/abc", map[string]any{"name": "X"}, adminUser())
	r.SetPathValue("id", "abc")
	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRooms_Update_EmptyName(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRoomsHandler(mq, slog.Default())

	w, r := newRequest("PATCH", "/rooms/1", map[string]any{"name": ""}, adminUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRooms_Update_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRoomsHandler(mq, slog.Default())

	mq.On("UpdateRoom", mock.Anything, mock.AnythingOfType("UpdateRoomParams")).
		Return(dbsqlc.Room{}, errors.New("fk error"))

	w, r := newRequest("PATCH", "/rooms/1", map[string]any{"name": "X"}, adminUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Delete ───────────────────────────────────────────────────────────────────

func TestRooms_Delete_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRoomsHandler(mq, slog.Default())

	mq.On("DeleteRoom", mock.Anything, int64(7)).Return(nil)

	w, r := newRequest("DELETE", "/rooms/7", nil, adminUser())
	r.SetPathValue("id", "7")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}

func TestRooms_Delete_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRoomsHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/rooms/xyz", nil, adminUser())
	r.SetPathValue("id", "xyz")
	h.Delete(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRooms_Delete_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewRoomsHandler(mq, slog.Default())

	mq.On("DeleteRoom", mock.Anything, int64(1)).Return(errors.New("fk violation"))

	w, r := newRequest("DELETE", "/rooms/1", nil, adminUser())
	r.SetPathValue("id", "1")
	h.Delete(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
