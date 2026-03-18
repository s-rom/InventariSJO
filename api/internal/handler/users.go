package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"

	"golang.org/x/crypto/bcrypt"
)

type UsersHandler struct {
	queries Querier
	logger  *slog.Logger
}

func NewUsersHandler(queries Querier, logger *slog.Logger) *UsersHandler {
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
	if !requireRole(w, r, RoleAdmin) {
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		RoleID   string `json:"role_id"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Username == "" || req.Password == "" || req.RoleID == "" {
		respondError(w, http.StatusBadRequest, "username, password and role_id required")
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
		RoleID:       req.RoleID,
	})
	if err != nil {
		h.logger.Error("create user", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, user)
}

func (h *UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin) {
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		Username *string `json:"username"`
		RoleID   *string `json:"role_id"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.queries.UpdateUser(r.Context(), dbsqlc.UpdateUserParams{
		AppUserID: id,
		Username:  toPgText(req.Username),
		RoleID:    toPgText(req.RoleID),
	})
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		h.logger.Error("update user", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func (h *UsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin) {
		return
	}

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
	w.WriteHeader(http.StatusNoContent)
}
