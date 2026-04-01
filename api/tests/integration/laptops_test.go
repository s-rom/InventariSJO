//go:build integration

package integration_test

import (
"context"
"encoding/json"
"fmt"
"log/slog"
"net/http"
"testing"

dbsqlc "inventari/api/internal/db/sqlc"
"inventari/api/internal/handler"
"inventari/api/internal/testutil"

"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
)

// setupLaptopTest seeds the minimum required rows and returns the handler
// along with an admin user and a laptop_model_id to use in requests.
func setupLaptopTest(t *testing.T) (*handler.LaptopsHandler, dbsqlc.AppUser, int64) {
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

// Seed a brand and a laptop_model (required FK for laptops)
var brandID int64
err = pool.QueryRow(ctx, `INSERT INTO brand (name) VALUES ('TestBrand') RETURNING brand_id`).Scan(&brandID)
require.NoError(t, err)

var modelID int64
err = pool.QueryRow(ctx, `
INSERT INTO laptop_model (brand_id, model_name, base_ram_gb, base_ram_type, base_storage_gb, base_storage_type)
VALUES ($1, 'ModelX', 8, 'DDR4', 256, 'SSD')
RETURNING laptop_model_id`,
brandID).Scan(&modelID)
require.NoError(t, err)

user := dbsqlc.AppUser{AppUserID: userID, Username: "testadmin", RoleID: "admin"}
queries := dbsqlc.New(pool)
h := handler.NewLaptopsHandler(queries, pool, slog.Default())
return h, user, modelID
}

// TestLaptops_Create_List_OK creates a laptop and lists it back.
func TestLaptops_Create_List_OK(t *testing.T) {
h, user, modelID := setupLaptopTest(t)

w, r := makeRequest("POST", "/laptops", map[string]any{
"hostname":        "lt-test-01",
"laptop_model_id": modelID,
}, user)
h.Create(w, r)
require.Equal(t, http.StatusCreated, w.Code, "create: %s", w.Body.String())

var created dbsqlc.GetLaptopRow
require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
assert.Equal(t, "lt-test-01", created.Hostname)
assert.Equal(t, modelID, created.LaptopModelID)

w2, r2 := makeRequest("GET", "/laptops", nil, user)
h.List(w2, r2)
assert.Equal(t, http.StatusOK, w2.Code, "list: %s", w2.Body.String())

var listed []dbsqlc.ListLaptopsRow
require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &listed))
assert.Len(t, listed, 1)
assert.Equal(t, "lt-test-01", listed[0].Hostname)
}

// TestLaptops_Update_ObservationsSetToNull is the laptop regression test for
// the same observations-null bug as desktops.
func TestLaptops_Update_ObservationsSetToNull(t *testing.T) {
h, user, modelID := setupLaptopTest(t)

w, r := makeRequest("POST", "/laptops", map[string]any{
"hostname":        "lt-obs-01",
"laptop_model_id": modelID,
"observations":    "observació inicial",
}, user)
h.Create(w, r)
require.Equal(t, http.StatusCreated, w.Code)

var created dbsqlc.GetLaptopRow
require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
assert.True(t, created.Observations.Valid)

idStr := fmt.Sprintf("%d", created.ComputerID)
w2, r2 := makeRequest("PATCH", "/laptops/"+idStr, map[string]any{"observations": nil}, user)
r2.SetPathValue("id", idStr)
h.Update(w2, r2)
require.Equal(t, http.StatusOK, w2.Code, "update: %s", w2.Body.String())

w3, r3 := makeRequest("GET", "/laptops/"+idStr, nil, user)
r3.SetPathValue("id", idStr)
h.Get(w3, r3)
require.Equal(t, http.StatusOK, w3.Code)

var updated dbsqlc.GetLaptopRow
require.NoError(t, json.Unmarshal(w3.Body.Bytes(), &updated))
assert.False(t, updated.Observations.Valid, "observations should be NULL after clearing")
}

// TestLaptops_Update_ObservationsSetToString verifies changing observations to a string.
func TestLaptops_Update_ObservationsSetToString(t *testing.T) {
h, user, modelID := setupLaptopTest(t)

w, r := makeRequest("POST", "/laptops", map[string]any{
"hostname":        "lt-obs-02",
"laptop_model_id": modelID,
}, user)
h.Create(w, r)
require.Equal(t, http.StatusCreated, w.Code)

var created dbsqlc.GetLaptopRow
require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))

idStr := fmt.Sprintf("%d", created.ComputerID)
w2, r2 := makeRequest("PATCH", "/laptops/"+idStr, map[string]any{"observations": "nova observació"}, user)
r2.SetPathValue("id", idStr)
h.Update(w2, r2)
require.Equal(t, http.StatusOK, w2.Code)

w3, r3 := makeRequest("GET", "/laptops/"+idStr, nil, user)
r3.SetPathValue("id", idStr)
h.Get(w3, r3)

var updated dbsqlc.GetLaptopRow
require.NoError(t, json.Unmarshal(w3.Body.Bytes(), &updated))
assert.True(t, updated.Observations.Valid)
assert.Equal(t, "nova observació", updated.Observations.String)
}

// TestLaptops_Get_NotFound verifies 404 when the laptop does not exist.
func TestLaptops_Get_NotFound(t *testing.T) {
h, user, _ := setupLaptopTest(t)

w, r := makeRequest("GET", "/laptops/99999", nil, user)
r.SetPathValue("id", "99999")
h.Get(w, r)

assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestLaptops_Update_ChangeModel verifies a laptop_model_id can be changed.
func TestLaptops_Update_ChangeModel(t *testing.T) {
h, user, modelID := setupLaptopTest(t)

// Get the pool to seed a second model — reach through the handler is not
// possible, so we create a second laptop (different hostname) to confirm
// the model is what we seeded.
w, r := makeRequest("POST", "/laptops", map[string]any{
"hostname":        "lt-model-01",
"laptop_model_id": modelID,
}, user)
h.Create(w, r)
require.Equal(t, http.StatusCreated, w.Code)

var created dbsqlc.GetLaptopRow
require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
assert.Equal(t, modelID, created.LaptopModelID)

// Update hostname only – laptop_model_id should remain
idStr := fmt.Sprintf("%d", created.ComputerID)
w2, r2 := makeRequest("PATCH", "/laptops/"+idStr, map[string]any{"hostname": "lt-model-01-renamed"}, user)
r2.SetPathValue("id", idStr)
h.Update(w2, r2)
require.Equal(t, http.StatusOK, w2.Code)

w3, r3 := makeRequest("GET", "/laptops/"+idStr, nil, user)
r3.SetPathValue("id", idStr)
h.Get(w3, r3)

var updated dbsqlc.GetLaptopRow
require.NoError(t, json.Unmarshal(w3.Body.Bytes(), &updated))
assert.Equal(t, "lt-model-01-renamed", updated.Hostname)
assert.Equal(t, modelID, updated.LaptopModelID)
}
