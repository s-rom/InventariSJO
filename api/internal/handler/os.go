package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"
)

type OSHandler struct {
	queries Querier
	logger  *slog.Logger
}

func NewOSHandler(queries Querier, logger *slog.Logger) *OSHandler {
	return &OSHandler{queries: queries, logger: logger}
}

func (h *OSHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.queries.ListOS(r.Context())
	if err != nil {
		h.logger.Error("list os", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, list)
}

func (h *OSHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Name == "" {
		respondError(w, http.StatusBadRequest, "name required")
		return
	}
	entry, err := h.queries.CreateOS(r.Context(), req.Name)
	if err != nil {
		h.logger.Error("create os", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, entry)
}

func (h *OSHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	entry, err := h.queries.UpdateOS(r.Context(), dbsqlc.UpdateOSParams{
		OsID: id,
		Name: req.Name,
	})
	if err != nil {
		h.logger.Error("update os", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, entry)
}

func (h *OSHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.queries.DeleteOS(r.Context(), id); err != nil {
		h.logger.Error("delete os", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}
