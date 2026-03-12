package server

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	dbsqlc "inventari/api/internal/db/sqlc"
	"inventari/api/internal/handler"
	"inventari/api/internal/middleware"
	"inventari/api/internal/session"
)

func New(pool *pgxpool.Pool, logger *slog.Logger) http.Handler {
	queries := dbsqlc.New(pool)
	sessions := session.NewStore()

	mux := http.NewServeMux()

	authH := handler.NewAuthHandler(queries, sessions, logger)
	usersH := handler.NewUsersHandler(queries, logger)
	centersH := handler.NewCentersHandler(queries, logger)
	roomsH := handler.NewRoomsHandler(queries, logger)
	cpusH := handler.NewCPUsHandler(queries, logger)
	osH := handler.NewOSHandler(queries, logger)
	equipUsersH := handler.NewEquipmentUsersHandler(queries, logger)
	computersH := handler.NewComputersHandler(queries, pool, logger)

	// Auth (public)
	mux.HandleFunc("POST /auth/login", authH.Login)
	mux.HandleFunc("POST /auth/logout", authH.Logout)

	// Users
	mux.HandleFunc("GET /users", usersH.List)
	mux.HandleFunc("POST /users", usersH.Create)
	mux.HandleFunc("PATCH /users/{id}", usersH.Update)
	mux.HandleFunc("DELETE /users/{id}", usersH.Delete)

	// Centers
	mux.HandleFunc("GET /centers", centersH.List)
	mux.HandleFunc("POST /centers", centersH.Create)
	mux.HandleFunc("PATCH /centers/{id}", centersH.Update)
	mux.HandleFunc("DELETE /centers/{id}", centersH.Delete)

	// Rooms nested under centers + standalone
	mux.HandleFunc("GET /centers/{centerId}/rooms", roomsH.ListByCenter)
	mux.HandleFunc("POST /centers/{centerId}/rooms", roomsH.Create)
	mux.HandleFunc("PATCH /rooms/{id}", roomsH.Update)
	mux.HandleFunc("DELETE /rooms/{id}", roomsH.Delete)

	// CPUs
	mux.HandleFunc("GET /cpus", cpusH.List)
	mux.HandleFunc("POST /cpus", cpusH.Create)
	mux.HandleFunc("PATCH /cpus/{id}", cpusH.Update)
	mux.HandleFunc("DELETE /cpus/{id}", cpusH.Delete)

	// Operating Systems
	mux.HandleFunc("GET /os", osH.List)
	mux.HandleFunc("POST /os", osH.Create)
	mux.HandleFunc("PATCH /os/{id}", osH.Update)
	mux.HandleFunc("DELETE /os/{id}", osH.Delete)

	// Equipment Users
	mux.HandleFunc("GET /equipment-users", equipUsersH.List)
	mux.HandleFunc("POST /equipment-users", equipUsersH.Create)
	mux.HandleFunc("PATCH /equipment-users/{id}", equipUsersH.Update)
	mux.HandleFunc("DELETE /equipment-users/{id}", equipUsersH.Delete)

	// Computers
	mux.HandleFunc("GET /computers", computersH.List)
	mux.HandleFunc("GET /computers/{id}", computersH.Get)
	mux.HandleFunc("POST /computers", computersH.Create)
	mux.HandleFunc("PATCH /computers/{id}", computersH.Update)
	mux.HandleFunc("DELETE /computers/{id}", computersH.Delete)
	mux.HandleFunc("GET /computers/{id}/audit", computersH.Audit)

	return middleware.CORS(
		middleware.Logger(logger)(
			middleware.Auth(sessions)(mux),
		),
	)
}
