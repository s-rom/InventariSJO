package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ComputersHandler struct {
	queries *dbsqlc.Queries
	pool    *pgxpool.Pool
	logger  *slog.Logger
}

func NewComputersHandler(queries *dbsqlc.Queries, pool *pgxpool.Pool, logger *slog.Logger) *ComputersHandler {
	return &ComputersHandler{queries: queries, pool: pool, logger: logger}
}

// List returns all computers (desktops + laptops) with their type and base fields.
func (h *ComputersHandler) List(w http.ResponseWriter, r *http.Request) {
	computers, err := h.queries.ListComputers(r.Context())
	if err != nil {
		h.logger.Error("list computers", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, computers)
}

// Get returns the base fields of a single computer plus its type.
func (h *ComputersHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	computer, err := h.queries.GetComputerBase(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "computer not found")
			return
		}
		h.logger.Error("get computer", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, computer)
}

// Delete removes a computer and cascades to desktop/laptop subtable.
func (h *ComputersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	user := currentUser(r)

	tx, err := h.pool.Begin(r.Context())
	if err != nil {
		h.logger.Error("begin tx", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	defer tx.Rollback(r.Context()) //nolint

	qtx := h.queries.WithTx(tx)

	old, err := qtx.GetComputerBase(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "computer not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := qtx.DeleteComputer(r.Context(), id); err != nil {
		h.logger.Error("delete computer", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	oldJSON, _ := json.Marshal(old)
	if err := writeAudit(r.Context(), qtx, user, auditEntry{
		TableName: "computer",
		RecordID:  id,
		EventType: dbsqlc.AuditEventEnumDeleted,
		OldValues: oldJSON,
	}); err != nil {
		h.logger.Error("write audit", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		h.logger.Error("commit tx", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
