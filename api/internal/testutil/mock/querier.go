// Package testmock provides a testify/mock implementation of handler.Querier
// for use in unit tests across any package.
package testmock

import (
	"context"

	dbsqlc "inventari/api/internal/db/sqlc"
	"inventari/api/internal/handler"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
)

// compile-time check
var _ handler.Querier = (*Querier)(nil)

// Querier is a testify/mock implementation of handler.Querier.
type Querier struct{ mock.Mock }

func (m *Querier) CreateAssignment(ctx context.Context, arg dbsqlc.CreateAssignmentParams) (dbsqlc.LaptopStudentAssignment, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.LaptopStudentAssignment)
	return v, a.Error(1)
}
func (m *Querier) DeleteAssignment(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *Querier) GetAssignment(ctx context.Context, id int64) (dbsqlc.LaptopStudentAssignment, error) {
	a := m.Called(ctx, id)
	v, _ := a.Get(0).(dbsqlc.LaptopStudentAssignment)
	return v, a.Error(1)
}
func (m *Querier) ListAssignmentsByClass(ctx context.Context, classID int64) ([]dbsqlc.ListAssignmentsByClassRow, error) {
	a := m.Called(ctx, classID)
	v, _ := a.Get(0).([]dbsqlc.ListAssignmentsByClassRow)
	return v, a.Error(1)
}
func (m *Querier) ListAssignmentsByClassAndYear(ctx context.Context, arg dbsqlc.ListAssignmentsByClassAndYearParams) ([]dbsqlc.ListAssignmentsByClassAndYearRow, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).([]dbsqlc.ListAssignmentsByClassAndYearRow)
	return v, a.Error(1)
}
func (m *Querier) ListAssignmentsByLaptop(ctx context.Context, computerID int64) ([]dbsqlc.ListAssignmentsByLaptopRow, error) {
	a := m.Called(ctx, computerID)
	v, _ := a.Get(0).([]dbsqlc.ListAssignmentsByLaptopRow)
	return v, a.Error(1)
}
func (m *Querier) ListAssignmentsByYear(ctx context.Context, academicYear string) ([]dbsqlc.ListAssignmentsByYearRow, error) {
	a := m.Called(ctx, academicYear)
	v, _ := a.Get(0).([]dbsqlc.ListAssignmentsByYearRow)
	return v, a.Error(1)
}
func (m *Querier) ListDistinctAcademicYears(ctx context.Context) ([]string, error) {
	a := m.Called(ctx)
	v, _ := a.Get(0).([]string)
	return v, a.Error(1)
}
func (m *Querier) UpdateAssignment(ctx context.Context, arg dbsqlc.UpdateAssignmentParams) (dbsqlc.LaptopStudentAssignment, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.LaptopStudentAssignment)
	return v, a.Error(1)
}
func (m *Querier) GetAuditLog(ctx context.Context, arg dbsqlc.GetAuditLogParams) ([]dbsqlc.AuditLog, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).([]dbsqlc.AuditLog)
	return v, a.Error(1)
}
func (m *Querier) InsertAuditLog(ctx context.Context, arg dbsqlc.InsertAuditLogParams) error {
	return m.Called(ctx, arg).Error(0)
}
func (m *Querier) GetUserByUsername(ctx context.Context, username string) (dbsqlc.AppUser, error) {
	a := m.Called(ctx, username)
	v, _ := a.Get(0).(dbsqlc.AppUser)
	return v, a.Error(1)
}
func (m *Querier) UpdateUserPassword(ctx context.Context, arg dbsqlc.UpdateUserPasswordParams) error {
	return m.Called(ctx, arg).Error(0)
}
func (m *Querier) CreateBrand(ctx context.Context, name string) (dbsqlc.Brand, error) {
	a := m.Called(ctx, name)
	v, _ := a.Get(0).(dbsqlc.Brand)
	return v, a.Error(1)
}
func (m *Querier) DeleteBrand(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *Querier) GetBrand(ctx context.Context, id int64) (dbsqlc.Brand, error) {
	a := m.Called(ctx, id)
	v, _ := a.Get(0).(dbsqlc.Brand)
	return v, a.Error(1)
}
func (m *Querier) ListBrands(ctx context.Context) ([]dbsqlc.Brand, error) {
	a := m.Called(ctx)
	v, _ := a.Get(0).([]dbsqlc.Brand)
	return v, a.Error(1)
}
func (m *Querier) UpdateBrand(ctx context.Context, arg dbsqlc.UpdateBrandParams) (dbsqlc.Brand, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Brand)
	return v, a.Error(1)
}
func (m *Querier) CreateCenter(ctx context.Context, name string) (dbsqlc.Center, error) {
	a := m.Called(ctx, name)
	v, _ := a.Get(0).(dbsqlc.Center)
	return v, a.Error(1)
}
func (m *Querier) DeleteCenter(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *Querier) ListCenters(ctx context.Context) ([]dbsqlc.Center, error) {
	a := m.Called(ctx)
	v, _ := a.Get(0).([]dbsqlc.Center)
	return v, a.Error(1)
}
func (m *Querier) UpdateCenter(ctx context.Context, arg dbsqlc.UpdateCenterParams) (dbsqlc.Center, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Center)
	return v, a.Error(1)
}
func (m *Querier) CreateComputer(ctx context.Context, arg dbsqlc.CreateComputerParams) (dbsqlc.Computer, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Computer)
	return v, a.Error(1)
}
func (m *Querier) DeleteComputer(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *Querier) GetComputerBase(ctx context.Context, id int64) (dbsqlc.Computer, error) {
	a := m.Called(ctx, id)
	v, _ := a.Get(0).(dbsqlc.Computer)
	return v, a.Error(1)
}
func (m *Querier) ListComputers(ctx context.Context) ([]dbsqlc.ListComputersRow, error) {
	a := m.Called(ctx)
	v, _ := a.Get(0).([]dbsqlc.ListComputersRow)
	return v, a.Error(1)
}
func (m *Querier) UpdateComputerBase(ctx context.Context, arg dbsqlc.UpdateComputerBaseParams) (dbsqlc.Computer, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Computer)
	return v, a.Error(1)
}
func (m *Querier) CreateCPU(ctx context.Context, arg dbsqlc.CreateCPUParams) (dbsqlc.Cpu, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Cpu)
	return v, a.Error(1)
}
func (m *Querier) DeleteCPU(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *Querier) ListCPUs(ctx context.Context) ([]dbsqlc.Cpu, error) {
	a := m.Called(ctx)
	v, _ := a.Get(0).([]dbsqlc.Cpu)
	return v, a.Error(1)
}
func (m *Querier) UpdateCPU(ctx context.Context, arg dbsqlc.UpdateCPUParams) (dbsqlc.Cpu, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Cpu)
	return v, a.Error(1)
}
func (m *Querier) CreateCycle(ctx context.Context, name string) (dbsqlc.Cycle, error) {
	a := m.Called(ctx, name)
	v, _ := a.Get(0).(dbsqlc.Cycle)
	return v, a.Error(1)
}
func (m *Querier) DeleteCycle(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *Querier) ListCycles(ctx context.Context) ([]dbsqlc.Cycle, error) {
	a := m.Called(ctx)
	v, _ := a.Get(0).([]dbsqlc.Cycle)
	return v, a.Error(1)
}
func (m *Querier) UpdateCycle(ctx context.Context, arg dbsqlc.UpdateCycleParams) (dbsqlc.Cycle, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Cycle)
	return v, a.Error(1)
}
func (m *Querier) CreateDesktopModel(ctx context.Context, arg dbsqlc.CreateDesktopModelParams) (dbsqlc.DesktopModel, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.DesktopModel)
	return v, a.Error(1)
}
func (m *Querier) DeleteDesktopModel(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *Querier) GetDesktopModel(ctx context.Context, id int64) (dbsqlc.GetDesktopModelRow, error) {
	a := m.Called(ctx, id)
	v, _ := a.Get(0).(dbsqlc.GetDesktopModelRow)
	return v, a.Error(1)
}
func (m *Querier) ListDesktopModels(ctx context.Context) ([]dbsqlc.ListDesktopModelsRow, error) {
	a := m.Called(ctx)
	v, _ := a.Get(0).([]dbsqlc.ListDesktopModelsRow)
	return v, a.Error(1)
}
func (m *Querier) UpdateDesktopModel(ctx context.Context, arg dbsqlc.UpdateDesktopModelParams) (dbsqlc.DesktopModel, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.DesktopModel)
	return v, a.Error(1)
}
func (m *Querier) CreateDesktop(ctx context.Context, arg dbsqlc.CreateDesktopParams) (dbsqlc.Desktop, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Desktop)
	return v, a.Error(1)
}
func (m *Querier) GetDesktop(ctx context.Context, id int64) (dbsqlc.GetDesktopRow, error) {
	a := m.Called(ctx, id)
	v, _ := a.Get(0).(dbsqlc.GetDesktopRow)
	return v, a.Error(1)
}
func (m *Querier) ListDesktops(ctx context.Context) ([]dbsqlc.ListDesktopsRow, error) {
	a := m.Called(ctx)
	v, _ := a.Get(0).([]dbsqlc.ListDesktopsRow)
	return v, a.Error(1)
}
func (m *Querier) UpdateDesktop(ctx context.Context, arg dbsqlc.UpdateDesktopParams) (dbsqlc.Desktop, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Desktop)
	return v, a.Error(1)
}
func (m *Querier) CreateEquipmentUser(ctx context.Context, name string) (dbsqlc.EquipmentUser, error) {
	a := m.Called(ctx, name)
	v, _ := a.Get(0).(dbsqlc.EquipmentUser)
	return v, a.Error(1)
}
func (m *Querier) DeleteEquipmentUser(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *Querier) ListEquipmentUsers(ctx context.Context) ([]dbsqlc.EquipmentUser, error) {
	a := m.Called(ctx)
	v, _ := a.Get(0).([]dbsqlc.EquipmentUser)
	return v, a.Error(1)
}
func (m *Querier) UpdateEquipmentUser(ctx context.Context, arg dbsqlc.UpdateEquipmentUserParams) (dbsqlc.EquipmentUser, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.EquipmentUser)
	return v, a.Error(1)
}
func (m *Querier) CreateLaptopModel(ctx context.Context, arg dbsqlc.CreateLaptopModelParams) (dbsqlc.LaptopModel, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.LaptopModel)
	return v, a.Error(1)
}
func (m *Querier) DeleteLaptopModel(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *Querier) GetLaptopModel(ctx context.Context, id int64) (dbsqlc.GetLaptopModelRow, error) {
	a := m.Called(ctx, id)
	v, _ := a.Get(0).(dbsqlc.GetLaptopModelRow)
	return v, a.Error(1)
}
func (m *Querier) ListLaptopModels(ctx context.Context) ([]dbsqlc.ListLaptopModelsRow, error) {
	a := m.Called(ctx)
	v, _ := a.Get(0).([]dbsqlc.ListLaptopModelsRow)
	return v, a.Error(1)
}
func (m *Querier) UpdateLaptopModel(ctx context.Context, arg dbsqlc.UpdateLaptopModelParams) (dbsqlc.LaptopModel, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.LaptopModel)
	return v, a.Error(1)
}
func (m *Querier) CreateLaptop(ctx context.Context, arg dbsqlc.CreateLaptopParams) (dbsqlc.Laptop, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Laptop)
	return v, a.Error(1)
}
func (m *Querier) GetLaptop(ctx context.Context, id int64) (dbsqlc.GetLaptopRow, error) {
	a := m.Called(ctx, id)
	v, _ := a.Get(0).(dbsqlc.GetLaptopRow)
	return v, a.Error(1)
}
func (m *Querier) ListLaptops(ctx context.Context) ([]dbsqlc.ListLaptopsRow, error) {
	a := m.Called(ctx)
	v, _ := a.Get(0).([]dbsqlc.ListLaptopsRow)
	return v, a.Error(1)
}
func (m *Querier) UpdateLaptop(ctx context.Context, arg dbsqlc.UpdateLaptopParams) (dbsqlc.Laptop, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Laptop)
	return v, a.Error(1)
}
func (m *Querier) CreateOS(ctx context.Context, name string) (dbsqlc.O, error) {
	a := m.Called(ctx, name)
	v, _ := a.Get(0).(dbsqlc.O)
	return v, a.Error(1)
}
func (m *Querier) DeleteOS(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *Querier) ListOS(ctx context.Context) ([]dbsqlc.O, error) {
	a := m.Called(ctx)
	v, _ := a.Get(0).([]dbsqlc.O)
	return v, a.Error(1)
}
func (m *Querier) UpdateOS(ctx context.Context, arg dbsqlc.UpdateOSParams) (dbsqlc.O, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.O)
	return v, a.Error(1)
}
func (m *Querier) CreateRole(ctx context.Context, arg dbsqlc.CreateRoleParams) (dbsqlc.Role, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Role)
	return v, a.Error(1)
}
func (m *Querier) DeleteRole(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *Querier) ListRoles(ctx context.Context) ([]dbsqlc.Role, error) {
	a := m.Called(ctx)
	v, _ := a.Get(0).([]dbsqlc.Role)
	return v, a.Error(1)
}
func (m *Querier) CreateRoom(ctx context.Context, arg dbsqlc.CreateRoomParams) (dbsqlc.Room, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Room)
	return v, a.Error(1)
}
func (m *Querier) DeleteRoom(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *Querier) ListRoomsByCenter(ctx context.Context, centerID int64) ([]dbsqlc.Room, error) {
	a := m.Called(ctx, centerID)
	v, _ := a.Get(0).([]dbsqlc.Room)
	return v, a.Error(1)
}
func (m *Querier) UpdateRoom(ctx context.Context, arg dbsqlc.UpdateRoomParams) (dbsqlc.Room, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Room)
	return v, a.Error(1)
}
func (m *Querier) CreateClass(ctx context.Context, arg dbsqlc.CreateClassParams) (dbsqlc.SchoolClass, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.SchoolClass)
	return v, a.Error(1)
}
func (m *Querier) DeleteClass(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *Querier) GetClass(ctx context.Context, id int64) (dbsqlc.SchoolClass, error) {
	a := m.Called(ctx, id)
	v, _ := a.Get(0).(dbsqlc.SchoolClass)
	return v, a.Error(1)
}
func (m *Querier) ListClasses(ctx context.Context) ([]dbsqlc.ListClassesRow, error) {
	a := m.Called(ctx)
	v, _ := a.Get(0).([]dbsqlc.ListClassesRow)
	return v, a.Error(1)
}
func (m *Querier) ListClassesByCycle(ctx context.Context, cycleID int64) ([]dbsqlc.ListClassesByCycleRow, error) {
	a := m.Called(ctx, cycleID)
	v, _ := a.Get(0).([]dbsqlc.ListClassesByCycleRow)
	return v, a.Error(1)
}
func (m *Querier) ListClassesByTutor(ctx context.Context, tutorID pgtype.Int8) ([]dbsqlc.ListClassesByTutorRow, error) {
	a := m.Called(ctx, tutorID)
	v, _ := a.Get(0).([]dbsqlc.ListClassesByTutorRow)
	return v, a.Error(1)
}
func (m *Querier) UpdateClass(ctx context.Context, arg dbsqlc.UpdateClassParams) (dbsqlc.SchoolClass, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.SchoolClass)
	return v, a.Error(1)
}
func (m *Querier) CreateStudent(ctx context.Context, arg dbsqlc.CreateStudentParams) (dbsqlc.Student, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Student)
	return v, a.Error(1)
}
func (m *Querier) DeleteStudent(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *Querier) GetStudent(ctx context.Context, id int64) (dbsqlc.Student, error) {
	a := m.Called(ctx, id)
	v, _ := a.Get(0).(dbsqlc.Student)
	return v, a.Error(1)
}
func (m *Querier) ListStudentsByClass(ctx context.Context, classID int64) ([]dbsqlc.Student, error) {
	a := m.Called(ctx, classID)
	v, _ := a.Get(0).([]dbsqlc.Student)
	return v, a.Error(1)
}
func (m *Querier) UpdateStudent(ctx context.Context, arg dbsqlc.UpdateStudentParams) (dbsqlc.Student, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.Student)
	return v, a.Error(1)
}
func (m *Querier) CreateUser(ctx context.Context, arg dbsqlc.CreateUserParams) (dbsqlc.CreateUserRow, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.CreateUserRow)
	return v, a.Error(1)
}
func (m *Querier) DeleteUser(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *Querier) ListUsers(ctx context.Context) ([]dbsqlc.ListUsersRow, error) {
	a := m.Called(ctx)
	v, _ := a.Get(0).([]dbsqlc.ListUsersRow)
	return v, a.Error(1)
}
func (m *Querier) UpdateUser(ctx context.Context, arg dbsqlc.UpdateUserParams) (dbsqlc.UpdateUserRow, error) {
	a := m.Called(ctx, arg)
	v, _ := a.Get(0).(dbsqlc.UpdateUserRow)
	return v, a.Error(1)
}
