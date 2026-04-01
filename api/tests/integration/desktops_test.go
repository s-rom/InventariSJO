//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	dbsqlc "inventari/api/internal/db/sqlc"
	"inventari/api/internal/handler"
	"inventari/api/internal/middleware"
	"inventari/api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPasswordHash = "$2a$10$aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func makeRequest(method, target string, body any, user dbsqlc.AppUser) (*httptest.ResponseRecorder, *http.Request) {
var buf bytes.Buffer
if body != nil {
_ = json.NewEncoder(&buf).Encode(body)
}
req := httptest.NewRequest(method, target, &buf)
req.Header.Set("Content-Type", "application/json")
ctx := context.WithValue(req.Context(), middleware.CtxUser, user)
req = req.WithContext(ctx)
return httptest.NewRecorder(), req
}

// setupDesktopTest seeds the minimum required data in a fresh DB and returns
// the handler and an admin user to embed in request contexts.
func setupDesktopTest(t *testing.T) (*handler.DesktopsHandler, dbsqlc.AppUser) {
t.Helper()
pool, teardown := testutil.SetupTestDB(t)
t.Cleanup(teardown)
testutil.TruncateAll(t, pool)

ctx := context.Background()

_, err := pool.Exec(ctx,
`INSERT INTO role (role_id, description) VALUES ('admin', 'Administrador')`)
require.NoError(t, err)

_, err = pool.Exec(ctx,
`INSERT INTO app_user (username, password_hash, role_id) VALUES ($1, $2, 'admin')`,
"testadmin", testPasswordHash)
require.NoError(t, err)

var userID int64
err = pool.QueryRow(ctx, `SELECT app_user_id FROM app_user WHERE username = 'testadmin'`).Scan(&userID)
require.NoError(t, err)

user := dbsqlc.AppUser{AppUserID: userID, Username: "testadmin", RoleID: "admin"}
queries := dbsqlc.New(pool)
h := handler.NewDesktopsHandler(queries, pool, slog.Default())
return h, user
}

// TestDesktops_CreateWithoutModel_ListOK is a regression test for the NULL scan
// bug: a desktop without a desktop_model has NULL ram_gb/storage_gb/etc. from the
// COALESCE+LEFT JOIN. ListDesktops must not crash when scanning those NULLs.
func TestDesktops_CreateWithoutModel_ListOK(t *testing.T) {
h, user := setupDesktopTest(t)

w, r := makeRequest("POST", "/desktops", map[string]any{
"hostname":      "pc-test-01",
"has_wifi_card": false,
}, user)
h.Create(w, r)
require.Equal(t, http.StatusCreated, w.Code, "create: %s", w.Body.String())

var created dbsqlc.GetDesktopRow
require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
assert.Equal(t, "pc-test-01", created.Hostname)

// List must not crash even though ram_gb / storage_gb are NULL
w2, r2 := makeRequest("GET", "/desktops", nil, user)
h.List(w2, r2)
assert.Equal(t, http.StatusOK, w2.Code, "list: %s", w2.Body.String())

var listed []dbsqlc.ListDesktopsRow
require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &listed))
assert.Len(t, listed, 1)
}

// TestDesktops_Update_ObservationsSetToNull is a regression test for the bug
// where setting observations to null via PATCH left the old value in place
// because of the COALESCE in UpdateComputerBase.
func TestDesktops_Update_ObservationsSetToNull(t *testing.T) {
h, user := setupDesktopTest(t)

w, r := makeRequest("POST", "/desktops", map[string]any{
"hostname":      "pc-obs-01",
"observations":  "observació inicial",
"has_wifi_card": false,
}, user)
h.Create(w, r)
require.Equal(t, http.StatusCreated, w.Code)

var created dbsqlc.GetDesktopRow
require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
assert.True(t, created.Observations.Valid)

// Clear observations by sending null
idStr := fmt.Sprintf("%d", created.ComputerID)
w2, r2 := makeRequest("PATCH", "/desktops/"+idStr, map[string]any{"observations": nil}, user)
r2.SetPathValue("id", idStr)
h.Update(w2, r2)
require.Equal(t, http.StatusOK, w2.Code, "update: %s", w2.Body.String())

// Verify via GET
w3, r3 := makeRequest("GET", "/desktops/"+idStr, nil, user)
r3.SetPathValue("id", idStr)
h.Get(w3, r3)
require.Equal(t, http.StatusOK, w3.Code)

var updated dbsqlc.GetDesktopRow
require.NoError(t, json.Unmarshal(w3.Body.Bytes(), &updated))
assert.False(t, updated.Observations.Valid, "observations should be NULL after clearing")
}

// TestDesktops_Update_ObservationsSetToString verifies that observations can
// be changed to a new non-null string value.
func TestDesktops_Update_ObservationsSetToString(t *testing.T) {
h, user := setupDesktopTest(t)

w, r := makeRequest("POST", "/desktops", map[string]any{
"hostname":      "pc-obs-02",
"has_wifi_card": false,
}, user)
h.Create(w, r)
require.Equal(t, http.StatusCreated, w.Code)

var created dbsqlc.GetDesktopRow
require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))

idStr := fmt.Sprintf("%d", created.ComputerID)
w2, r2 := makeRequest("PATCH", "/desktops/"+idStr, map[string]any{"observations": "nova observació"}, user)
r2.SetPathValue("id", idStr)
h.Update(w2, r2)
require.Equal(t, http.StatusOK, w2.Code)

w3, r3 := makeRequest("GET", "/desktops/"+idStr, nil, user)
r3.SetPathValue("id", idStr)
h.Get(w3, r3)

var updated dbsqlc.GetDesktopRow
require.NoError(t, json.Unmarshal(w3.Body.Bytes(), &updated))
assert.True(t, updated.Observations.Valid)
assert.Equal(t, "nova observació", updated.Observations.String)
}
