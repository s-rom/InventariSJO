package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"
)

type EquipmentUsersHandler struct {
	queries *dbsqlc.Queries
	logger  *slog.Logger
}

func NewEquipmentUsersHandler(queries *dbsqlc.Queries, logger *slog.Logger) *EquipmentUsersHandler {
	return &EquipmentUsersHandler{queries: queries, logger: logger}
}

func (h *EquipmentUsersHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.queries.ListEquipmentUsers(r.Context())
	if err != nil {
		h.logger.Error("list equipment users", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, list)
}

func (h *EquipmentUsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Name == "" {
		respondError(w, http.StatusBadRequest, "name required")
		return
	}
	eu, err := h.queries.CreateEquipmentUser(r.Context(), req.Name)
	if err != nil {
		h.logger.Error("create equipment user", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, eu)
}

func (h *EquipmentUsersHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	eu, err := h.queries.UpdateEquipmentUser(r.Context(), dbsqlc.UpdateEquipmentUserParams{
		EquipmentUserID: id,
		Name:            req.Name,
	})
	if err != nil {
		h.logger.Error("update equipment user", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, eu)
}

func (h *EquipmentUsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.queries.DeleteEquipmentUser(r.Context(), id); err != nil {
		h.logger.Error("delete equipment user", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}
