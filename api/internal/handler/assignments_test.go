package handler_test

import (
"log/slog"
"net/http"
"testing"

dbsqlc "inventari/api/internal/db/sqlc"
"inventari/api/internal/handler"
testmock "inventari/api/internal/testutil/mock"

"github.com/jackc/pgx/v5/pgtype"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/mock"
)

func TestAssignments_Create_Admin(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewAssignmentsHandler(mq, slog.Default())

mq.On("CreateAssignment", mock.Anything, mock.AnythingOfType("CreateAssignmentParams")).
Return(dbsqlc.LaptopStudentAssignment{AssignmentID: 1}, nil)

body := map[string]any{
"student_id":    1,
"class_id":      2,
"academic_year": "2024-25",
}
w, r := newRequest("POST", "/laptops/10/assignments", body, adminUser())
r.SetPathValue("laptopId", "10")
h.Create(w, r)

assert.Equal(t, http.StatusCreated, w.Code)
mq.AssertExpectations(t)
}

func TestAssignments_Create_TutorOwnsClass(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewAssignmentsHandler(mq, slog.Default())

// GetClass: tutor 3 owns class 2
mq.On("GetClass", mock.Anything, int64(2)).Return(dbsqlc.SchoolClass{
ClassID:        2,
TutorAppUserID: pgtype.Int8{Int64: 3, Valid: true},
}, nil)
mq.On("CreateAssignment", mock.Anything, mock.AnythingOfType("CreateAssignmentParams")).
Return(dbsqlc.LaptopStudentAssignment{AssignmentID: 2}, nil)

body := map[string]any{
"student_id":    1,
"class_id":      2,
"academic_year": "2024-25",
}
w, r := newRequest("POST", "/laptops/10/assignments", body, tutorUser())
r.SetPathValue("laptopId", "10")
h.Create(w, r)

assert.Equal(t, http.StatusCreated, w.Code)
mq.AssertExpectations(t)
}

func TestAssignments_Create_TutorNotOwner_Forbidden(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewAssignmentsHandler(mq, slog.Default())

mq.On("GetClass", mock.Anything, int64(2)).Return(dbsqlc.SchoolClass{
ClassID:        2,
TutorAppUserID: pgtype.Int8{Int64: 999, Valid: true},
}, nil)

body := map[string]any{
"student_id":    1,
"class_id":      2,
"academic_year": "2024-25",
}
w, r := newRequest("POST", "/laptops/10/assignments", body, tutorUser())
r.SetPathValue("laptopId", "10")
h.Create(w, r)

assert.Equal(t, http.StatusForbidden, w.Code)
mq.AssertNotCalled(t, "CreateAssignment", mock.Anything, mock.Anything)
}

func TestAssignments_Create_Readonly_Forbidden(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewAssignmentsHandler(mq, slog.Default())

body := map[string]any{
"student_id":    1,
"class_id":      2,
"academic_year": "2024-25",
}
w, r := newRequest("POST", "/laptops/10/assignments", body, readonlyUser())
r.SetPathValue("laptopId", "10")
h.Create(w, r)

assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAssignments_Create_MissingFields(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewAssignmentsHandler(mq, slog.Default())

body := map[string]any{"student_id": 1} // missing class_id, academic_year
w, r := newRequest("POST", "/laptops/10/assignments", body, adminUser())
r.SetPathValue("laptopId", "10")
h.Create(w, r)

assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAssignments_Update_TutorOwnsClass(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewAssignmentsHandler(mq, slog.Default())

// GetAssignment returns existing with class_id=2
mq.On("GetAssignment", mock.Anything, int64(5)).Return(dbsqlc.LaptopStudentAssignment{
AssignmentID: 5, ClassID: 2,
}, nil)
// tutor 3 owns class 2
mq.On("GetClass", mock.Anything, int64(2)).Return(dbsqlc.SchoolClass{
ClassID:        2,
TutorAppUserID: pgtype.Int8{Int64: 3, Valid: true},
}, nil)
mq.On("UpdateAssignment", mock.Anything, mock.AnythingOfType("UpdateAssignmentParams")).
Return(dbsqlc.LaptopStudentAssignment{AssignmentID: 5}, nil)

body := map[string]any{"academic_year": "2025-26"}
w, r := newRequest("PATCH", "/assignments/5", body, tutorUser())
r.SetPathValue("id", "5")
h.Update(w, r)

assert.Equal(t, http.StatusOK, w.Code)
mq.AssertExpectations(t)
}

func TestAssignments_Update_TutorNotOwner_Forbidden(t *testing.T) {
mq := new(testmock.Querier)
h := handler.NewAssignmentsHandler(mq, slog.Default())

mq.On("GetAssignment", mock.Anything, int64(5)).Return(dbsqlc.LaptopStudentAssignment{
AssignmentID: 5, ClassID: 2,
}, nil)
mq.On("GetClass", mock.Anything, int64(2)).Return(dbsqlc.SchoolClass{
ClassID:        2,
TutorAppUserID: pgtype.Int8{Int64: 999, Valid: true},
}, nil)

body := map[string]any{"academic_year": "2025-26"}
w, r := newRequest("PATCH", "/assignments/5", body, tutorUser())
r.SetPathValue("id", "5")
h.Update(w, r)

assert.Equal(t, http.StatusForbidden, w.Code)
mq.AssertNotCalled(t, "UpdateAssignment", mock.Anything, mock.Anything)
}
