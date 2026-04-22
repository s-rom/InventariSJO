package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	dbsqlc "inventari/api/internal/db/sqlc"
	"inventari/api/internal/session"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleAuthHandler handles the Google OAuth 2.0 flow.
type GoogleAuthHandler struct {
	queries       Querier
	sessions      *session.Store
	logger        *slog.Logger
	oauthConfig   *oauth2.Config
	allowedDomain string
	frontendURL   string
	states        oauthStateStore
}

// oauthStateStore is a short-lived in-memory map of CSRF state tokens.
type oauthStateStore struct {
	mu     sync.Mutex
	states map[string]time.Time
}

func (s *oauthStateStore) set(state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Expire old entries on every write to avoid unbounded growth
	now := time.Now()
	for k, v := range s.states {
		if now.After(v) {
			delete(s.states, k)
		}
	}
	s.states[state] = now.Add(10 * time.Minute)
}

func (s *oauthStateStore) verify(state string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	exp, ok := s.states[state]
	if !ok {
		return false
	}
	delete(s.states, state)
	return time.Now().Before(exp)
}

// NewGoogleAuthHandler reads OAuth config from environment variables:
//
//	GOOGLE_CLIENT_ID      — OAuth client ID from Google Cloud Console
//	GOOGLE_CLIENT_SECRET  — OAuth client secret
//	GOOGLE_ALLOWED_DOMAIN — corporate domain to restrict (e.g. "escola.cat")
//	GOOGLE_CALLBACK_URL   — full callback URL registered in Google Cloud Console
//	FRONTEND_URL          — base URL of the frontend app (for final redirect)
func NewGoogleAuthHandler(queries Querier, sessions *session.Store, logger *slog.Logger) *GoogleAuthHandler {
	cfg := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_CALLBACK_URL"),
		Scopes: []string{
			"openid",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return &GoogleAuthHandler{
		queries:       queries,
		sessions:      sessions,
		logger:        logger,
		oauthConfig:   cfg,
		allowedDomain: os.Getenv("GOOGLE_ALLOWED_DOMAIN"),
		frontendURL:   os.Getenv("FRONTEND_URL"),
		states:        oauthStateStore{states: make(map[string]time.Time)},
	}
}

// Redirect initiates the OAuth flow by redirecting the user to Google.
// GET /auth/google
func (h *GoogleAuthHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	state, err := generateToken() // reuses the helper from auth.go
	if err != nil {
		h.logger.Error("generate oauth state", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	h.states.set(state)
	url := h.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Callback handles the redirect from Google after the user grants consent.
// GET /auth/google/callback
func (h *GoogleAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	errRedirect := func(reason string) {
		http.Redirect(w, r, h.frontendURL+"/login?error="+reason, http.StatusTemporaryRedirect)
	}

	if !h.states.verify(r.URL.Query().Get("state")) {
		errRedirect("invalid_state")
		return
	}

	token, err := h.oauthConfig.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		h.logger.Error("oauth token exchange", "error", err)
		errRedirect("oauth_exchange")
		return
	}

	userInfo, err := fetchGoogleUserInfo(r.Context(), h.oauthConfig, token)
	if err != nil {
		h.logger.Error("fetch google userinfo", "error", err)
		errRedirect("userinfo")
		return
	}

	if h.allowedDomain != "" && !strings.HasSuffix(userInfo.Email, "@"+h.allowedDomain) {
		h.logger.Warn("google login rejected: domain not allowed", "email", userInfo.Email)
		errRedirect("domain_not_allowed")
		return
	}

	user, err := h.findOrCreateUser(r.Context(), userInfo)
	if err != nil {
		h.logger.Error("find or create google user", "error", err)
		errRedirect("user_error")
		return
	}

	sessionToken, err := generateToken()
	if err != nil {
		h.logger.Error("generate session token", "error", err)
		errRedirect("internal")
		return
	}

	h.sessions.Set(sessionToken, user)

	redirectURL := fmt.Sprintf("%s/login?token=%s&role=%s&username=%s",
		h.frontendURL, sessionToken, user.RoleID, user.Username)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

type googleUserInfo struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func fetchGoogleUserInfo(ctx context.Context, cfg *oauth2.Config, token *oauth2.Token) (*googleUserInfo, error) {
	client := cfg.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var info googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// findOrCreateUser looks up the user by google_sub, then by email (and links
// google_sub), and finally auto-provisions a new user with 'readonly' role.
func (h *GoogleAuthHandler) findOrCreateUser(ctx context.Context, info *googleUserInfo) (dbsqlc.AppUser, error) {
	// 1. Known Google user
	user, err := h.queries.GetUserByGoogleSub(ctx, pgtype.Text{String: info.Sub, Valid: true})
	if err == nil {
		return user, nil
	}

	// 2. Existing user with matching email (first Google login → link account)
	user, err = h.queries.GetUserByEmail(ctx, pgtype.Text{String: info.Email, Valid: true})
	if err == nil {
		return h.queries.SetGoogleSub(ctx, dbsqlc.SetGoogleSubParams{
			GoogleSub: pgtype.Text{String: info.Sub, Valid: true},
			Email:     pgtype.Text{String: info.Email, Valid: true},
		})
	}

	// 3. Auto-provision: create user with 'readonly' role
	// Username is derived from the email local part; a suffix is appended on conflict by the DB UNIQUE constraint.
	username := strings.Split(info.Email, "@")[0]
	return h.queries.CreateGoogleUser(ctx, dbsqlc.CreateGoogleUserParams{
		Username:  username,
		Email:     pgtype.Text{String: info.Email, Valid: true},
		GoogleSub: pgtype.Text{String: info.Sub, Valid: true},
	})
}
