package testmock

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

// MockDB implements handler.DB (i.e. a pool that exposes Begin) using testify/mock.
type MockDB struct{ mock.Mock }

func (m *MockDB) Begin(ctx context.Context) (pgx.Tx, error) {
	a := m.Called(ctx)
	tx, _ := a.Get(0).(pgx.Tx)
	return tx, a.Error(1)
}

// ErrRow is a pgx.Row that always returns a preset error on Scan.
type ErrRow struct{ Err error }

func (r ErrRow) Scan(_ ...any) error { return r.Err }

// NoRowsTx implements pgx.Tx. Its QueryRow always returns pgx.ErrNoRows.
// Use it to simulate "record not found" inside a transaction in unit tests.
type NoRowsTx struct{}

func (NoRowsTx) Begin(_ context.Context) (pgx.Tx, error) { return NoRowsTx{}, nil }
func (NoRowsTx) Commit(_ context.Context) error          { return nil }
func (NoRowsTx) Rollback(_ context.Context) error        { return nil }

func (NoRowsTx) QueryRow(_ context.Context, _ string, _ ...any) pgx.Row {
	return ErrRow{pgx.ErrNoRows}
}
func (NoRowsTx) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return nil, pgx.ErrNoRows
}
func (NoRowsTx) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}
func (NoRowsTx) CopyFrom(_ context.Context, _ pgx.Identifier, _ []string, _ pgx.CopyFromSource) (int64, error) {
	return 0, nil
}
func (NoRowsTx) SendBatch(_ context.Context, _ *pgx.Batch) pgx.BatchResults { return nil }
func (NoRowsTx) LargeObjects() pgx.LargeObjects                             { return pgx.LargeObjects{} }
func (NoRowsTx) Prepare(_ context.Context, _, _ string) (*pgconn.StatementDescription, error) {
	return nil, nil
}
func (NoRowsTx) Conn() *pgx.Conn { return nil }
