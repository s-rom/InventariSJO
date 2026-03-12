package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DesktopsHandler struct {
	queries *dbsqlc.Queries
	pool    *pgxpool.Pool
	logger  *slog.Logger
}

func NewDesktopsHandler(queries *dbsqlc.Queries, pool *pgxpool.Pool, logger *slog.Logger) *DesktopsHandler {
	return &DesktopsHandler{queries: queries, pool: pool, logger: logger}
}

func (h *DesktopsHandler) List(w http.ResponseWriter, r *http.Request) {
	desktops, err := h.queries.ListDesktops(r.Context())
	if err != nil {
		h.logger.Error("list desktops", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, desktops)
}

func (h *DesktopsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	desktop, err := h.queries.GetDesktop(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "desktop not found")
			return
		}
		h.logger.Error("get desktop", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, desktop)
}

func (h *DesktopsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	user := currentUser(r)

	var req struct {
		Hostname        string  `json:"hostname"`
		RoomID          *int64  `json:"room_id"`
		Observations    *string `json:"observations"`
		DesktopModelID  *int64  `json:"desktop_model_id"`
		CpuID           *int64  `json:"cpu_id"`
		RamGb           *int32  `json:"ram_gb"`
		RamType         *string `json:"ram_type"`
		StorageGb       *int32  `json:"storage_gb"`
		StorageType     *string `json:"storage_type"`
		OsID            *int64  `json:"os_id"`
		EquipmentUserID *int64  `json:"equipment_user_id"`
		HasWifiCard     bool    `json:"has_wifi_card"`
		MacAddress      *string `json:"mac_address"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Hostname == "" {
		respondError(w, http.StatusBadRequest, "hostname required")
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
	qtx := h.queries.WithTx(tx)

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

	desktop, err := qtx.CreateDesktop(r.Context(), dbsqlc.CreateDesktopParams{
		ComputerID:      computer.ComputerID,
		DesktopModelID:  toPgInt8(req.DesktopModelID),
		CpuID:           toPgInt8(req.CpuID),
		RamGb:           toPgInt4(req.RamGb),
		RamType:         ramType,
		StorageGb:       toPgInt4(req.StorageGb),
		StorageType:     storType,
		OsID:            toPgInt8(req.OsID),
		EquipmentUserID: toPgInt8(req.EquipmentUserID),
		HasWifiCard:     req.HasWifiCard,
		MacAddress:      toPgText(req.MacAddress),
	})
	if err != nil {
		h.logger.Error("create desktop", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	newJSON, _ := json.Marshal(desktop)
	if err := writeAudit(r.Context(), qtx, user, auditEntry{
		TableName: "desktop",
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

	result, _ := h.queries.GetDesktop(r.Context(), computer.ComputerID)
	respondJSON(w, http.StatusCreated, result)
}

func (h *DesktopsHandler) Update(w http.ResponseWriter, r *http.Request) {
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
		DesktopModelID  *int64  `json:"desktop_model_id"`
		CpuID           *int64  `json:"cpu_id"`
		RamGb           *int32  `json:"ram_gb"`
		RamType         *string `json:"ram_type"`
		StorageGb       *int32  `json:"storage_gb"`
		StorageType     *string `json:"storage_type"`
		OsID            *int64  `json:"os_id"`
		EquipmentUserID *int64  `json:"equipment_user_id"`
		HasWifiCard     *bool   `json:"has_wifi_card"`
		MacAddress      *string `json:"mac_address"`
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
	qtx := h.queries.WithTx(tx)

	old, err := qtx.GetDesktop(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "desktop not found")
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

	updated, err := qtx.UpdateDesktop(r.Context(), dbsqlc.UpdateDesktopParams{
		ComputerID:      id,
		DesktopModelID:  toPgInt8(req.DesktopModelID),
		CpuID:           toPgInt8(req.CpuID),
		RamGb:           toPgInt4(req.RamGb),
		RamType:         ramType,
		StorageGb:       toPgInt4(req.StorageGb),
		StorageType:     storType,
		OsID:            toPgInt8(req.OsID),
		EquipmentUserID: toPgInt8(req.EquipmentUserID),
		HasWifiCard:     toPgBool(req.HasWifiCard),
		MacAddress:      toPgText(req.MacAddress),
	})
	if err != nil {
		h.logger.Error("update desktop", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	oldJSON, _ := json.Marshal(old)
	newJSON, _ := json.Marshal(updated)
	if err := writeAudit(r.Context(), qtx, user, auditEntry{
		TableName: "desktop",
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

	result, _ := h.queries.GetDesktop(r.Context(), id)
	respondJSON(w, http.StatusOK, result)
}
