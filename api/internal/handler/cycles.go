package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"
)

type CyclesHandler struct {
	queries Querier
	logger  *slog.Logger
}

func NewCyclesHandler(queries Querier, logger *slog.Logger) *CyclesHandler {
	return &CyclesHandler{queries: queries, logger: logger}
}

func (h *CyclesHandler) List(w http.ResponseWriter, r *http.Request) {
	cycles, err := h.queries.ListCycles(r.Context())
	if err != nil {
		h.logger.Error("list cycles", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, cycles)
}

func (h *CyclesHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Name == "" {
		respondError(w, http.StatusBadRequest, "name required")
		return
	}
	cycle, err := h.queries.CreateCycle(r.Context(), req.Name)
	if err != nil {
		h.logger.Error("create cycle", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, cycle)
}

func (h *CyclesHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Name == "" {
		respondError(w, http.StatusBadRequest, "name required")
		return
	}
	cycle, err := h.queries.UpdateCycle(r.Context(), dbsqlc.UpdateCycleParams{
		CycleID: id,
		Name:    req.Name,
	})
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "cycle not found")
			return
		}
		h.logger.Error("update cycle", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, cycle)
}

func (h *CyclesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.queries.DeleteCycle(r.Context(), id); err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "cycle not found")
			return
		}
		h.logger.Error("delete cycle", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
