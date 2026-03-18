package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"
)

type BrandsHandler struct {
	queries Querier
	logger  *slog.Logger
}

func NewBrandsHandler(queries Querier, logger *slog.Logger) *BrandsHandler {
	return &BrandsHandler{queries: queries, logger: logger}
}

func (h *BrandsHandler) List(w http.ResponseWriter, r *http.Request) {
	brands, err := h.queries.ListBrands(r.Context())
	if err != nil {
		h.logger.Error("list brands", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, brands)
}

func (h *BrandsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Name == "" {
		respondError(w, http.StatusBadRequest, "name required")
		return
	}
	brand, err := h.queries.CreateBrand(r.Context(), req.Name)
	if err != nil {
		h.logger.Error("create brand", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, brand)
}

func (h *BrandsHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Name == "" {
		respondError(w, http.StatusBadRequest, "name required")
		return
	}
	brand, err := h.queries.UpdateBrand(r.Context(), dbsqlc.UpdateBrandParams{
		BrandID: id,
		Name:    req.Name,
	})
	if err != nil {
		if isNotFound(err) {
			respondError(w, http.StatusNotFound, "brand not found")
			return
		}
		h.logger.Error("update brand", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, brand)
}

func (h *BrandsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, RoleAdmin, RoleEditor) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.queries.DeleteBrand(r.Context(), id); err != nil {
		h.logger.Error("delete brand", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
