package server

import (
	"log/slog"
	"net/http"

	dbsqlc "inventari/api/internal/db/sqlc"
	"inventari/api/internal/handler"
	"inventari/api/internal/middleware"
	"inventari/api/internal/session"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(pool *pgxpool.Pool, logger *slog.Logger) http.Handler {
	queries := dbsqlc.New(pool)
	sessions := session.NewStore()

	mux := http.NewServeMux()

	// ─── Handlers ────────────────────────────────────────────────────────────
	authH := handler.NewAuthHandler(queries, sessions, logger)
	usersH := handler.NewUsersHandler(queries, logger)
	rolesH := handler.NewRolesHandler(queries, logger)
	centersH := handler.NewCentersHandler(queries, logger)
	roomsH := handler.NewRoomsHandler(queries, logger)
	cpusH := handler.NewCPUsHandler(queries, logger)
	osH := handler.NewOSHandler(queries, logger)
	equipUsersH := handler.NewEquipmentUsersHandler(queries, logger)
	brandsH := handler.NewBrandsHandler(queries, logger)
	laptopModH := handler.NewLaptopModelsHandler(queries, logger)
	desktopModH := handler.NewDesktopModelsHandler(queries, logger)
	computersH := handler.NewComputersHandler(queries, pool, logger)
	desktopsH := handler.NewDesktopsHandler(queries, pool, logger)
	laptopsH := handler.NewLaptopsHandler(queries, pool, logger)
	cyclesH := handler.NewCyclesHandler(queries, logger)
	classesH := handler.NewClassesHandler(queries, logger)
	studentsH := handler.NewStudentsHandler(queries, logger)
	assignmentsH := handler.NewAssignmentsHandler(queries, logger)
	auditH := handler.NewAuditHandler(queries, logger)
	printerModH := handler.NewPrinterModelsHandler(queries, logger)
	printerSupH := handler.NewPrinterSuppliesHandler(queries, logger)
	printersH := handler.NewPrintersHandler(queries, pool, logger)
	projectorModH := handler.NewProjectorModelsHandler(queries, logger)
	projectorsH := handler.NewProjectorsHandler(queries, pool, logger)

	// ─── Auth (public) ────────────────────────────────────────────────────────
	mux.HandleFunc("POST /auth/login", authH.Login)
	mux.HandleFunc("POST /auth/logout", authH.Logout)
	mux.HandleFunc("GET /auth/me", authH.Me)
	mux.HandleFunc("POST /auth/change-password", authH.ChangePassword)

	// ─── Users ────────────────────────────────────────────────────────────────
	mux.HandleFunc("GET /users", usersH.List)
	mux.HandleFunc("POST /users", usersH.Create)
	mux.HandleFunc("PATCH /users/{id}", usersH.Update)
	mux.HandleFunc("DELETE /users/{id}", usersH.Delete)

	// ─── Roles ────────────────────────────────────────────────────────────────
	mux.HandleFunc("GET /roles", rolesH.List)
	mux.HandleFunc("POST /roles", rolesH.Create)
	mux.HandleFunc("DELETE /roles/{id}", rolesH.Delete)

	// ─── Centers ──────────────────────────────────────────────────────────────
	mux.HandleFunc("GET /centers", centersH.List)
	mux.HandleFunc("POST /centers", centersH.Create)
	mux.HandleFunc("PATCH /centers/{id}", centersH.Update)
	mux.HandleFunc("DELETE /centers/{id}", centersH.Delete)

	// ─── Rooms ────────────────────────────────────────────────────────────────
	mux.HandleFunc("GET /centers/{centerId}/rooms", roomsH.ListByCenter)
	mux.HandleFunc("POST /centers/{centerId}/rooms", roomsH.Create)
	mux.HandleFunc("PATCH /rooms/{id}", roomsH.Update)
	mux.HandleFunc("DELETE /rooms/{id}", roomsH.Delete)

	// ─── CPUs ─────────────────────────────────────────────────────────────────
	mux.HandleFunc("GET /cpus", cpusH.List)
	mux.HandleFunc("POST /cpus", cpusH.Create)
	mux.HandleFunc("PATCH /cpus/{id}", cpusH.Update)
	mux.HandleFunc("DELETE /cpus/{id}", cpusH.Delete)

	// ─── Operating Systems ────────────────────────────────────────────────────
	mux.HandleFunc("GET /os", osH.List)
	mux.HandleFunc("POST /os", osH.Create)
	mux.HandleFunc("PATCH /os/{id}", osH.Update)
	mux.HandleFunc("DELETE /os/{id}", osH.Delete)

	// ─── Equipment Users ──────────────────────────────────────────────────────
	mux.HandleFunc("GET /equipment-users", equipUsersH.List)
	mux.HandleFunc("POST /equipment-users", equipUsersH.Create)
	mux.HandleFunc("PATCH /equipment-users/{id}", equipUsersH.Update)
	mux.HandleFunc("DELETE /equipment-users/{id}", equipUsersH.Delete)

	// ─── Brands ───────────────────────────────────────────────────────────────
	mux.HandleFunc("GET /brands", brandsH.List)
	mux.HandleFunc("POST /brands", brandsH.Create)
	mux.HandleFunc("PATCH /brands/{id}", brandsH.Update)
	mux.HandleFunc("DELETE /brands/{id}", brandsH.Delete)

	// ─── Laptop Models ────────────────────────────────────────────────────────
	mux.HandleFunc("GET /laptop-models", laptopModH.List)
	mux.HandleFunc("POST /laptop-models", laptopModH.Create)
	mux.HandleFunc("GET /laptop-models/{id}", laptopModH.Get)
	mux.HandleFunc("PATCH /laptop-models/{id}", laptopModH.Update)
	mux.HandleFunc("DELETE /laptop-models/{id}", laptopModH.Delete)

	// ─── Desktop Models ───────────────────────────────────────────────────────
	mux.HandleFunc("GET /desktop-models", desktopModH.List)
	mux.HandleFunc("POST /desktop-models", desktopModH.Create)
	mux.HandleFunc("GET /desktop-models/{id}", desktopModH.Get)
	mux.HandleFunc("PATCH /desktop-models/{id}", desktopModH.Update)
	mux.HandleFunc("DELETE /desktop-models/{id}", desktopModH.Delete)

	// ─── Computers (base) ─────────────────────────────────────────────────────
	mux.HandleFunc("GET /computers", computersH.List)
	mux.HandleFunc("GET /computers/{id}", computersH.Get)
	mux.HandleFunc("DELETE /computers/{id}", computersH.Delete)

	// ─── Desktops ─────────────────────────────────────────────────────────────
	mux.HandleFunc("GET /desktops", desktopsH.List)
	mux.HandleFunc("POST /desktops", desktopsH.Create)
	mux.HandleFunc("GET /desktops/{id}", desktopsH.Get)
	mux.HandleFunc("PATCH /desktops/{id}", desktopsH.Update)

	// ─── Laptops ──────────────────────────────────────────────────────────────
	mux.HandleFunc("GET /laptops", laptopsH.List)
	mux.HandleFunc("POST /laptops", laptopsH.Create)
	mux.HandleFunc("GET /laptops/{id}", laptopsH.Get)
	mux.HandleFunc("PATCH /laptops/{id}", laptopsH.Update)

	// ─── Laptop Assignments ───────────────────────────────────────────────────
	mux.HandleFunc("GET /laptops/{laptopId}/assignments", assignmentsH.List)
	mux.HandleFunc("POST /laptops/{laptopId}/assignments", assignmentsH.Create)
	mux.HandleFunc("GET /classes/{classId}/assignments", assignmentsH.ListByClass)
	mux.HandleFunc("GET /assignments", assignmentsH.ListByYear)
	mux.HandleFunc("GET /academic-years", assignmentsH.ListAcademicYears)
	mux.HandleFunc("GET /assignments/{id}", assignmentsH.Get)
	mux.HandleFunc("PATCH /assignments/{id}", assignmentsH.Update)
	mux.HandleFunc("DELETE /assignments/{id}", assignmentsH.Delete)

	// ─── Cycles ───────────────────────────────────────────────────────────────
	mux.HandleFunc("GET /cycles", cyclesH.List)
	mux.HandleFunc("POST /cycles", cyclesH.Create)
	mux.HandleFunc("PATCH /cycles/{id}", cyclesH.Update)
	mux.HandleFunc("DELETE /cycles/{id}", cyclesH.Delete)

	// ─── Classes ──────────────────────────────────────────────────────────────
	mux.HandleFunc("GET /classes", classesH.List)
	mux.HandleFunc("GET /cycles/{cycleId}/classes", classesH.List)
	mux.HandleFunc("POST /cycles/{cycleId}/classes", classesH.Create)
	mux.HandleFunc("GET /classes/{id}", classesH.Get)
	mux.HandleFunc("PATCH /classes/{id}", classesH.Update)
	mux.HandleFunc("DELETE /classes/{id}", classesH.Delete)

	// ─── Tutor ────────────────────────────────────────────────────────────────
	mux.HandleFunc("GET /tutor/classes", classesH.Mine)

	// ─── Students ─────────────────────────────────────────────────────────────
	mux.HandleFunc("GET /classes/{classId}/students", studentsH.List)
	mux.HandleFunc("POST /classes/{classId}/students", studentsH.Create)
	mux.HandleFunc("POST /classes/{classId}/students/import", studentsH.ImportCSV)
	mux.HandleFunc("GET /students/{id}", studentsH.Get)
	mux.HandleFunc("PATCH /students/{id}", studentsH.Update)
	mux.HandleFunc("DELETE /students/{id}", studentsH.Delete)

	// ─── Audit ────────────────────────────────────────────────────────────────
	mux.HandleFunc("GET /audit", auditH.Get)

	// ─── Printer Models ───────────────────────────────────────────────────────
	mux.HandleFunc("GET /printer-models", printerModH.List)
	mux.HandleFunc("POST /printer-models", printerModH.Create)
	mux.HandleFunc("GET /printer-models/{id}", printerModH.Get)
	mux.HandleFunc("PATCH /printer-models/{id}", printerModH.Update)
	mux.HandleFunc("DELETE /printer-models/{id}", printerModH.Delete)

	// Supplies linked to a specific printer model
	mux.HandleFunc("GET /printer-models/{modelId}/supplies", printerSupH.ListByModel)
	mux.HandleFunc("POST /printer-models/{modelId}/supplies", printerSupH.AddToModel)
	mux.HandleFunc("DELETE /printer-models/{modelId}/supplies/{supplyId}", printerSupH.RemoveFromModel)

	// ─── Printer Supplies (consumibles) catalog ───────────────────────────────
	mux.HandleFunc("GET /printer-supplies", printerSupH.List)
	mux.HandleFunc("POST /printer-supplies", printerSupH.Create)
	mux.HandleFunc("PATCH /printer-supplies/{id}", printerSupH.Update)
	mux.HandleFunc("DELETE /printer-supplies/{id}", printerSupH.Delete)

	// ─── Printers ─────────────────────────────────────────────────────────────
	mux.HandleFunc("GET /printers", printersH.List)
	mux.HandleFunc("POST /printers", printersH.Create)
	mux.HandleFunc("GET /printers/{id}", printersH.Get)
	mux.HandleFunc("PATCH /printers/{id}", printersH.Update)
	mux.HandleFunc("DELETE /printers/{id}", printersH.Delete)

	// ─── Projectors ───────────────────────────────────────────────────────────
	mux.HandleFunc("GET /projectors", projectorsH.List)
	mux.HandleFunc("POST /projectors", projectorsH.Create)
	mux.HandleFunc("GET /projectors/{id}", projectorsH.Get)
	mux.HandleFunc("PATCH /projectors/{id}", projectorsH.Update)
	mux.HandleFunc("DELETE /projectors/{id}", projectorsH.Delete)

	// ─── Projector Models ─────────────────────────────────────────────────────
	mux.HandleFunc("GET /projector-models", projectorModH.List)
	mux.HandleFunc("POST /projector-models", projectorModH.Create)
	mux.HandleFunc("GET /projector-models/{id}", projectorModH.Get)
	mux.HandleFunc("PATCH /projector-models/{id}", projectorModH.Update)
	mux.HandleFunc("DELETE /projector-models/{id}", projectorModH.Delete)

	return middleware.CORS(
		middleware.Logger(logger)(
			middleware.Auth(sessions)(mux),
		),
	)
}
