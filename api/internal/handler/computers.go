package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	dbsqlc "inventari/api/internal/db/sqlc"
	"inventari/api/internal/middleware"
)

type ComputersHandler struct {
	queries *dbsqlc.Queries
	pool    *pgxpool.Pool
	logger  *slog.Logger
}

func NewComputersHandler(queries *dbsqlc.Queries, pool *pgxpool.Pool, logger *slog.Logger) *ComputersHandler {
	return &ComputersHandler{queries: queries, pool: pool, logger: logger}
}

// ComputerResponse combines the computer row with its OS list.
type ComputerResponse struct {
	dbsqlc.Computer
	OperatingSystems []dbsqlc.O `json:"operating_systems"`
}

// helpers for nullable pgtype conversions
func toPgInt8(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}

func toPgInt4(v *int32) pgtype.Int4 {
	if v == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: *v, Valid: true}
}

func toPgText(v *string) pgtype.Text {
	if v == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *v, Valid: true}
}

func (h *ComputersHandler) List(w http.ResponseWriter, r *http.Request) {
	computers, err := h.queries.ListComputers(r.Context())
	if err != nil {
		h.logger.Error("list computers", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	result := make([]ComputerResponse, 0, len(computers))
	for _, c := range computers {
		osList, err := h.queries.ListComputerOS(r.Context(), c.ComputerID)
		if err != nil {
			h.logger.Error("list computer os", "error", err, "computer_id", c.ComputerID)
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}
		result = append(result, ComputerResponse{Computer: c, OperatingSystems: osList})
	}

	respondJSON(w, http.StatusOK, result)
}

func (h *ComputersHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	computer, err := h.queries.GetComputer(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "computer not found")
		return
	}

	osList, err := h.queries.ListComputerOS(r.Context(), id)
	if err != nil {
		h.logger.Error("list computer os", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	respondJSON(w, http.StatusOK, ComputerResponse{Computer: computer, OperatingSystems: osList})
}

func (h *ComputersHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.CtxUser).(dbsqlc.AppUser)

	var req struct {
		Hostname        string  `json:"hostname"`
		CpuID           *int64  `json:"cpu_id"`
		RamGb           int32   `json:"ram_gb"`
		RamType         string  `json:"ram_type"`
		StorageGb       int32   `json:"storage_gb"`
		StorageType     string  `json:"storage_type"`
		ComputerType    string  `json:"computer_type"`
		Observations    *string `json:"observations"`
		EquipmentUserID *int64  `json:"equipment_user_id"`
		RoomID          *int64  `json:"room_id"`
		MacAddress      *string `json:"mac_address"`
		OsIDs           []int64 `json:"os_ids"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Hostname == "" || req.ComputerType == "" {
		respondError(w, http.StatusBadRequest, "hostname and computer_type required")
		return
	}

	ramType := dbsqlc.RamTypeEnumNone
	if req.RamType != "" {
		ramType = dbsqlc.RamTypeEnum(req.RamType)
	}
	storageType := dbsqlc.StorageTypeEnum(req.StorageType)
	if storageType == "" {
		storageType = dbsqlc.StorageTypeEnumSSD
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
		CpuID:              toPgInt8(req.CpuID),
		RamGb:              req.RamGb,
		RamType:            ramType,
		StorageGb:          req.StorageGb,
		StorageType:        storageType,
		ComputerType:       dbsqlc.ComputerTypeEnum(req.ComputerType),
		Observations:       toPgText(req.Observations),
		EquipmentUserID:    toPgInt8(req.EquipmentUserID),
		RoomID:             toPgInt8(req.RoomID),
		MacAddress:         toPgText(req.MacAddress),
		CreatedByAppUserID: user.AppUserID,
	})
	if err != nil {
		h.logger.Error("create computer", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	for _, osID := range req.OsIDs {
		if err := qtx.AddComputerOS(r.Context(), dbsqlc.AddComputerOSParams{
			ComputerID: computer.ComputerID,
			OsID:       osID,
		}); err != nil {
			h.logger.Error("add computer os", "error", err)
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}

	if err := tx.Commit(r.Context()); err != nil {
		h.logger.Error("commit tx", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	osList, _ := h.queries.ListComputerOS(r.Context(), computer.ComputerID)
	respondJSON(w, http.StatusCreated, ComputerResponse{Computer: computer, OperatingSystems: osList})
}

func (h *ComputersHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	user := r.Context().Value(middleware.CtxUser).(dbsqlc.AppUser)

	var req struct {
		Hostname        *string `json:"hostname"`
		CpuID           *int64  `json:"cpu_id"`
		RamGb           *int32  `json:"ram_gb"`
		RamType         *string `json:"ram_type"`
		StorageGb       *int32  `json:"storage_gb"`
		StorageType     *string `json:"storage_type"`
		ComputerType    *string `json:"computer_type"`
		Observations    *string `json:"observations"`
		EquipmentUserID *int64  `json:"equipment_user_id"`
		RoomID          *int64  `json:"room_id"`
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
	var storageType dbsqlc.NullStorageTypeEnum
	if req.StorageType != nil {
		storageType = dbsqlc.NullStorageTypeEnum{StorageTypeEnum: dbsqlc.StorageTypeEnum(*req.StorageType), Valid: true}
	}
	var computerType dbsqlc.NullComputerTypeEnum
	if req.ComputerType != nil {
		computerType = dbsqlc.NullComputerTypeEnum{ComputerTypeEnum: dbsqlc.ComputerTypeEnum(*req.ComputerType), Valid: true}
	}

	tx, err := h.pool.Begin(r.Context())
	if err != nil {
		h.logger.Error("begin tx", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	defer tx.Rollback(r.Context()) //nolint

	qtx := h.queries.WithTx(tx)

	// Snapshot before update for audit
	old, err := qtx.GetComputer(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "computer not found")
		return
	}

	computer, err := qtx.UpdateComputer(r.Context(), dbsqlc.UpdateComputerParams{
		ComputerID:      id,
		Hostname:        toPgText(req.Hostname),
		CpuID:           toPgInt8(req.CpuID),
		RamGb:           toPgInt4(req.RamGb),
		RamType:         ramType,
		StorageGb:       toPgInt4(req.StorageGb),
		StorageType:     storageType,
		ComputerType:    computerType,
		Observations:    toPgText(req.Observations),
		EquipmentUserID: toPgInt8(req.EquipmentUserID),
		RoomID:          toPgInt8(req.RoomID),
		MacAddress:      toPgText(req.MacAddress),
	})
	if err != nil {
		h.logger.Error("update computer", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	oldJSON, _ := json.Marshal(old)
	newJSON, _ := json.Marshal(computer)

	if err := qtx.InsertComputerAudit(r.Context(), dbsqlc.InsertComputerAuditParams{
		EventType:          dbsqlc.AuditEventEnumUpdated,
		ComputerID:         id,
		OldValues:          json.RawMessage(oldJSON),
		NewValues:          newJSON,
		ChangedByAppUserID: user.AppUserID,
	}); err != nil {
		h.logger.Error("insert audit", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		h.logger.Error("commit tx", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	osList, _ := h.queries.ListComputerOS(r.Context(), computer.ComputerID)
	respondJSON(w, http.StatusOK, ComputerResponse{Computer: computer, OperatingSystems: osList})
}

func (h *ComputersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	user := r.Context().Value(middleware.CtxUser).(dbsqlc.AppUser)

	tx, err := h.pool.Begin(r.Context())
	if err != nil {
		h.logger.Error("begin tx", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	defer tx.Rollback(r.Context()) //nolint

	qtx := h.queries.WithTx(tx)

	// Snapshot before delete for audit
	old, err := qtx.GetComputer(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "computer not found")
		return
	}

	oldJSON, _ := json.Marshal(old)

	if err := qtx.InsertComputerAudit(r.Context(), dbsqlc.InsertComputerAuditParams{
		EventType:          dbsqlc.AuditEventEnumDeleted,
		ComputerID:         id,
		OldValues:          json.RawMessage(oldJSON),
		NewValues:          nil,
		ChangedByAppUserID: user.AppUserID,
	}); err != nil {
		h.logger.Error("insert audit", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := qtx.DeleteComputer(r.Context(), id); err != nil {
		h.logger.Error("delete computer", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		h.logger.Error("commit tx", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

func (h *ComputersHandler) Audit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	records, err := h.queries.GetComputerAudit(r.Context(), id)
	if err != nil {
		h.logger.Error("get audit", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, records)
}
