package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	dbsqlc "inventari/api/internal/db/sqlc"
)

type RoomsHandler struct {
	queries *dbsqlc.Queries
	logger  *slog.Logger
}

func NewRoomsHandler(queries *dbsqlc.Queries, logger *slog.Logger) *RoomsHandler {
	return &RoomsHandler{queries: queries, logger: logger}
}

func (h *RoomsHandler) ListByCenter(w http.ResponseWriter, r *http.Request) {
	centerID, err := strconv.ParseInt(r.PathValue("centerId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid centerId")
		return
	}
	rooms, err := h.queries.ListRoomsByCenter(r.Context(), centerID)
	if err != nil {
		h.logger.Error("list rooms", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, rooms)
}

func (h *RoomsHandler) Create(w http.ResponseWriter, r *http.Request) {
	centerID, err := strconv.ParseInt(r.PathValue("centerId"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid centerId")
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Name == "" {
		respondError(w, http.StatusBadRequest, "name required")
		return
	}
	room, err := h.queries.CreateRoom(r.Context(), dbsqlc.CreateRoomParams{
		CenterID: centerID,
		Name:     req.Name,
	})
	if err != nil {
		h.logger.Error("create room", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusCreated, room)
}

func (h *RoomsHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	room, err := h.queries.UpdateRoom(r.Context(), dbsqlc.UpdateRoomParams{
		RoomID: id,
		Name:   req.Name,
	})
	if err != nil {
		h.logger.Error("update room", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, room)
}

func (h *RoomsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.queries.DeleteRoom(r.Context(), id); err != nil {
		h.logger.Error("delete room", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}
