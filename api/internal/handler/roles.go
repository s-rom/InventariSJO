package handler

import (
	"log/slog"
	"net/http"

	dbsqlc "inventari/api/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RolesHandler struct {
	queries *dbsqlc.Queries
	pool    *pgxpool.Pool
	logger  *slog.Logger
}

func NewRolesHandler(queries *dbsqlc.Queries, pool *pgxpool.Pool, logger *slog.Logger) *RolesHandler {
	return &RolesHandler{queries: queries, pool: pool, logger: logger}
}

func (h *RolesHandler) List(w http.ResponseWriter, r *http.Request) {
	roles, err := h.queries.ListRoles(r.Context())
	if err != nil {
		h.logger.Error("list roles", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, roles)
}

func (h *RolesHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin) {
		return
	}
	var req struct {
		RoleID      string  `json:"role_id"`
		Description *string `json:"description"`
	}
	if err := decodeJSON(r, &req); err != nil || req.RoleID == "" {
		respondError(w, http.StatusBadRequest, "role_id required")
		return
	}
	role, err := h.queries.CreateRole(r.Context(), dbsqlc.CreateRoleParams{
		RoleID:      req.RoleID,
		Description: toPgText(req.Description),
	})
	if err != nil {
		h.logger.Error("create role", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, role)
}

func (h *RolesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin) {
		return
	}
	id := r.PathValue("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "role id required")
		return
	}
	if err := h.queries.DeleteRole(r.Context(), id); err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "role not found")
			return
		}
		h.logger.Error("delete role", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
