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

func TestComputers_List_OK(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewComputersHandler(mq, nil, slog.Default())

rows := []dbsqlc.ListComputersRow{
{ComputerID: 1, Hostname: "pc-01"},
{ComputerID: 2, Hostname: "lt-01"},
}
mq.On("ListComputers", mock.Anything).Return(rows, nil)

w, r := newRequest("GET", "/computers", nil, readonlyUser())
h.List(w, r)

assert.Equal(t, http.StatusOK, w.Code)
var got []dbsqlc.ListComputersRow
_ = json.Unmarshal(w.Body.Bytes(), &got)
assert.Len(t, got, 2)
mq.AssertExpectations(t)
}

func TestComputers_List_InternalError(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewComputersHandler(mq, nil, slog.Default())

mq.On("ListComputers", mock.Anything).Return([]dbsqlc.ListComputersRow(nil), errors.New("db down"))

w, r := newRequest("GET", "/computers", nil, readonlyUser())
h.List(w, r)

assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Get ──────────────────────────────────────────────────────────────────────

func TestComputers_Get_OK(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewComputersHandler(mq, nil, slog.Default())

mq.On("GetComputerBase", mock.Anything, int64(3)).
Return(dbsqlc.Computer{ComputerID: 3, Hostname: "pc-03"}, nil)

w, r := newRequest("GET", "/computers/3", nil, readonlyUser())
r.SetPathValue("id", "3")
h.Get(w, r)

assert.Equal(t, http.StatusOK, w.Code)
mq.AssertExpectations(t)
}

func TestComputers_Get_NotFound(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewComputersHandler(mq, nil, slog.Default())

mq.On("GetComputerBase", mock.Anything, int64(99)).
Return(dbsqlc.Computer{}, pgx.ErrNoRows)

w, r := newRequest("GET", "/computers/99", nil, readonlyUser())
r.SetPathValue("id", "99")
h.Get(w, r)

assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestComputers_Get_BadID(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewComputersHandler(mq, nil, slog.Default())

w, r := newRequest("GET", "/computers/abc", nil, readonlyUser())
r.SetPathValue("id", "abc")
h.Get(w, r)

assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── Delete ───────────────────────────────────────────────────────────────────

func TestComputers_Delete_Forbidden_Readonly(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewComputersHandler(mq, nil, slog.Default())

w, r := newRequest("DELETE", "/computers/1", nil, readonlyUser())
r.SetPathValue("id", "1")
h.Delete(w, r)

assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestComputers_Delete_Forbidden_Tutor(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewComputersHandler(mq, nil, slog.Default())

w, r := newRequest("DELETE", "/computers/1", nil, tutorUser())
r.SetPathValue("id", "1")
h.Delete(w, r)

assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestComputers_Delete_BadID(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewComputersHandler(mq, nil, slog.Default())

w, r := newRequest("DELETE", "/computers/abc", nil, adminUser())
r.SetPathValue("id", "abc")
h.Delete(w, r)

assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestComputers_Delete_NotFound(t *testing.T) {
mq := new(testmock.Querier)
mockDB := new(testmock.MockDB)
h := handler.NewComputersHandler(mq, mockDB, slog.Default())

// NoRowsTx causes GetComputerBase inside the tx to return ErrNoRows → 404.
mockDB.On("Begin", mock.Anything).Return(testmock.NoRowsTx{}, nil)

w, r := newRequest("DELETE", "/computers/42", nil, adminUser())
r.SetPathValue("id", "42")
h.Delete(w, r)

assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestComputers_Delete_BeginFails(t *testing.T) {
mq := new(testmock.Querier)
mockDB := new(testmock.MockDB)
h := handler.NewComputersHandler(mq, mockDB, slog.Default())

mockDB.On("Begin", mock.Anything).Return(nil, errors.New("no conn"))

w, r := newRequest("DELETE", "/computers/1", nil, adminUser())
r.SetPathValue("id", "1")
h.Delete(w, r)

assert.Equal(t, http.StatusInternalServerError, w.Code)
}
