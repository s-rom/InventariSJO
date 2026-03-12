package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"
)

type CentersHandler struct {
	queries *dbsqlc.Queries
	logger  *slog.Logger
}

func NewCentersHandler(queries *dbsqlc.Queries, logger *slog.Logger) *CentersHandler {
	return &CentersHandler{queries: queries, logger: logger}
}

func (h *CentersHandler) List(w http.ResponseWriter, r *http.Request) {
	centers, err := h.queries.ListCenters(r.Context())
	if err != nil {
		h.logger.Error("list centers", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, centers)
}

func (h *CentersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Name == "" {
		respondError(w, http.StatusBadRequest, "name required")
		return
	}
	center, err := h.queries.CreateCenter(r.Context(), req.Name)
	if err != nil {
		h.logger.Error("create center", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, center)
}

func (h *CentersHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	center, err := h.queries.UpdateCenter(r.Context(), dbsqlc.UpdateCenterParams{
		CenterID: id,
		Name:     req.Name,
	})
	if err != nil {
		h.logger.Error("update center", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, center)
}

func (h *CentersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.queries.DeleteCenter(r.Context(), id); err != nil {
		h.logger.Error("delete center", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}
