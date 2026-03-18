package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ClassesHandler struct {
	queries *dbsqlc.Queries
	pool    *pgxpool.Pool
	logger  *slog.Logger
}

func NewClassesHandler(queries *dbsqlc.Queries, pool *pgxpool.Pool, logger *slog.Logger) *ClassesHandler {
	return &ClassesHandler{queries: queries, pool: pool, logger: logger}
}

// List returns all classes, optionally filtered by cycle via path value {cycleId}.
func (h *ClassesHandler) List(w http.ResponseWriter, r *http.Request) {
	if cycleIDStr := r.PathValue("cycleId"); cycleIDStr != "" {
		cycleID, err := strconv.ParseInt(cycleIDStr, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid cycleId")
			return
		}
		classes, err := h.queries.ListClassesByCycle(r.Context(), cycleID)
		if err != nil {
			h.logger.Error("list classes by cycle", "error", err)
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}
		respondJSON(w, http.StatusOK, classes)
		return
	}
	classes, err := h.queries.ListClasses(r.Context())
	if err != nil {
		h.logger.Error("list classes", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, classes)
}

func (h *ClassesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	class, err := h.queries.GetClass(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "class not found")
			return
		}
		h.logger.Error("get class", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, class)
}

// Create is called via POST /cycles/{cycleId}/classes
func (h *ClassesHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	cycleID, err := strconv.ParseInt(r.PathValue("cycleId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid cycleId")
		return
	}
	var req struct {
		Course         int16  `json:"course"`
		ClassLabel     string `json:"class_label"`
		Shift          string `json:"shift"`
		TutorAppUserID *int64 `json:"tutor_app_user_id"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Course == 0 || req.Shift == "" {
		respondError(w, http.StatusBadRequest, "course and shift required")
		return
	}
	if req.ClassLabel == "" {
		req.ClassLabel = "A"
	}
	class, err := h.queries.CreateClass(r.Context(), dbsqlc.CreateClassParams{
		CycleID:        cycleID,
		Course:         req.Course,
		ClassLabel:     req.ClassLabel,
		Shift:          dbsqlc.ShiftEnum(req.Shift),
		TutorAppUserID: toPgInt8(req.TutorAppUserID),
	})
	if err != nil {
		h.logger.Error("create class", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, class)
}

// Update is called via PATCH /classes/{id}
// Tutors can only update their own class; editor/admin can update any.
func (h *ClassesHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor, RoleTutor) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user := currentUser(r)

	// Tutors may only update their own class.
	if user.RoleID == RoleTutor {
		existing, err := h.queries.GetClass(r.Context(), id)
		if err != nil {
			if isNotFound(err) {
				respondError(w, http.StatusNotFound, "class not found")
				return
			}
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}
		if !existing.TutorAppUserID.Valid || existing.TutorAppUserID.Int64 != user.AppUserID {
			respondError(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	var req struct {
		ClassLabel     *string `json:"class_label"`
		Shift          *string `json:"shift"`
		TutorAppUserID *int64  `json:"tutor_app_user_id"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var shift dbsqlc.NullShiftEnum
	if req.Shift != nil {
		shift = dbsqlc.NullShiftEnum{ShiftEnum: dbsqlc.ShiftEnum(*req.Shift), Valid: true}
	}

	class, err := h.queries.UpdateClass(r.Context(), dbsqlc.UpdateClassParams{
		ClassID:        id,
		ClassLabel:     toPgText(req.ClassLabel),
		Shift:          shift,
		TutorAppUserID: toPgInt8(req.TutorAppUserID),
	})
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "class not found")
			return
		}
		h.logger.Error("update class", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, class)
}

func (h *ClassesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.queries.DeleteClass(r.Context(), id); err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "class not found")
			return
		}
		h.logger.Error("delete class", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Mine returns the classes assigned to the currently logged-in tutor.
func (h *ClassesHandler) Mine(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	classes, err := h.queries.ListClassesByTutor(r.Context(), toPgInt8(&user.AppUserID))
	if err != nil {
		h.logger.Error("list classes by tutor", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, classes)
}
