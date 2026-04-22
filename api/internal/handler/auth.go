package handler

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strings"

	dbsqlc "inventari/api/internal/db/sqlc"
	"inventari/api/internal/session"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	queries  Querier
	sessions *session.Store
	logger   *slog.Logger
}

func NewAuthHandler(queries Querier, sessions *session.Store, logger *slog.Logger) *AuthHandler {
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

	if !user.PasswordHash.Valid {
		// Google-only account: cannot log in with password
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(req.Password)); err != nil {
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

// ChangePassword is called via POST /auth/change-password
// Any authenticated user can change their own password.
// Requires current_password (for verification) and new_password.
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := decodeJSON(r, &req); err != nil || req.CurrentPassword == "" || req.NewPassword == "" {
		respondError(w, http.StatusBadRequest, "current_password and new_password required")
		return
	}
	if len(req.NewPassword) < 8 {
		respondError(w, http.StatusBadRequest, "new password must be at least 8 characters")
		return
	}

	user := currentUser(r)

	// Re-fetch to get the stored hash (currentUser comes from session, no hash there)
	dbUser, err := h.queries.GetUserByUsername(r.Context(), user.Username)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if !dbUser.PasswordHash.Valid {
		respondError(w, http.StatusBadRequest, "account uses Google login, cannot change password")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash.String), []byte(req.CurrentPassword)); err != nil {
		respondError(w, http.StatusUnauthorized, "current password is incorrect")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("bcrypt generate", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if err := h.queries.UpdateUserPassword(r.Context(), dbsqlc.UpdateUserPasswordParams{
		AppUserID:    user.AppUserID,
		PasswordHash: pgtype.Text{String: string(hash), Valid: true},
	}); err != nil {
		h.logger.Error("update password", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "password updated"})
}
