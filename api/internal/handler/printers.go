package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"
)

// ── Printer Models ────────────────────────────────────────────────────────────

type PrinterModelsHandler struct {
	queries Querier
	logger  *slog.Logger
}

func NewPrinterModelsHandler(queries Querier, logger *slog.Logger) *PrinterModelsHandler {
	return &PrinterModelsHandler{queries: queries, logger: logger}
}

func (h *PrinterModelsHandler) List(w http.ResponseWriter, r *http.Request) {
	models, err := h.queries.ListPrinterModels(r.Context())
	if err != nil {
		h.logger.Error("list printer models", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, models)
}

func (h *PrinterModelsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := h.queries.GetPrinterModel(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "printer model not found")
			return
		}
		h.logger.Error("get printer model", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, m)
}

func (h *PrinterModelsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	var req struct {
		BrandID     int64  `json:"brand_id"`
		ModelName   string `json:"model_name"`
		PrinterType string `json:"printer_type"`
		PrintColor  string `json:"print_color"`
	}
	if err := decodeJSON(r, &req); err != nil || req.BrandID == 0 || req.ModelName == "" || req.PrinterType == "" {
		respondError(w, http.StatusBadRequest, "brand_id, model_name and printer_type required")
		return
	}
	m, err := h.queries.CreatePrinterModel(r.Context(), dbsqlc.CreatePrinterModelParams{
		BrandID:     req.BrandID,
		ModelName:   req.ModelName,
		PrinterType: dbsqlc.PrinterTypeEnum(req.PrinterType),
		PrintColor:  dbsqlc.PrintColorEnum(req.PrintColor),
	})
	if err != nil {
		h.logger.Error("create printer model", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, m)
}

func (h *PrinterModelsHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		BrandID     *int64  `json:"brand_id"`
		ModelName   *string `json:"model_name"`
		PrinterType *string `json:"printer_type"`
		PrintColor  *string `json:"print_color"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var pt dbsqlc.NullPrinterTypeEnum
	if req.PrinterType != nil {
		pt = dbsqlc.NullPrinterTypeEnum{PrinterTypeEnum: dbsqlc.PrinterTypeEnum(*req.PrinterType), Valid: true}
	}
	var pc dbsqlc.NullPrintColorEnum
	if req.PrintColor != nil {
		pc = dbsqlc.NullPrintColorEnum{PrintColorEnum: dbsqlc.PrintColorEnum(*req.PrintColor), Valid: true}
	}

	m, err := h.queries.UpdatePrinterModel(r.Context(), dbsqlc.UpdatePrinterModelParams{
		PrinterModelID: id,
		BrandID:        toPgInt8(req.BrandID),
		ModelName:      toPgText(req.ModelName),
		PrinterType:    pt,
		PrintColor:     pc,
	})
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "printer model not found")
			return
		}
		h.logger.Error("update printer model", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, m)
}

func (h *PrinterModelsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.queries.DeletePrinterModel(r.Context(), id); err != nil {
		h.logger.Error("delete printer model", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

// ── Printer Supplies (consumibles) ────────────────────────────────────────────

type PrinterSuppliesHandler struct {
	queries Querier
	logger  *slog.Logger
}

func NewPrinterSuppliesHandler(queries Querier, logger *slog.Logger) *PrinterSuppliesHandler {
	return &PrinterSuppliesHandler{queries: queries, logger: logger}
}

func (h *PrinterSuppliesHandler) List(w http.ResponseWriter, r *http.Request) {
	supplies, err := h.queries.ListPrinterSupplies(r.Context())
	if err != nil {
		h.logger.Error("list printer supplies", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, supplies)
}

func (h *PrinterSuppliesHandler) ListByModel(w http.ResponseWriter, r *http.Request) {
	modelID, err := strconv.ParseInt(r.PathValue("modelId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid modelId")
		return
	}
	supplies, err := h.queries.ListSuppliesByPrinterModel(r.Context(), modelID)
	if err != nil {
		h.logger.Error("list supplies by printer model", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, supplies)
}

func (h *PrinterSuppliesHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	var req struct {
		Name       string `json:"name"`
		SupplyType string `json:"supply_type"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Name == "" || req.SupplyType == "" {
		respondError(w, http.StatusBadRequest, "name and supply_type required")
		return
	}
	s, err := h.queries.CreatePrinterSupply(r.Context(), dbsqlc.CreatePrinterSupplyParams{
		Name:       req.Name,
		SupplyType: dbsqlc.PrinterTypeEnum(req.SupplyType),
	})
	if err != nil {
		h.logger.Error("create printer supply", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, s)
}

func (h *PrinterSuppliesHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Name       *string `json:"name"`
		SupplyType *string `json:"supply_type"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var st dbsqlc.NullPrinterTypeEnum
	if req.SupplyType != nil {
		st = dbsqlc.NullPrinterTypeEnum{PrinterTypeEnum: dbsqlc.PrinterTypeEnum(*req.SupplyType), Valid: true}
	}

	s, err := h.queries.UpdatePrinterSupply(r.Context(), dbsqlc.UpdatePrinterSupplyParams{
		PrinterSupplyID: id,
		Name:            toPgText(req.Name),
		SupplyType:      st,
	})
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "printer supply not found")
			return
		}
		h.logger.Error("update printer supply", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, s)
}

func (h *PrinterSuppliesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.queries.DeletePrinterSupply(r.Context(), id); err != nil {
		h.logger.Error("delete printer supply", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

// AddToModel links a supply to a printer model.
func (h *PrinterSuppliesHandler) AddToModel(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	modelID, err := strconv.ParseInt(r.PathValue("modelId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid modelId")
		return
	}
	var req struct {
		PrinterSupplyID int64 `json:"printer_supply_id"`
	}
	if err := decodeJSON(r, &req); err != nil || req.PrinterSupplyID == 0 {
		respondError(w, http.StatusBadRequest, "printer_supply_id required")
		return
	}
	if err := h.queries.AddPrinterModelSupply(r.Context(), dbsqlc.AddPrinterModelSupplyParams{
		PrinterModelID:  modelID,
		PrinterSupplyID: req.PrinterSupplyID,
	}); err != nil {
		h.logger.Error("add printer model supply", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

// RemoveFromModel unlinks a supply from a printer model.
func (h *PrinterSuppliesHandler) RemoveFromModel(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	modelID, err := strconv.ParseInt(r.PathValue("modelId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid modelId")
		return
	}
	supplyID, err := strconv.ParseInt(r.PathValue("supplyId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid supplyId")
		return
	}
	if err := h.queries.RemovePrinterModelSupply(r.Context(), dbsqlc.RemovePrinterModelSupplyParams{
		PrinterModelID:  modelID,
		PrinterSupplyID: supplyID,
	}); err != nil {
		h.logger.Error("remove printer model supply", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

// ── Printers ──────────────────────────────────────────────────────────────────

type PrintersHandler struct {
	queries Querier
	pool    DB
	logger  *slog.Logger
}

func NewPrintersHandler(queries Querier, pool DB, logger *slog.Logger) *PrintersHandler {
	return &PrintersHandler{queries: queries, pool: pool, logger: logger}
}

func (h *PrintersHandler) List(w http.ResponseWriter, r *http.Request) {
	printers, err := h.queries.ListPrinters(r.Context())
	if err != nil {
		h.logger.Error("list printers", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, printers)
}

func (h *PrintersHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	p, err := h.queries.GetPrinter(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "printer not found")
			return
		}
		h.logger.Error("get printer", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, p)
}

func (h *PrintersHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	user := currentUser(r)

	var req struct {
		PrinterModelID       int64   `json:"printer_model_id"`
		Status               string  `json:"status"`
		HasNetworkCapability bool    `json:"has_network_capability"`
		UsesNetwork          bool    `json:"uses_network"`
		IPAddress            *string `json:"ip_address"`
		RoomID               *int64  `json:"room_id"`
		EquipmentUserID      *int64  `json:"equipment_user_id"`
		Observations         *string `json:"observations"`
	}
	if err := decodeJSON(r, &req); err != nil || req.PrinterModelID == 0 {
		respondError(w, http.StatusBadRequest, "printer_model_id required")
		return
	}
	if req.Status == "" {
		req.Status = "actiu"
	}

	p, err := h.queries.CreatePrinter(r.Context(), dbsqlc.CreatePrinterParams{
		PrinterModelID:       req.PrinterModelID,
		Status:               dbsqlc.DeviceStatusEnum(req.Status),
		HasNetworkCapability: req.HasNetworkCapability,
		UsesNetwork:          req.UsesNetwork,
		IpAddress:            toPgText(req.IPAddress),
		RoomID:               toPgInt8(req.RoomID),
		EquipmentUserID:      toPgInt8(req.EquipmentUserID),
		Observations:         toPgText(req.Observations),
		CreatedByAppUserID:   user.AppUserID,
	})
	if err != nil {
		h.logger.Error("create printer", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	newJSON, _ := json.Marshal(p)
	if err := writeAudit(r.Context(), h.queries, user, auditEntry{
		TableName: "printer",
		RecordID:  p.PrinterID,
		EventType: dbsqlc.AuditEventEnumCreated,
		NewValues: newJSON,
	}); err != nil {
		h.logger.Error("write audit", "error", err)
	}

	result, _ := h.queries.GetPrinter(r.Context(), p.PrinterID)
	respondJSON(w, http.StatusCreated, result)
}

func (h *PrintersHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	user := currentUser(r)

	var req struct {
		PrinterModelID       *int64  `json:"printer_model_id"`
		Status               *string `json:"status"`
		HasNetworkCapability *bool   `json:"has_network_capability"`
		UsesNetwork          *bool   `json:"uses_network"`
		IPAddress            *string `json:"ip_address"`
		RoomID               *int64  `json:"room_id"`
		EquipmentUserID      *int64  `json:"equipment_user_id"`
		Observations         *string `json:"observations"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	old, err := h.queries.GetPrinter(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "printer not found")
			return
		}
		h.logger.Error("get printer for update", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	var st dbsqlc.NullDeviceStatusEnum
	if req.Status != nil {
		st = dbsqlc.NullDeviceStatusEnum{DeviceStatusEnum: dbsqlc.DeviceStatusEnum(*req.Status), Valid: true}
	}

	p, err := h.queries.UpdatePrinter(r.Context(), dbsqlc.UpdatePrinterParams{
		PrinterID:            id,
		PrinterModelID:       toPgInt8(req.PrinterModelID),
		Status:               st,
		HasNetworkCapability: toPgBool(req.HasNetworkCapability),
		UsesNetwork:          toPgBool(req.UsesNetwork),
		IpAddress:            toPgText(req.IPAddress),
		RoomID:               toPgInt8(req.RoomID),
		EquipmentUserID:      toPgInt8(req.EquipmentUserID),
		Observations:         toPgText(req.Observations),
	})
	if err != nil {
		h.logger.Error("update printer", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	oldJSON, _ := json.Marshal(old)
	newJSON, _ := json.Marshal(p)
	if err := writeAudit(r.Context(), h.queries, user, auditEntry{
		TableName: "printer",
		RecordID:  id,
		EventType: dbsqlc.AuditEventEnumUpdated,
		OldValues: oldJSON,
		NewValues: newJSON,
	}); err != nil {
		h.logger.Error("write audit", "error", err)
	}

	result, _ := h.queries.GetPrinter(r.Context(), id)
	respondJSON(w, http.StatusOK, result)
}

func (h *PrintersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	user := currentUser(r)

	old, err := h.queries.GetPrinter(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "printer not found")
			return
		}
		h.logger.Error("get printer for delete", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := h.queries.DeletePrinter(r.Context(), id); err != nil {
		h.logger.Error("delete printer", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	oldJSON, _ := json.Marshal(old)
	if err := writeAudit(r.Context(), h.queries, user, auditEntry{
		TableName: "printer",
		RecordID:  id,
		EventType: dbsqlc.AuditEventEnumDeleted,
		OldValues: oldJSON,
	}); err != nil {
		h.logger.Error("write audit", "error", err)
	}

	respondJSON(w, http.StatusNoContent, nil)
}
