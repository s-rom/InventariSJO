package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"
)

type AuditHandler struct {
	queries Querier
	logger  *slog.Logger
}

func NewAuditHandler(queries Querier, logger *slog.Logger) *AuditHandler {
	return &AuditHandler{queries: queries, logger: logger}
}

// Get handles GET /audit?table=<tableName>&record_id=<id>
// Requires admin role.
func (h *AuditHandler) Get(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin) {
		return
	}

	tableName := r.URL.Query().Get("table")
	recordIDStr := r.URL.Query().Get("record_id")
	if tableName == "" || recordIDStr == "" {
		respondError(w, http.StatusBadRequest, "table and record_id query params required")
		return
	}

	recordID, err := strconv.ParseInt(recordIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid record_id")
		return
	}

	entries, err := h.queries.GetAuditLog(r.Context(), dbsqlc.GetAuditLogParams{
		TableName: tableName,
		RecordID:  recordID,
	})
	if err != nil {
		h.logger.Error("get audit log", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	respondJSON(w, http.StatusOK, entries)
}
