package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"
)

// ============================================================
// Projector Models
// ============================================================

type ProjectorModelsHandler struct {
	queries Querier
	logger  *slog.Logger
}

func NewProjectorModelsHandler(queries Querier, logger *slog.Logger) *ProjectorModelsHandler {
	return &ProjectorModelsHandler{queries: queries, logger: logger}
}

func (h *ProjectorModelsHandler) List(w http.ResponseWriter, r *http.Request) {
	models, err := h.queries.ListProjectorModels(r.Context())
	if err != nil {
		h.logger.Error("list projector models", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, models)
}

func (h *ProjectorModelsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := h.queries.GetProjectorModel(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "projector model not found")
			return
		}
		h.logger.Error("get projector model", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, m)
}

func (h *ProjectorModelsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	var req struct {
		BrandID   int64  `json:"brand_id"`
		ModelName string `json:"model_name"`
	}
	if err := decodeJSON(r, &req); err != nil || req.BrandID == 0 || req.ModelName == "" {
		respondError(w, http.StatusBadRequest, "brand_id and model_name required")
		return
	}
	m, err := h.queries.CreateProjectorModel(r.Context(), dbsqlc.CreateProjectorModelParams{
		BrandID:   req.BrandID,
		ModelName: req.ModelName,
	})
	if err != nil {
		h.logger.Error("create projector model", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, m)
}

func (h *ProjectorModelsHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		BrandID   *int64  `json:"brand_id"`
		ModelName *string `json:"model_name"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	m, err := h.queries.UpdateProjectorModel(r.Context(), dbsqlc.UpdateProjectorModelParams{
		ProjectorModelID: id,
		BrandID:          toPgInt8(req.BrandID),
		ModelName:        toPgText(req.ModelName),
	})
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "projector model not found")
			return
		}
		h.logger.Error("update projector model", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, m)
}

func (h *ProjectorModelsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.queries.DeleteProjectorModel(r.Context(), id); err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "projector model not found")
			return
		}
		h.logger.Error("delete projector model", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

// ============================================================
// Projectors (unitats)
// ============================================================

type ProjectorsHandler struct {
	queries Querier
	pool    DB
	logger  *slog.Logger
}

func NewProjectorsHandler(queries Querier, pool DB, logger *slog.Logger) *ProjectorsHandler {
	return &ProjectorsHandler{queries: queries, pool: pool, logger: logger}
}

func (h *ProjectorsHandler) List(w http.ResponseWriter, r *http.Request) {
	projectors, err := h.queries.ListProjectors(r.Context())
	if err != nil {
		h.logger.Error("list projectors", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, projectors)
}

func (h *ProjectorsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	p, err := h.queries.GetProjector(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "projector not found")
			return
		}
		h.logger.Error("get projector", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, p)
}

func (h *ProjectorsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	user := currentUser(r)

	var req struct {
		ProjectorModelID int64   `json:"projector_model_id"`
		SerialNumber     *string `json:"serial_number"`
		Status           string  `json:"status"`
		RoomID           *int64  `json:"room_id"`
		EquipmentUserID  *int64  `json:"equipment_user_id"`
		Observations     *string `json:"observations"`
	}
	if err := decodeJSON(r, &req); err != nil || req.ProjectorModelID == 0 {
		respondError(w, http.StatusBadRequest, "projector_model_id required")
		return
	}
	if req.Status == "" {
		req.Status = "actiu"
	}

	p, err := h.queries.CreateProjector(r.Context(), dbsqlc.CreateProjectorParams{
		ProjectorModelID:   req.ProjectorModelID,
		SerialNumber:       toPgText(req.SerialNumber),
		Status:             dbsqlc.DeviceStatusEnum(req.Status),
		RoomID:             toPgInt8(req.RoomID),
		EquipmentUserID:    toPgInt8(req.EquipmentUserID),
		Observations:       toPgText(req.Observations),
		CreatedByAppUserID: user.AppUserID,
	})
	if err != nil {
		h.logger.Error("create projector", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	newJSON, _ := json.Marshal(p)
	if err := writeAudit(r.Context(), h.queries, user, auditEntry{
		TableName: "projector",
		RecordID:  p.ProjectorID,
		EventType: dbsqlc.AuditEventEnumCreated,
		NewValues: newJSON,
	}); err != nil {
		h.logger.Error("write audit", "error", err)
	}

	result, _ := h.queries.GetProjector(r.Context(), p.ProjectorID)
	respondJSON(w, http.StatusCreated, result)
}

func (h *ProjectorsHandler) Update(w http.ResponseWriter, r *http.Request) {
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
		ProjectorModelID *int64  `json:"projector_model_id"`
		SerialNumber     *string `json:"serial_number"`
		Status           *string `json:"status"`
		RoomID           *int64  `json:"room_id"`
		EquipmentUserID  *int64  `json:"equipment_user_id"`
		Observations     *string `json:"observations"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	old, err := h.queries.GetProjector(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "projector not found")
			return
		}
		h.logger.Error("get projector for update", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	var st dbsqlc.NullDeviceStatusEnum
	if req.Status != nil {
		st = dbsqlc.NullDeviceStatusEnum{DeviceStatusEnum: dbsqlc.DeviceStatusEnum(*req.Status), Valid: true}
	}

	p, err := h.queries.UpdateProjector(r.Context(), dbsqlc.UpdateProjectorParams{
		ProjectorID:      id,
		ProjectorModelID: toPgInt8(req.ProjectorModelID),
		SerialNumber:     toPgText(req.SerialNumber),
		Status:           st,
		RoomID:           toPgInt8(req.RoomID),
		EquipmentUserID:  toPgInt8(req.EquipmentUserID),
		Observations:     toPgText(req.Observations),
	})
	if err != nil {
		h.logger.Error("update projector", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	oldJSON, _ := json.Marshal(old)
	newJSON, _ := json.Marshal(p)
	if err := writeAudit(r.Context(), h.queries, user, auditEntry{
		TableName: "projector",
		RecordID:  id,
		EventType: dbsqlc.AuditEventEnumUpdated,
		OldValues: oldJSON,
		NewValues: newJSON,
	}); err != nil {
		h.logger.Error("write audit", "error", err)
	}

	result, _ := h.queries.GetProjector(r.Context(), id)
	respondJSON(w, http.StatusOK, result)
}

func (h *ProjectorsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	user := currentUser(r)

	old, err := h.queries.GetProjector(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "projector not found")
			return
		}
		h.logger.Error("get projector for delete", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := h.queries.DeleteProjector(r.Context(), id); err != nil {
		h.logger.Error("delete projector", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	oldJSON, _ := json.Marshal(old)
	if err := writeAudit(r.Context(), h.queries, user, auditEntry{
		TableName: "projector",
		RecordID:  id,
		EventType: dbsqlc.AuditEventEnumDeleted,
		OldValues: oldJSON,
	}); err != nil {
		h.logger.Error("write audit", "error", err)
	}

	respondJSON(w, http.StatusNoContent, nil)
}
