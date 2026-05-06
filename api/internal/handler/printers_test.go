package handler_test

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"testing"

	dbsqlc "inventari/api/internal/db/sqlc"
	"inventari/api/internal/handler"
	testmock "inventari/api/internal/testutil/mock"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ═══════════════════════════════════════════════════════════════
// PrinterModels
// ═══════════════════════════════════════════════════════════════

func TestPrinterModels_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterModelsHandler(mq, slog.Default())

	rows := []dbsqlc.ListPrinterModelsRow{{PrinterModelID: 1, ModelName: "LaserJet"}}
	mq.On("ListPrinterModels", mock.Anything).Return(rows, nil)

	w, r := newRequest("GET", "/printer-models", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.ListPrinterModelsRow
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 1)
	mq.AssertExpectations(t)
}

func TestPrinterModels_List_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterModelsHandler(mq, slog.Default())

	mq.On("ListPrinterModels", mock.Anything).Return([]dbsqlc.ListPrinterModelsRow(nil), errors.New("db down"))

	w, r := newRequest("GET", "/printer-models", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPrinterModels_Get_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterModelsHandler(mq, slog.Default())

	mq.On("GetPrinterModel", mock.Anything, int64(3)).
		Return(dbsqlc.GetPrinterModelRow{PrinterModelID: 3, ModelName: "M404n"}, nil)

	w, r := newRequest("GET", "/printer-models/3", nil, readonlyUser())
	r.SetPathValue("id", "3")
	h.Get(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestPrinterModels_Get_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterModelsHandler(mq, slog.Default())

	mq.On("GetPrinterModel", mock.Anything, int64(99)).
		Return(dbsqlc.GetPrinterModelRow{}, pgx.ErrNoRows)

	w, r := newRequest("GET", "/printer-models/99", nil, readonlyUser())
	r.SetPathValue("id", "99")
	h.Get(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPrinterModels_Get_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterModelsHandler(mq, slog.Default())

	w, r := newRequest("GET", "/printer-models/abc", nil, readonlyUser())
	r.SetPathValue("id", "abc")
	h.Get(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPrinterModels_Create_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterModelsHandler(mq, slog.Default())

	created := dbsqlc.PrinterModel{PrinterModelID: 1, ModelName: "LaserJet", PrinterType: "toner", PrintColor: "Color"}
	mq.On("CreatePrinterModel", mock.Anything, mock.MatchedBy(func(p dbsqlc.CreatePrinterModelParams) bool {
		return p.BrandID == 2 && p.ModelName == "LaserJet"
	})).Return(created, nil)

	body := map[string]any{"brand_id": 2, "model_name": "LaserJet", "printer_type": "toner", "print_color": "Color"}
	w, r := newRequest("POST", "/printer-models", body, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestPrinterModels_Create_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterModelsHandler(mq, slog.Default())

	w, r := newRequest("POST", "/printer-models", map[string]any{"brand_id": 1, "model_name": "x", "printer_type": "toner"}, readonlyUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPrinterModels_Create_MissingFields(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterModelsHandler(mq, slog.Default())

	// missing model_name
	w, r := newRequest("POST", "/printer-models", map[string]any{"brand_id": 1, "printer_type": "toner"}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPrinterModels_Update_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterModelsHandler(mq, slog.Default())

	updated := dbsqlc.PrinterModel{PrinterModelID: 1, ModelName: "NewName"}
	mq.On("UpdatePrinterModel", mock.Anything, mock.Anything).Return(updated, nil)

	body := map[string]any{"model_name": "NewName"}
	w, r := newRequest("PATCH", "/printer-models/1", body, adminUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestPrinterModels_Update_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterModelsHandler(mq, slog.Default())

	mq.On("UpdatePrinterModel", mock.Anything, mock.Anything).Return(dbsqlc.PrinterModel{}, pgx.ErrNoRows)

	w, r := newRequest("PATCH", "/printer-models/99", map[string]any{"model_name": "x"}, adminUser())
	r.SetPathValue("id", "99")
	h.Update(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPrinterModels_Update_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterModelsHandler(mq, slog.Default())

	w, r := newRequest("PATCH", "/printer-models/1", map[string]any{"model_name": "x"}, readonlyUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPrinterModels_Delete_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterModelsHandler(mq, slog.Default())

	mq.On("DeletePrinterModel", mock.Anything, int64(1)).Return(nil)

	w, r := newRequest("DELETE", "/printer-models/1", nil, adminUser())
	r.SetPathValue("id", "1")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}

func TestPrinterModels_Delete_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterModelsHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/printer-models/1", nil, editorUser())
	r.SetPathValue("id", "1")
	h.Delete(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ═══════════════════════════════════════════════════════════════
// PrinterSupplies
// ═══════════════════════════════════════════════════════════════

func TestPrinterSupplies_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterSuppliesHandler(mq, slog.Default())

	rows := []dbsqlc.PrinterSupply{{PrinterSupplyID: 1, Name: "Toner HP"}}
	mq.On("ListPrinterSupplies", mock.Anything).Return(rows, nil)

	w, r := newRequest("GET", "/printer-supplies", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestPrinterSupplies_ListByModel_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterSuppliesHandler(mq, slog.Default())

	rows := []dbsqlc.PrinterSupply{{PrinterSupplyID: 2, Name: "Ink Cyan"}}
	mq.On("ListSuppliesByPrinterModel", mock.Anything, int64(5)).Return(rows, nil)

	w, r := newRequest("GET", "/printer-models/5/supplies", nil, readonlyUser())
	r.SetPathValue("modelId", "5")
	h.ListByModel(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.PrinterSupply
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 1)
	mq.AssertExpectations(t)
}

func TestPrinterSupplies_ListByModel_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterSuppliesHandler(mq, slog.Default())

	w, r := newRequest("GET", "/printer-models/abc/supplies", nil, readonlyUser())
	r.SetPathValue("modelId", "abc")
	h.ListByModel(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPrinterSupplies_Create_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterSuppliesHandler(mq, slog.Default())

	created := dbsqlc.PrinterSupply{PrinterSupplyID: 1, Name: "Toner Negre", SupplyType: "toner"}
	mq.On("CreatePrinterSupply", mock.Anything, mock.MatchedBy(func(p dbsqlc.CreatePrinterSupplyParams) bool {
		return p.Name == "Toner Negre" && p.SupplyType == "toner"
	})).Return(created, nil)

	body := map[string]any{"name": "Toner Negre", "supply_type": "toner"}
	w, r := newRequest("POST", "/printer-supplies", body, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestPrinterSupplies_Create_MissingFields(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterSuppliesHandler(mq, slog.Default())

	w, r := newRequest("POST", "/printer-supplies", map[string]any{"name": "x"}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPrinterSupplies_AddToModel_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterSuppliesHandler(mq, slog.Default())

	mq.On("AddPrinterModelSupply", mock.Anything, dbsqlc.AddPrinterModelSupplyParams{
		PrinterModelID: 3, PrinterSupplyID: 7,
	}).Return(nil)

	body := map[string]any{"printer_supply_id": 7}
	w, r := newRequest("POST", "/printer-models/3/supplies", body, adminUser())
	r.SetPathValue("modelId", "3")
	h.AddToModel(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}

func TestPrinterSupplies_AddToModel_MissingSupplyID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterSuppliesHandler(mq, slog.Default())

	w, r := newRequest("POST", "/printer-models/3/supplies", map[string]any{}, adminUser())
	r.SetPathValue("modelId", "3")
	h.AddToModel(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPrinterSupplies_RemoveFromModel_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterSuppliesHandler(mq, slog.Default())

	mq.On("RemovePrinterModelSupply", mock.Anything, dbsqlc.RemovePrinterModelSupplyParams{
		PrinterModelID: 3, PrinterSupplyID: 7,
	}).Return(nil)

	w, r := newRequest("DELETE", "/printer-models/3/supplies/7", nil, adminUser())
	r.SetPathValue("modelId", "3")
	r.SetPathValue("supplyId", "7")
	h.RemoveFromModel(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}

func TestPrinterSupplies_RemoveFromModel_BadSupplyID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrinterSuppliesHandler(mq, slog.Default())

	w, r := newRequest("DELETE", "/printer-models/3/supplies/abc", nil, adminUser())
	r.SetPathValue("modelId", "3")
	r.SetPathValue("supplyId", "abc")
	h.RemoveFromModel(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ═══════════════════════════════════════════════════════════════
// Printers
// ═══════════════════════════════════════════════════════════════

func TestPrinters_List_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrintersHandler(mq, nil, slog.Default())

	rows := []dbsqlc.ListPrintersRow{
		{PrinterID: 1, PrinterModelID: 2},
		{PrinterID: 2, PrinterModelID: 2},
	}
	mq.On("ListPrinters", mock.Anything).Return(rows, nil)

	w, r := newRequest("GET", "/printers", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dbsqlc.ListPrintersRow
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 2)
	mq.AssertExpectations(t)
}

func TestPrinters_List_InternalError(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrintersHandler(mq, nil, slog.Default())

	mq.On("ListPrinters", mock.Anything).Return([]dbsqlc.ListPrintersRow(nil), errors.New("db down"))

	w, r := newRequest("GET", "/printers", nil, readonlyUser())
	h.List(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPrinters_Get_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrintersHandler(mq, nil, slog.Default())

	mq.On("GetPrinter", mock.Anything, int64(5)).
		Return(dbsqlc.GetPrinterRow{PrinterID: 5}, nil)

	w, r := newRequest("GET", "/printers/5", nil, readonlyUser())
	r.SetPathValue("id", "5")
	h.Get(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestPrinters_Get_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrintersHandler(mq, nil, slog.Default())

	mq.On("GetPrinter", mock.Anything, int64(99)).
		Return(dbsqlc.GetPrinterRow{}, pgx.ErrNoRows)

	w, r := newRequest("GET", "/printers/99", nil, readonlyUser())
	r.SetPathValue("id", "99")
	h.Get(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPrinters_Get_BadID(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrintersHandler(mq, nil, slog.Default())

	w, r := newRequest("GET", "/printers/abc", nil, readonlyUser())
	r.SetPathValue("id", "abc")
	h.Get(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPrinters_Create_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrintersHandler(mq, nil, slog.Default())

	body := map[string]any{"printer_model_id": 1}
	w, r := newRequest("POST", "/printers", body, readonlyUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPrinters_Create_MissingModel(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrintersHandler(mq, nil, slog.Default())

	w, r := newRequest("POST", "/printers", map[string]any{}, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPrinters_Create_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrintersHandler(mq, nil, slog.Default())

	created := dbsqlc.Printer{PrinterID: 10, PrinterModelID: 2}
	mq.On("CreatePrinter", mock.Anything, mock.MatchedBy(func(p dbsqlc.CreatePrinterParams) bool {
		return p.PrinterModelID == 2
	})).Return(created, nil)
	mq.On("InsertAuditLog", mock.Anything, mock.Anything).Return(nil)
	mq.On("GetPrinter", mock.Anything, int64(10)).Return(dbsqlc.GetPrinterRow{PrinterID: 10}, nil)

	body := map[string]any{"printer_model_id": 2, "status": "actiu"}
	w, r := newRequest("POST", "/printers", body, adminUser())
	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mq.AssertExpectations(t)
}

func TestPrinters_Update_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrintersHandler(mq, nil, slog.Default())

	w, r := newRequest("PATCH", "/printers/1", map[string]any{}, readonlyUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPrinters_Update_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrintersHandler(mq, nil, slog.Default())

	mq.On("GetPrinter", mock.Anything, int64(99)).Return(dbsqlc.GetPrinterRow{}, pgx.ErrNoRows)

	w, r := newRequest("PATCH", "/printers/99", map[string]any{}, adminUser())
	r.SetPathValue("id", "99")
	h.Update(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mq.AssertExpectations(t)
}

func TestPrinters_Update_ClearsObservations(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrintersHandler(mq, nil, slog.Default())

	old := dbsqlc.GetPrinterRow{PrinterID: 1}
	updated := dbsqlc.Printer{PrinterID: 1}
	final := dbsqlc.GetPrinterRow{PrinterID: 1}

	mq.On("GetPrinter", mock.Anything, int64(1)).Return(old, nil).Once()
	mq.On("UpdatePrinter", mock.Anything, mock.MatchedBy(func(p dbsqlc.UpdatePrinterParams) bool {
		// observations must be NULL (not Valid) when sent as null from client
		return p.PrinterID == 1 && !p.Observations.Valid
	})).Return(updated, nil)
	mq.On("InsertAuditLog", mock.Anything, mock.Anything).Return(nil)
	mq.On("GetPrinter", mock.Anything, int64(1)).Return(final, nil).Once()

	// client sends observations: null explicitly
	body := map[string]any{"observations": nil}
	w, r := newRequest("PATCH", "/printers/1", body, adminUser())
	r.SetPathValue("id", "1")
	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mq.AssertExpectations(t)
}

func TestPrinters_Delete_Forbidden(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrintersHandler(mq, nil, slog.Default())

	w, r := newRequest("DELETE", "/printers/1", nil, editorUser())
	r.SetPathValue("id", "1")
	h.Delete(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPrinters_Delete_NotFound(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrintersHandler(mq, nil, slog.Default())

	mq.On("GetPrinter", mock.Anything, int64(99)).Return(dbsqlc.GetPrinterRow{}, pgx.ErrNoRows)

	w, r := newRequest("DELETE", "/printers/99", nil, adminUser())
	r.SetPathValue("id", "99")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPrinters_Delete_OK(t *testing.T) {
	mq := new(testmock.Querier)
	h := handler.NewPrintersHandler(mq, nil, slog.Default())

	old := dbsqlc.GetPrinterRow{PrinterID: 1}
	mq.On("GetPrinter", mock.Anything, int64(1)).Return(old, nil)
	mq.On("DeletePrinter", mock.Anything, int64(1)).Return(nil)
	mq.On("InsertAuditLog", mock.Anything, mock.Anything).Return(nil)

	w, r := newRequest("DELETE", "/printers/1", nil, adminUser())
	r.SetPathValue("id", "1")
	h.Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mq.AssertExpectations(t)
}
