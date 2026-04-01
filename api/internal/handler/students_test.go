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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─── List ─────────────────────────────────────────────────────────────────────

func TestStudents_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	mq.On("ListStudentsByClass", mock.Anything, int64(2)).Return([]dbsqlc.Student{
		{StudentID: 1, FullName: "Anna Garcia", ClassID: 2},
		{StudentID: 2, FullName: "Pau Mas", ClassID: 2},
	}, nil)

	w, r := newRequest("GET", "/classes/2/students", nil, readonlyUser())
	r.SetPathValue("classId", "2")
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.Student
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 2)
	mq.AssertExpectations(t)
}

func TestStudents_List_BadClassID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	w, r := newRequest("GET", "/classes/abc/students", nil, readonlyUser())
	r.SetPathValue("classId", "abc")
	h.List(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStudents_List_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	mq.On("ListStudentsByClass", mock.Anything, int64(1)).Return([]dbsqlc.Student(nil), errors.New("db down"))

	w, r := newRequest("GET", "/classes/1/students", nil, readonlyUser())
	r.SetPathValue("classId", "1")
	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Get ──────────────────────────────────────────────────────────────────────

func TestStudents_Get_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	mq.On("GetStudent", mock.Anything, int64(5)).Return(dbsqlc.Student{StudentID: 5, FullName: "Anna", ClassID: 1}, nil)

	w, r := newRequest("GET", "/students/5", nil, readonlyUser())
	r.SetPathValue("id", "5")
	h.Get(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestStudents_Get_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	w, r := newRequest("GET", "/students/abc", nil, readonlyUser())
	r.SetPathValue("id", "abc")
	h.Get(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStudents_Get_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	mq.On("GetStudent", mock.Anything, int64(99)).Return(dbsqlc.Student{}, pgx.ErrNoRows)

	w, r := newRequest("GET", "/students/99", nil, readonlyUser())
	r.SetPathValue("id", "99")
	h.Get(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── Create ───────────────────────────────────────────────────────────────────

func TestStudents_Create_OK_Admin(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	mq.On("CreateStudent", mock.Anything, dbsqlc.CreateStudentParams{FullName: "Marta Puig", ClassID: 3}).
		Return(dbsqlc.Student{StudentID: 10, FullName: "Marta Puig", ClassID: 3}, nil)

	w, r := newRequest("POST", "/classes/3/students", map[string]any{"full_name": "Marta Puig"}, adminUser())
	r.SetPathValue("classId", "3")
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestStudents_Create_Forbidden_Readonly(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	w, r := newRequest("POST", "/classes/1/students", map[string]any{"full_name": "X"}, readonlyUser())
	r.SetPathValue("classId", "1")
	h.Create(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestStudents_Create_Tutor_OwnClass(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	// tutorUser() has AppUserID = 3; class is owned by tutor 3
	mq.On("GetClass", mock.Anything, int64(7)).Return(dbsqlc.SchoolClass{
		ClassID:        7,
		TutorAppUserID: pgtype.Int8{Int64: 3, Valid: true},
	}, nil)
	mq.On("CreateStudent", mock.Anything, dbsqlc.CreateStudentParams{FullName: "Joan Gil", ClassID: 7}).
		Return(dbsqlc.Student{StudentID: 11, FullName: "Joan Gil", ClassID: 7}, nil)

	w, r := newRequest("POST", "/classes/7/students", map[string]any{"full_name": "Joan Gil"}, tutorUser())
	r.SetPathValue("classId", "7")
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestStudents_Create_Tutor_OtherClass(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	// Class belongs to tutor 99, not tutor 3
	mq.On("GetClass", mock.Anything, int64(7)).Return(dbsqlc.SchoolClass{
		ClassID:        7,
		TutorAppUserID: pgtype.Int8{Int64: 99, Valid: true},
	}, nil)

	w, r := newRequest("POST", "/classes/7/students", map[string]any{"full_name": "X"}, tutorUser())
	r.SetPathValue("classId", "7")
	h.Create(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestStudents_Create_EmptyName(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	w, r := newRequest("POST", "/classes/1/students", map[string]any{"full_name": ""}, adminUser())
	r.SetPathValue("classId", "1")
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── Update ───────────────────────────────────────────────────────────────────

func TestStudents_Update_OK_Editor(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	mq.On("UpdateStudent", mock.Anything, mock.AnythingOfType("UpdateStudentParams")).
		Return(dbsqlc.Student{StudentID: 5, FullName: "Updated Name", ClassID: 1}, nil)

	w, r := newRequest("PATCH", "/students/5", map[string]any{"full_name": "Updated Name"}, editorUser())
	r.SetPathValue("id", "5")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestStudents_Update_Forbidden_Readonly(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	w, r := newRequest("PATCH", "/students/5", map[string]any{"full_name": "X"}, readonlyUser())
	r.SetPathValue("id", "5")
	h.Update(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestStudents_Update_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	w, r := newRequest("PATCH", "/students/abc", map[string]any{"full_name": "X"}, adminUser())
	r.SetPathValue("id", "abc")
	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStudents_Update_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	mq.On("UpdateStudent", mock.Anything, mock.AnythingOfType("UpdateStudentParams")).
		Return(dbsqlc.Student{}, pgx.ErrNoRows)

	w, r := newRequest("PATCH", "/students/99", map[string]any{"full_name": "Ghost"}, editorUser())
	r.SetPathValue("id", "99")
	h.Update(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestStudents_Update_Tutor_OwnClass(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	// tutorUser() AppUserID = 3
	mq.On("GetStudent", mock.Anything, int64(5)).Return(dbsqlc.Student{StudentID: 5, FullName: "Anna", ClassID: 2}, nil)
	mq.On("GetClass", mock.Anything, int64(2)).Return(dbsqlc.SchoolClass{
		ClassID:        2,
		TutorAppUserID: pgtype.Int8{Int64: 3, Valid: true},
	}, nil)
	mq.On("UpdateStudent", mock.Anything, mock.AnythingOfType("UpdateStudentParams")).
		Return(dbsqlc.Student{StudentID: 5, FullName: "Anna B", ClassID: 2}, nil)

	w, r := newRequest("PATCH", "/students/5", map[string]any{"full_name": "Anna B"}, tutorUser())
	r.SetPathValue("id", "5")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestStudents_Update_Tutor_OtherClass(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	mq.On("GetStudent", mock.Anything, int64(5)).Return(dbsqlc.Student{StudentID: 5, FullName: "Anna", ClassID: 2}, nil)
	mq.On("GetClass", mock.Anything, int64(2)).Return(dbsqlc.SchoolClass{
		ClassID:        2,
		TutorAppUserID: pgtype.Int8{Int64: 99, Valid: true},
	}, nil)

	w, r := newRequest("PATCH", "/students/5", map[string]any{"full_name": "X"}, tutorUser())
	r.SetPathValue("id", "5")
	h.Update(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ─── Delete ───────────────────────────────────────────────────────────────────

func TestStudents_Delete_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	mq.On("DeleteStudent", mock.Anything, int64(8)).Return(nil)

	w, r := newRequest("DELETE", "/students/8", nil, adminUser())
	r.SetPathValue("id", "8")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}

func TestStudents_Delete_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/students/8", nil, editorUser())
	r.SetPathValue("id", "8")
	h.Delete(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestStudents_Delete_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/students/abc", nil, adminUser())
	r.SetPathValue("id", "abc")
	h.Delete(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStudents_Delete_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewStudentsHandler(mq, slog.Default())

	mq.On("DeleteStudent", mock.Anything, int64(99)).Return(pgx.ErrNoRows)

	w, r := newRequest("DELETE", "/students/99", nil, adminUser())
	r.SetPathValue("id", "99")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
