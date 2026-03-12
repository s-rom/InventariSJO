package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	dbsqlc "inventari/api/internal/db/sqlc"
	"inventari/api/internal/middleware"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ── HTTP helpers ─────────────────────────────────────────────────────────────

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}

func decodeJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

// isNotFound returns true when the error represents a missing row.
func isNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

// ── pgtype helpers ────────────────────────────────────────────────────────────

func toPgInt8(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}

func toPgInt4(v *int32) pgtype.Int4 {
	if v == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: *v, Valid: true}
}

func toPgText(v *string) pgtype.Text {
	if v == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *v, Valid: true}
}

func toPgBool(v *bool) pgtype.Bool {
	if v == nil {
		return pgtype.Bool{}
	}
	return pgtype.Bool{Bool: *v, Valid: true}
}

// ── RBAC helpers ──────────────────────────────────────────────────────────────

const (
	RoleReadonly = "readonly"
	RoleEditor   = "editor"
	RoleAdmin    = "admin"
	RoleTutor    = "tutor"
)

// currentUser extracts the authenticated user from context.
func currentUser(r *http.Request) dbsqlc.AppUser {
	return r.Context().Value(middleware.CtxUser).(dbsqlc.AppUser)
}

// requireRole checks that the current user has one of the given roles.
// Returns false (and writes 403) if the check fails.
func requireRole(w http.ResponseWriter, r *http.Request, roles ...string) bool {
	user := currentUser(r)
	for _, role := range roles {
		if user.RoleID == role {
			return true
		}
	}
	respondError(w, http.StatusForbidden, "insufficient permissions")
	return false
}

// ── Audit helper ──────────────────────────────────────────────────────────────

type auditEntry struct {
	TableName string
	RecordID  int64
	EventType dbsqlc.AuditEventEnum
	OldValues json.RawMessage // nil for 'created'
	NewValues json.RawMessage // nil for 'deleted'
}

func writeAudit(ctx context.Context, q *dbsqlc.Queries, user dbsqlc.AppUser, e auditEntry) error {
	return q.InsertAuditLog(ctx, dbsqlc.InsertAuditLogParams{
		TableName:          e.TableName,
		RecordID:           e.RecordID,
		EventType:          e.EventType,
		OldValues:          e.OldValues,
		NewValues:          e.NewValues,
		ChangedByAppUserID: user.AppUserID,
		ChangedByUsername:  user.Username,
	})
}
