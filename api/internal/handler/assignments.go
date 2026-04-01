package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"
)

type AssignmentsHandler struct {
	queries Querier
	logger  *slog.Logger
}

func NewAssignmentsHandler(queries Querier, logger *slog.Logger) *AssignmentsHandler {
	return &AssignmentsHandler{queries: queries, logger: logger}
}

// tutorOwnsClass returns true if the current user is the tutor of classID.
func (h *AssignmentsHandler) tutorOwnsClass(r *http.Request, classID int64) (bool, error) {
	user := currentUser(r)
	class, err := h.queries.GetClass(r.Context(), classID)
	if err != nil {
		return false, err
	}
	return class.TutorAppUserID.Valid && class.TutorAppUserID.Int64 == user.AppUserID, nil
}

// List is called via GET /laptops/{laptopId}/assignments
func (h *AssignmentsHandler) List(w http.ResponseWriter, r *http.Request) {
	laptopID, err := strconv.ParseInt(r.PathValue("laptopId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid laptopId")
		return
	}
	assignments, err := h.queries.ListAssignmentsByLaptop(r.Context(), laptopID)
	if err != nil {
		h.logger.Error("list assignments", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, assignments)
}

// ListByClass is called via GET /classes/{classId}/assignments
// Accepts optional query param ?year=2025-2026. If present, returns only
// assignments for that academic year with the laptop hostname included.
func (h *AssignmentsHandler) ListByClass(w http.ResponseWriter, r *http.Request) {
	classID, err := strconv.ParseInt(r.PathValue("classId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid classId")
		return
	}
	if year := r.URL.Query().Get("year"); year != "" {
		assignments, err := h.queries.ListAssignmentsByClassAndYear(r.Context(), dbsqlc.ListAssignmentsByClassAndYearParams{
			ClassID:      classID,
			AcademicYear: year,
		})
		if err != nil {
			h.logger.Error("list assignments by class and year", "error", err)
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}
		respondJSON(w, http.StatusOK, assignments)
		return
	}
	assignments, err := h.queries.ListAssignmentsByClass(r.Context(), classID)
	if err != nil {
		h.logger.Error("list assignments by class", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, assignments)
}

// ListByYear is called via GET /assignments?year=2025-2026
// Returns all assignments for a given academic year across all classes.
func (h *AssignmentsHandler) ListByYear(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	if year == "" {
		respondError(w, http.StatusBadRequest, "year query param required")
		return
	}
	assignments, err := h.queries.ListAssignmentsByYear(r.Context(), year)
	if err != nil {
		h.logger.Error("list assignments by year", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, assignments)
}

// ListAcademicYears is called via GET /academic-years
// Returns the distinct academic years that have at least one assignment.
func (h *AssignmentsHandler) ListAcademicYears(w http.ResponseWriter, r *http.Request) {
	years, err := h.queries.ListDistinctAcademicYears(r.Context())
	if err != nil {
		h.logger.Error("list academic years", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if years == nil {
		years = []string{}
	}
	respondJSON(w, http.StatusOK, years)
}

func (h *AssignmentsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	assignment, err := h.queries.GetAssignment(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "assignment not found")
			return
		}
		h.logger.Error("get assignment", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, assignment)
}

// Create is called via POST /laptops/{laptopId}/assignments
func (h *AssignmentsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor, RoleTutor) {
		return
	}
	laptopID, err := strconv.ParseInt(r.PathValue("laptopId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid laptopId")
		return
	}

	var req struct {
		StudentID    int64  `json:"student_id"`
		ClassID      int64  `json:"class_id"`
		AcademicYear string `json:"academic_year"`
	}
	if err := decodeJSON(r, &req); err != nil || req.StudentID == 0 || req.ClassID == 0 || req.AcademicYear == "" {
		respondError(w, http.StatusBadRequest, "student_id, class_id, and academic_year required")
		return
	}

	user := currentUser(r)
	if user.RoleID == RoleTutor {
		ok, err := h.tutorOwnsClass(r, req.ClassID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}
		if !ok {
			respondError(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	assignment, err := h.queries.CreateAssignment(r.Context(), dbsqlc.CreateAssignmentParams{
		ComputerID:   laptopID,
		StudentID:    req.StudentID,
		ClassID:      req.ClassID,
		AcademicYear: req.AcademicYear,
	})
	if err != nil {
		h.logger.Error("create assignment", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, assignment)
}

// Update is called via PATCH /assignments/{id}
func (h *AssignmentsHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor, RoleTutor) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user := currentUser(r)
	if user.RoleID == RoleTutor {
		existing, err := h.queries.GetAssignment(r.Context(), id)
		if err != nil {
			if isNotFound(err) {
				respondError(w, http.StatusNotFound, "assignment not found")
				return
			}
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}
		ok, err := h.tutorOwnsClass(r, existing.ClassID)
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
		StudentID    *int64  `json:"student_id"`
		ClassID      *int64  `json:"class_id"`
		AcademicYear *string `json:"academic_year"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	assignment, err := h.queries.UpdateAssignment(r.Context(), dbsqlc.UpdateAssignmentParams{
		AssignmentID: id,
		StudentID:    toPgInt8(req.StudentID),
		ClassID:      toPgInt8(req.ClassID),
		AcademicYear: toPgText(req.AcademicYear),
	})
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "assignment not found")
			return
		}
		h.logger.Error("update assignment", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, assignment)
}

// Delete is called via DELETE /assignments/{id}
func (h *AssignmentsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor, RoleTutor) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user := currentUser(r)
	if user.RoleID == RoleTutor {
		existing, err := h.queries.GetAssignment(r.Context(), id)
		if err != nil {
			if isNotFound(err) {
				respondError(w, http.StatusNotFound, "assignment not found")
				return
			}
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}
		ok, err := h.tutorOwnsClass(r, existing.ClassID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}
		if !ok {
			respondError(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	if err := h.queries.DeleteAssignment(r.Context(), id); err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "assignment not found")
			return
		}
		h.logger.Error("delete assignment", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
