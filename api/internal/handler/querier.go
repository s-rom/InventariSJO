package handler

import (
	"context"

	dbsqlc "inventari/api/internal/db/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Querier is satisfied by *dbsqlc.Queries and can be mocked in unit tests.
type Querier interface {
	// Assignments
	CreateAssignment(ctx context.Context, arg dbsqlc.CreateAssignmentParams) (dbsqlc.LaptopStudentAssignment, error)
	DeleteAssignment(ctx context.Context, assignmentID int64) error
	GetAssignment(ctx context.Context, assignmentID int64) (dbsqlc.LaptopStudentAssignment, error)
	ListAssignmentsByClass(ctx context.Context, classID int64) ([]dbsqlc.ListAssignmentsByClassRow, error)
	ListAssignmentsByClassAndYear(ctx context.Context, arg dbsqlc.ListAssignmentsByClassAndYearParams) ([]dbsqlc.ListAssignmentsByClassAndYearRow, error)
	ListAssignmentsByLaptop(ctx context.Context, computerID int64) ([]dbsqlc.ListAssignmentsByLaptopRow, error)
	ListAssignmentsByYear(ctx context.Context, academicYear string) ([]dbsqlc.ListAssignmentsByYearRow, error)
	ListDistinctAcademicYears(ctx context.Context) ([]string, error)
	UpdateAssignment(ctx context.Context, arg dbsqlc.UpdateAssignmentParams) (dbsqlc.LaptopStudentAssignment, error)
	// Audit
	GetAuditLog(ctx context.Context, arg dbsqlc.GetAuditLogParams) ([]dbsqlc.AuditLog, error)
	InsertAuditLog(ctx context.Context, arg dbsqlc.InsertAuditLogParams) error
	// Auth
	GetUserByUsername(ctx context.Context, username string) (dbsqlc.AppUser, error)
	UpdateUserPassword(ctx context.Context, arg dbsqlc.UpdateUserPasswordParams) error
	// Brands
	CreateBrand(ctx context.Context, name string) (dbsqlc.Brand, error)
	DeleteBrand(ctx context.Context, brandID int64) error
	GetBrand(ctx context.Context, brandID int64) (dbsqlc.Brand, error)
	ListBrands(ctx context.Context) ([]dbsqlc.Brand, error)
	UpdateBrand(ctx context.Context, arg dbsqlc.UpdateBrandParams) (dbsqlc.Brand, error)
	// Centers
	CreateCenter(ctx context.Context, name string) (dbsqlc.Center, error)
	DeleteCenter(ctx context.Context, centerID int64) error
	ListCenters(ctx context.Context) ([]dbsqlc.Center, error)
	UpdateCenter(ctx context.Context, arg dbsqlc.UpdateCenterParams) (dbsqlc.Center, error)
	// Computers
	CreateComputer(ctx context.Context, arg dbsqlc.CreateComputerParams) (dbsqlc.Computer, error)
	DeleteComputer(ctx context.Context, computerID int64) error
	GetComputerBase(ctx context.Context, computerID int64) (dbsqlc.Computer, error)
	ListComputers(ctx context.Context) ([]dbsqlc.ListComputersRow, error)
	UpdateComputerBase(ctx context.Context, arg dbsqlc.UpdateComputerBaseParams) (dbsqlc.Computer, error)
	// CPUs
	CreateCPU(ctx context.Context, arg dbsqlc.CreateCPUParams) (dbsqlc.Cpu, error)
	DeleteCPU(ctx context.Context, cpuID int64) error
	ListCPUs(ctx context.Context) ([]dbsqlc.Cpu, error)
	UpdateCPU(ctx context.Context, arg dbsqlc.UpdateCPUParams) (dbsqlc.Cpu, error)
	// Cycles
	CreateCycle(ctx context.Context, name string) (dbsqlc.Cycle, error)
	DeleteCycle(ctx context.Context, cycleID int64) error
	ListCycles(ctx context.Context) ([]dbsqlc.Cycle, error)
	UpdateCycle(ctx context.Context, arg dbsqlc.UpdateCycleParams) (dbsqlc.Cycle, error)
	// Desktop Models
	CreateDesktopModel(ctx context.Context, arg dbsqlc.CreateDesktopModelParams) (dbsqlc.DesktopModel, error)
	DeleteDesktopModel(ctx context.Context, desktopModelID int64) error
	GetDesktopModel(ctx context.Context, desktopModelID int64) (dbsqlc.GetDesktopModelRow, error)
	ListDesktopModels(ctx context.Context) ([]dbsqlc.ListDesktopModelsRow, error)
	UpdateDesktopModel(ctx context.Context, arg dbsqlc.UpdateDesktopModelParams) (dbsqlc.DesktopModel, error)
	// Desktops
	CreateDesktop(ctx context.Context, arg dbsqlc.CreateDesktopParams) (dbsqlc.Desktop, error)
	GetDesktop(ctx context.Context, computerID int64) (dbsqlc.GetDesktopRow, error)
	ListDesktops(ctx context.Context) ([]dbsqlc.ListDesktopsRow, error)
	UpdateDesktop(ctx context.Context, arg dbsqlc.UpdateDesktopParams) (dbsqlc.Desktop, error)
	// Equipment Users
	CreateEquipmentUser(ctx context.Context, name string) (dbsqlc.EquipmentUser, error)
	DeleteEquipmentUser(ctx context.Context, equipmentUserID int64) error
	ListEquipmentUsers(ctx context.Context) ([]dbsqlc.EquipmentUser, error)
	UpdateEquipmentUser(ctx context.Context, arg dbsqlc.UpdateEquipmentUserParams) (dbsqlc.EquipmentUser, error)
	// Laptop Models
	CreateLaptopModel(ctx context.Context, arg dbsqlc.CreateLaptopModelParams) (dbsqlc.LaptopModel, error)
	DeleteLaptopModel(ctx context.Context, laptopModelID int64) error
	GetLaptopModel(ctx context.Context, laptopModelID int64) (dbsqlc.GetLaptopModelRow, error)
	ListLaptopModels(ctx context.Context) ([]dbsqlc.ListLaptopModelsRow, error)
	UpdateLaptopModel(ctx context.Context, arg dbsqlc.UpdateLaptopModelParams) (dbsqlc.LaptopModel, error)
	// Laptops
	CreateLaptop(ctx context.Context, arg dbsqlc.CreateLaptopParams) (dbsqlc.Laptop, error)
	GetLaptop(ctx context.Context, computerID int64) (dbsqlc.GetLaptopRow, error)
	ListLaptops(ctx context.Context) ([]dbsqlc.ListLaptopsRow, error)
	UpdateLaptop(ctx context.Context, arg dbsqlc.UpdateLaptopParams) (dbsqlc.Laptop, error)
	// Operating Systems
	CreateOS(ctx context.Context, name string) (dbsqlc.O, error)
	DeleteOS(ctx context.Context, osID int64) error
	ListOS(ctx context.Context) ([]dbsqlc.O, error)
	UpdateOS(ctx context.Context, arg dbsqlc.UpdateOSParams) (dbsqlc.O, error)
	// Roles
	CreateRole(ctx context.Context, arg dbsqlc.CreateRoleParams) (dbsqlc.Role, error)
	DeleteRole(ctx context.Context, roleID string) error
	ListRoles(ctx context.Context) ([]dbsqlc.Role, error)
	// Rooms
	CreateRoom(ctx context.Context, arg dbsqlc.CreateRoomParams) (dbsqlc.Room, error)
	DeleteRoom(ctx context.Context, roomID int64) error
	ListRoomsByCenter(ctx context.Context, centerID int64) ([]dbsqlc.Room, error)
	UpdateRoom(ctx context.Context, arg dbsqlc.UpdateRoomParams) (dbsqlc.Room, error)
	// School Classes
	CreateClass(ctx context.Context, arg dbsqlc.CreateClassParams) (dbsqlc.SchoolClass, error)
	DeleteClass(ctx context.Context, classID int64) error
	GetClass(ctx context.Context, classID int64) (dbsqlc.SchoolClass, error)
	ListClasses(ctx context.Context) ([]dbsqlc.ListClassesRow, error)
	ListClassesByCycle(ctx context.Context, cycleID int64) ([]dbsqlc.ListClassesByCycleRow, error)
	ListClassesByTutor(ctx context.Context, tutorAppUserID pgtype.Int8) ([]dbsqlc.ListClassesByTutorRow, error)
	UpdateClass(ctx context.Context, arg dbsqlc.UpdateClassParams) (dbsqlc.SchoolClass, error)
	// Students
	CreateStudent(ctx context.Context, arg dbsqlc.CreateStudentParams) (dbsqlc.Student, error)
	DeleteStudent(ctx context.Context, studentID int64) error
	GetStudent(ctx context.Context, studentID int64) (dbsqlc.Student, error)
	ListStudentsByClass(ctx context.Context, classID int64) ([]dbsqlc.Student, error)
	UpdateStudent(ctx context.Context, arg dbsqlc.UpdateStudentParams) (dbsqlc.Student, error)
	// Users
	CreateUser(ctx context.Context, arg dbsqlc.CreateUserParams) (dbsqlc.CreateUserRow, error)
	DeleteUser(ctx context.Context, appUserID int64) error
	ListUsers(ctx context.Context) ([]dbsqlc.ListUsersRow, error)
	UpdateUser(ctx context.Context, arg dbsqlc.UpdateUserParams) (dbsqlc.UpdateUserRow, error)
}

// DB abstracts the connection pool for handlers that need transactions.
// *pgxpool.Pool satisfies this interface.
type DB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// compile-time assertions
var (
	_ Querier = (*dbsqlc.Queries)(nil)
)
