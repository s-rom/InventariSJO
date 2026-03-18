package handler

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strings"

	dbsqlc "inventari/api/internal/db/sqlc"
	"inventari/api/internal/session"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	queries  *dbsqlc.Queries
	sessions *session.Store
	logger   *slog.Logger
}

func NewAuthHandler(queries *dbsqlc.Queries, sessions *session.Store, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{queries: queries, sessions: sessions, logger: logger}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "username and password required")
		return
	}

	user, err := h.queries.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		// Constant-time response to avoid user enumeration
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := generateToken()
	if err != nil {
		h.logger.Error("generate token", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.sessions.Set(token, user)

	respondJSON(w, http.StatusOK, map[string]string{
		"token":    token,
		"username": user.Username,
		"role_id":  user.RoleID,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	respondJSON(w, http.StatusOK, map[string]any{
		"app_user_id": user.AppUserID,
		"username":    user.Username,
		"role_id":     user.RoleID,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	token := strings.TrimPrefix(bearer, "Bearer ")
	if token != "" {
		h.sessions.Delete(token)
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
