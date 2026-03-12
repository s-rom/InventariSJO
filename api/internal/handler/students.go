package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentsHandler struct {
	queries *dbsqlc.Queries
	pool    *pgxpool.Pool
	logger  *slog.Logger
}

func NewStudentsHandler(queries *dbsqlc.Queries, pool *pgxpool.Pool, logger *slog.Logger) *StudentsHandler {
	return &StudentsHandler{queries: queries, pool: pool, logger: logger}
}

// isTutorOfClass checks if the given user is the tutor of the class
// that contains the student with studentID.
func (h *StudentsHandler) isTutorOfClass(r *http.Request, classID int64) (bool, error) {
	user := currentUser(r)
	class, err := h.queries.GetClass(r.Context(), classID)
	if err != nil {
		return false, err
	}
	return class.TutorAppUserID.Valid && class.TutorAppUserID.Int64 == user.AppUserID, nil
}

// List is called via GET /classes/{classId}/students
func (h *StudentsHandler) List(w http.ResponseWriter, r *http.Request) {
	classID, err := strconv.ParseInt(r.PathValue("classId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid classId")
		return
	}
	students, err := h.queries.ListStudentsByClass(r.Context(), classID)
	if err != nil {
		h.logger.Error("list students", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, students)
}

func (h *StudentsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	student, err := h.queries.GetStudent(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "student not found")
			return
		}
		h.logger.Error("get student", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, student)
}

// Create is called via POST /classes/{classId}/students
func (h *StudentsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor, RoleTutor) {
		return
	}
	classID, err := strconv.ParseInt(r.PathValue("classId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid classId")
		return
	}

	// Tutors can only add students to their own class.
	user := currentUser(r)
	if user.RoleID == RoleTutor {
		ok, err := h.isTutorOfClass(r, classID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}
		if !ok {
			respondError(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	var req struct {
		FullName string `json:"full_name"`
	}
	if err := decodeJSON(r, &req); err != nil || req.FullName == "" {
		respondError(w, http.StatusBadRequest, "full_name required")
		return
	}
	student, err := h.queries.CreateStudent(r.Context(), dbsqlc.CreateStudentParams{
		FullName: req.FullName,
		ClassID:  classID,
	})
	if err != nil {
		h.logger.Error("create student", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, student)
}

// Update is called via PATCH /students/{id}
func (h *StudentsHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor, RoleTutor) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Tutors: can only update students from their own class.
	user := currentUser(r)
	if user.RoleID == RoleTutor {
		existing, err := h.queries.GetStudent(r.Context(), id)
		if err != nil {
			if isNotFound(err) {
				respondError(w, http.StatusNotFound, "student not found")
				return
			}
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}
		ok, err := h.isTutorOfClass(r, existing.ClassID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}
		if !ok {
			respondError(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	var req struct {
		FullName *string `json:"full_name"`
		ClassID  *int64  `json:"class_id"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	student, err := h.queries.UpdateStudent(r.Context(), dbsqlc.UpdateStudentParams{
		StudentID: id,
		FullName:  toPgText(req.FullName),
		ClassID:   toPgInt8(req.ClassID),
	})
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "student not found")
			return
		}
		h.logger.Error("update student", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, student)
}

// Delete is called via DELETE /students/{id}
func (h *StudentsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.queries.DeleteStudent(r.Context(), id); err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "student not found")
			return
		}
		h.logger.Error("delete student", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
