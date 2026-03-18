package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"
)

// ── Laptop Models ─────────────────────────────────────────────────────────────

type LaptopModelsHandler struct {
	queries Querier
	logger  *slog.Logger
}

func NewLaptopModelsHandler(queries Querier, logger *slog.Logger) *LaptopModelsHandler {
	return &LaptopModelsHandler{queries: queries, logger: logger}
}

func (h *LaptopModelsHandler) List(w http.ResponseWriter, r *http.Request) {
	models, err := h.queries.ListLaptopModels(r.Context())
	if err != nil {
		h.logger.Error("list laptop models", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, models)
}

func (h *LaptopModelsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := h.queries.GetLaptopModel(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "laptop model not found")
			return
		}
		h.logger.Error("get laptop model", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, m)
}

func (h *LaptopModelsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	var req struct {
		BrandID      int64  `json:"brand_id"`
		ModelName    string `json:"model_name"`
		CpuID        *int64 `json:"cpu_id"`
		BaseRamGb    int32  `json:"base_ram_gb"`
		BaseRamType  string `json:"base_ram_type"`
		BaseStorGb   int32  `json:"base_storage_gb"`
		BaseStorType string `json:"base_storage_type"`
		BaseOsID     *int64 `json:"base_os_id"`
	}
	if err := decodeJSON(r, &req); err != nil || req.ModelName == "" || req.BrandID == 0 {
		respondError(w, http.StatusBadRequest, "brand_id and model_name required")
		return
	}
	m, err := h.queries.CreateLaptopModel(r.Context(), dbsqlc.CreateLaptopModelParams{
		BrandID:         req.BrandID,
		ModelName:       req.ModelName,
		CpuID:           toPgInt8(req.CpuID),
		BaseRamGb:       req.BaseRamGb,
		BaseRamType:     dbsqlc.RamTypeEnum(req.BaseRamType),
		BaseStorageGb:   req.BaseStorGb,
		BaseStorageType: dbsqlc.StorageTypeEnum(req.BaseStorType),
		BaseOsID:        toPgInt8(req.BaseOsID),
	})
	if err != nil {
		h.logger.Error("create laptop model", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, m)
}

func (h *LaptopModelsHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		BrandID      *int64  `json:"brand_id"`
		ModelName    *string `json:"model_name"`
		CpuID        *int64  `json:"cpu_id"`
		BaseRamGb    *int32  `json:"base_ram_gb"`
		BaseRamType  *string `json:"base_ram_type"`
		BaseStorGb   *int32  `json:"base_storage_gb"`
		BaseStorType *string `json:"base_storage_type"`
		BaseOsID     *int64  `json:"base_os_id"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	var ramType dbsqlc.NullRamTypeEnum
	if req.BaseRamType != nil {
		ramType = dbsqlc.NullRamTypeEnum{RamTypeEnum: dbsqlc.RamTypeEnum(*req.BaseRamType), Valid: true}
	}
	var storType dbsqlc.NullStorageTypeEnum
	if req.BaseStorType != nil {
		storType = dbsqlc.NullStorageTypeEnum{StorageTypeEnum: dbsqlc.StorageTypeEnum(*req.BaseStorType), Valid: true}
	}
	m, err := h.queries.UpdateLaptopModel(r.Context(), dbsqlc.UpdateLaptopModelParams{
		LaptopModelID:   id,
		BrandID:         toPgInt8(req.BrandID),
		ModelName:       toPgText(req.ModelName),
		CpuID:           toPgInt8(req.CpuID),
		BaseRamGb:       toPgInt4(req.BaseRamGb),
		BaseRamType:     ramType,
		BaseStorageGb:   toPgInt4(req.BaseStorGb),
		BaseStorageType: storType,
		BaseOsID:        toPgInt8(req.BaseOsID),
	})
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "laptop model not found")
			return
		}
		h.logger.Error("update laptop model", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, m)
}

func (h *LaptopModelsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.queries.DeleteLaptopModel(r.Context(), id); err != nil {
		h.logger.Error("delete laptop model", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── Desktop Models ────────────────────────────────────────────────────────────

type DesktopModelsHandler struct {
	queries Querier
	logger  *slog.Logger
}

func NewDesktopModelsHandler(queries Querier, logger *slog.Logger) *DesktopModelsHandler {
	return &DesktopModelsHandler{queries: queries, logger: logger}
}

func (h *DesktopModelsHandler) List(w http.ResponseWriter, r *http.Request) {
	models, err := h.queries.ListDesktopModels(r.Context())
	if err != nil {
		h.logger.Error("list desktop models", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, models)
}

func (h *DesktopModelsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := h.queries.GetDesktopModel(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "desktop model not found")
			return
		}
		h.logger.Error("get desktop model", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, m)
}

func (h *DesktopModelsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	var req struct {
		BrandID      int64  `json:"brand_id"`
		ModelName    string `json:"model_name"`
		CpuID        *int64 `json:"cpu_id"`
		BaseRamGb    int32  `json:"base_ram_gb"`
		BaseRamType  string `json:"base_ram_type"`
		BaseStorGb   int32  `json:"base_storage_gb"`
		BaseStorType string `json:"base_storage_type"`
		BaseOsID     *int64 `json:"base_os_id"`
	}
	if err := decodeJSON(r, &req); err != nil || req.ModelName == "" || req.BrandID == 0 {
		respondError(w, http.StatusBadRequest, "brand_id and model_name required")
		return
	}
	m, err := h.queries.CreateDesktopModel(r.Context(), dbsqlc.CreateDesktopModelParams{
		BrandID:         req.BrandID,
		ModelName:       req.ModelName,
		CpuID:           toPgInt8(req.CpuID),
		BaseRamGb:       req.BaseRamGb,
		BaseRamType:     dbsqlc.RamTypeEnum(req.BaseRamType),
		BaseStorageGb:   req.BaseStorGb,
		BaseStorageType: dbsqlc.StorageTypeEnum(req.BaseStorType),
		BaseOsID:        toPgInt8(req.BaseOsID),
	})
	if err != nil {
		h.logger.Error("create desktop model", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, m)
}

func (h *DesktopModelsHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		BrandID      *int64  `json:"brand_id"`
		ModelName    *string `json:"model_name"`
		CpuID        *int64  `json:"cpu_id"`
		BaseRamGb    *int32  `json:"base_ram_gb"`
		BaseRamType  *string `json:"base_ram_type"`
		BaseStorGb   *int32  `json:"base_storage_gb"`
		BaseStorType *string `json:"base_storage_type"`
		BaseOsID     *int64  `json:"base_os_id"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	var ramType dbsqlc.NullRamTypeEnum
	if req.BaseRamType != nil {
		ramType = dbsqlc.NullRamTypeEnum{RamTypeEnum: dbsqlc.RamTypeEnum(*req.BaseRamType), Valid: true}
	}
	var storType dbsqlc.NullStorageTypeEnum
	if req.BaseStorType != nil {
		storType = dbsqlc.NullStorageTypeEnum{StorageTypeEnum: dbsqlc.StorageTypeEnum(*req.BaseStorType), Valid: true}
	}
	m, err := h.queries.UpdateDesktopModel(r.Context(), dbsqlc.UpdateDesktopModelParams{
		DesktopModelID:  id,
		BrandID:         toPgInt8(req.BrandID),
		ModelName:       toPgText(req.ModelName),
		CpuID:           toPgInt8(req.CpuID),
		BaseRamGb:       toPgInt4(req.BaseRamGb),
		BaseRamType:     ramType,
		BaseStorageGb:   toPgInt4(req.BaseStorGb),
		BaseStorageType: storType,
		BaseOsID:        toPgInt8(req.BaseOsID),
	})
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "desktop model not found")
			return
		}
		h.logger.Error("update desktop model", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, m)
}

func (h *DesktopModelsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.queries.DeleteDesktopModel(r.Context(), id); err != nil {
		h.logger.Error("delete desktop model", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
