package handler_test

import (
	"encoding/json"
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

func TestClasses_List_All(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewClassesHandler(mq, slog.Default())

	expected := []dbsqlc.ListClassesRow{{ClassID: 1}}
	mq.On("ListClasses", mock.Anything).Return(expected, nil)

	w, r := newRequest("GET", "/classes", nil, adminUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.ListClassesRow
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 1)
	mq.AssertExpectations(t)
}

func TestClasses_Update_TutorOwnsClass(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewClassesHandler(mq, slog.Default())

	// GetClass returns a class owned by tutor (AppUserID=3)
	mq.On("GetClass", mock.Anything, int64(10)).Return(dbsqlc.SchoolClass{
		ClassID:        10,
		TutorAppUserID: pgtype.Int8{Int64: 3, Valid: true},
	}, nil)
	mq.On("UpdateClass", mock.Anything, mock.AnythingOfType("UpdateClassParams")).
		Return(dbsqlc.SchoolClass{ClassID: 10}, nil)

	body := map[string]any{"class_label": "B"}
	w, r := newRequest("PATCH", "/classes/10", body, tutorUser())
	r.SetPathValue("id", "10")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestClasses_Update_TutorNotOwner_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewClassesHandler(mq, slog.Default())

	mq.On("GetClass", mock.Anything, int64(10)).Return(dbsqlc.SchoolClass{
		ClassID:        10,
		TutorAppUserID: pgtype.Int8{Int64: 999, Valid: true}, // different tutor
	}, nil)

	body := map[string]any{"class_label": "B"}
	w, r := newRequest("PATCH", "/classes/10", body, tutorUser())
	r.SetPathValue("id", "10")
	h.Update(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mq.AssertNotCalled(t, "UpdateClass", mock.Anything, mock.Anything)
}

func TestClasses_Update_Admin_Bypass(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewClassesHandler(mq, slog.Default())

	// Admin should NOT call GetClass to check ownership
	mq.On("UpdateClass", mock.Anything, mock.AnythingOfType("UpdateClassParams")).
		Return(dbsqlc.SchoolClass{ClassID: 10}, nil)

	body := map[string]any{"class_label": "C"}
	w, r := newRequest("PATCH", "/classes/10", body, adminUser())
	r.SetPathValue("id", "10")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertNotCalled(t, "GetClass", mock.Anything, mock.Anything)
}

func TestClasses_Mine(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewClassesHandler(mq, slog.Default())

	expected := []dbsqlc.ListClassesByTutorRow{{ClassID: 5}}
	tutorID := pgtype.Int8{Int64: 3, Valid: true}
	mq.On("ListClassesByTutor", mock.Anything, tutorID).Return(expected, nil)

	w, r := newRequest("GET", "/tutor/classes", nil, tutorUser())
	h.Mine(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestClasses_Delete_OnlyAdmin(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewClassesHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/classes/1", nil, editorUser())
	r.SetPathValue("id", "1")
	h.Delete(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
