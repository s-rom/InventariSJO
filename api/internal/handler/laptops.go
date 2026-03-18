package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"
)

type LaptopsHandler struct {
	queries Querier
	pool    DB
	logger  *slog.Logger
}

func NewLaptopsHandler(queries Querier, pool DB, logger *slog.Logger) *LaptopsHandler {
	return &LaptopsHandler{queries: queries, pool: pool, logger: logger}
}

func (h *LaptopsHandler) List(w http.ResponseWriter, r *http.Request) {
	laptops, err := h.queries.ListLaptops(r.Context())
	if err != nil {
		h.logger.Error("list laptops", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, laptops)
}

func (h *LaptopsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	laptop, err := h.queries.GetLaptop(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "laptop not found")
			return
		}
		h.logger.Error("get laptop", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, laptop)
}

func (h *LaptopsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	user := currentUser(r)

	var req struct {
		Hostname        string  `json:"hostname"`
		RoomID          *int64  `json:"room_id"`
		Observations    *string `json:"observations"`
		LaptopModelID   int64   `json:"laptop_model_id"`
		SerialNumber    *string `json:"serial_number"`
		RamGb           *int32  `json:"ram_gb"`
		RamType         *string `json:"ram_type"`
		StorageGb       *int32  `json:"storage_gb"`
		StorageType     *string `json:"storage_type"`
		MacAddress      *string `json:"mac_address"`
		OsID            *int64  `json:"os_id"`
		EquipmentUserID *int64  `json:"equipment_user_id"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Hostname == "" || req.LaptopModelID == 0 {
		respondError(w, http.StatusBadRequest, "hostname and laptop_model_id required")
		return
	}

	var ramType dbsqlc.NullRamTypeEnum
	if req.RamType != nil {
		ramType = dbsqlc.NullRamTypeEnum{RamTypeEnum: dbsqlc.RamTypeEnum(*req.RamType), Valid: true}
	}
	var storType dbsqlc.NullStorageTypeEnum
	if req.StorageType != nil {
		storType = dbsqlc.NullStorageTypeEnum{StorageTypeEnum: dbsqlc.StorageTypeEnum(*req.StorageType), Valid: true}
	}

	tx, err := h.pool.Begin(r.Context())
	if err != nil {
		h.logger.Error("begin tx", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	defer tx.Rollback(r.Context()) //nolint
	qtx := dbsqlc.New(tx)

	computer, err := qtx.CreateComputer(r.Context(), dbsqlc.CreateComputerParams{
		Hostname:           req.Hostname,
		RoomID:             toPgInt8(req.RoomID),
		Observations:       toPgText(req.Observations),
		CreatedByAppUserID: user.AppUserID,
	})
	if err != nil {
		h.logger.Error("create computer", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	laptop, err := qtx.CreateLaptop(r.Context(), dbsqlc.CreateLaptopParams{
		ComputerID:      computer.ComputerID,
		LaptopModelID:   req.LaptopModelID,
		SerialNumber:    toPgText(req.SerialNumber),
		RamGb:           toPgInt4(req.RamGb),
		RamType:         ramType,
		StorageGb:       toPgInt4(req.StorageGb),
		StorageType:     storType,
		MacAddress:      toPgText(req.MacAddress),
		OsID:            toPgInt8(req.OsID),
		EquipmentUserID: toPgInt8(req.EquipmentUserID),
	})
	if err != nil {
		h.logger.Error("create laptop", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	newJSON, _ := json.Marshal(laptop)
	if err := writeAudit(r.Context(), qtx, user, auditEntry{
		TableName: "laptop",
		RecordID:  computer.ComputerID,
		EventType: dbsqlc.AuditEventEnumCreated,
		NewValues: newJSON,
	}); err != nil {
		h.logger.Error("write audit", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		h.logger.Error("commit tx", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	result, _ := h.queries.GetLaptop(r.Context(), computer.ComputerID)
	respondJSON(w, http.StatusCreated, result)
}

func (h *LaptopsHandler) Update(w http.ResponseWriter, r *http.Request) {
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
		Hostname        *string `json:"hostname"`
		RoomID          *int64  `json:"room_id"`
		Observations    *string `json:"observations"`
		LaptopModelID   *int64  `json:"laptop_model_id"`
		SerialNumber    *string `json:"serial_number"`
		RamGb           *int32  `json:"ram_gb"`
		RamType         *string `json:"ram_type"`
		StorageGb       *int32  `json:"storage_gb"`
		StorageType     *string `json:"storage_type"`
		MacAddress      *string `json:"mac_address"`
		OsID            *int64  `json:"os_id"`
		EquipmentUserID *int64  `json:"equipment_user_id"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var ramType dbsqlc.NullRamTypeEnum
	if req.RamType != nil {
		ramType = dbsqlc.NullRamTypeEnum{RamTypeEnum: dbsqlc.RamTypeEnum(*req.RamType), Valid: true}
	}
	var storType dbsqlc.NullStorageTypeEnum
	if req.StorageType != nil {
		storType = dbsqlc.NullStorageTypeEnum{StorageTypeEnum: dbsqlc.StorageTypeEnum(*req.StorageType), Valid: true}
	}

	tx, err := h.pool.Begin(r.Context())
	if err != nil {
		h.logger.Error("begin tx", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	defer tx.Rollback(r.Context()) //nolint
	qtx := dbsqlc.New(tx)

	old, err := qtx.GetLaptop(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "laptop not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if _, err := qtx.UpdateComputerBase(r.Context(), dbsqlc.UpdateComputerBaseParams{
		ComputerID:   id,
		Hostname:     toPgText(req.Hostname),
		RoomID:       toPgInt8(req.RoomID),
		Observations: toPgText(req.Observations),
	}); err != nil {
		h.logger.Error("update computer base", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	updated, err := qtx.UpdateLaptop(r.Context(), dbsqlc.UpdateLaptopParams{
		ComputerID:      id,
		LaptopModelID:   toPgInt8(req.LaptopModelID),
		SerialNumber:    toPgText(req.SerialNumber),
		RamGb:           toPgInt4(req.RamGb),
		RamType:         ramType,
		StorageGb:       toPgInt4(req.StorageGb),
		StorageType:     storType,
		MacAddress:      toPgText(req.MacAddress),
		OsID:            toPgInt8(req.OsID),
		EquipmentUserID: toPgInt8(req.EquipmentUserID),
	})
	if err != nil {
		h.logger.Error("update laptop", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	oldJSON, _ := json.Marshal(old)
	newJSON, _ := json.Marshal(updated)
	if err := writeAudit(r.Context(), qtx, user, auditEntry{
		TableName: "laptop",
		RecordID:  id,
		EventType: dbsqlc.AuditEventEnumUpdated,
		OldValues: oldJSON,
		NewValues: newJSON,
	}); err != nil {
		h.logger.Error("write audit", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		h.logger.Error("commit tx", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	result, _ := h.queries.GetLaptop(r.Context(), id)
	respondJSON(w, http.StatusOK, result)
}
