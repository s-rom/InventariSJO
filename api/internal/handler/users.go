package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	dbsqlc "inventari/api/internal/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

func toPgBool(v *bool) pgtype.Bool {
	if v == nil {
		return pgtype.Bool{}
	}
	return pgtype.Bool{Bool: *v, Valid: true}
}

type UsersHandler struct {
	queries *dbsqlc.Queries
	logger  *slog.Logger
}

func NewUsersHandler(queries *dbsqlc.Queries, logger *slog.Logger) *UsersHandler {
	return &UsersHandler{queries: queries, logger: logger}
}

func (h *UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.queries.ListUsers(r.Context())
	if err != nil {
		h.logger.Error("list users", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, users)
}

func (h *UsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		CanCreate bool   `json:"can_create"`
		CanUpdate bool   `json:"can_update"`
		CanDelete bool   `json:"can_delete"`
		IsMeta    bool   `json:"is_meta"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "username and password required")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("hash password", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	user, err := h.queries.CreateUser(r.Context(), dbsqlc.CreateUserParams{
		Username:     req.Username,
		PasswordHash: string(hash),
		CanCreate:    req.CanCreate,
		CanUpdate:    req.CanUpdate,
		CanDelete:    req.CanDelete,
		IsMeta:       req.IsMeta,
	})
	if err != nil {
		h.logger.Error("create user", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, user)
}

func (h *UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		Username  *string `json:"username"`
		CanCreate *bool   `json:"can_create"`
		CanUpdate *bool   `json:"can_update"`
		CanDelete *bool   `json:"can_delete"`
		IsMeta    *bool   `json:"is_meta"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.queries.UpdateUser(r.Context(), dbsqlc.UpdateUserParams{
		AppUserID: id,
		Username:  toPgText(req.Username),
		CanCreate: toPgBool(req.CanCreate),
		CanUpdate: toPgBool(req.CanUpdate),
		CanDelete: toPgBool(req.CanDelete),
		IsMeta:    toPgBool(req.IsMeta),
	})
	if err != nil {
		h.logger.Error("update user", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func (h *UsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.queries.DeleteUser(r.Context(), id); err != nil {
		h.logger.Error("delete user", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}
