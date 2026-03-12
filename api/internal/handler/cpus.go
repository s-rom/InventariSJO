package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"
)

type CPUsHandler struct {
	queries *dbsqlc.Queries
	logger  *slog.Logger
}

func NewCPUsHandler(queries *dbsqlc.Queries, logger *slog.Logger) *CPUsHandler {
	return &CPUsHandler{queries: queries, logger: logger}
}

func (h *CPUsHandler) List(w http.ResponseWriter, r *http.Request) {
	cpus, err := h.queries.ListCPUs(r.Context())
	if err != nil {
		h.logger.Error("list cpus", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, cpus)
}

func (h *CPUsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ModelName      string `json:"model_name"`
		BenchmarkScore *int32 `json:"benchmark_score"`
	}
	if err := decodeJSON(r, &req); err != nil || req.ModelName == "" {
		respondError(w, http.StatusBadRequest, "model_name required")
		return
	}
	cpu, err := h.queries.CreateCPU(r.Context(), dbsqlc.CreateCPUParams{
		ModelName:      req.ModelName,
		BenchmarkScore: toPgInt4(req.BenchmarkScore),
	})
	if err != nil {
		h.logger.Error("create cpu", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, cpu)
}

func (h *CPUsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		ModelName      *string `json:"model_name"`
		BenchmarkScore *int32  `json:"benchmark_score"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	cpu, err := h.queries.UpdateCPU(r.Context(), dbsqlc.UpdateCPUParams{
		CpuID:          id,
		ModelName:      toPgText(req.ModelName),
		BenchmarkScore: toPgInt4(req.BenchmarkScore),
	})
	if err != nil {
		h.logger.Error("update cpu", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, cpu)
}

func (h *CPUsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.queries.DeleteCPU(r.Context(), id); err != nil {
		h.logger.Error("delete cpu", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
