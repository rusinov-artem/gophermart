package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/gophermart/app/dto"
)

type StorageTestSuite struct {
	suite.Suite
	ctx     context.Context
	storage *RegistrationStorage
	pool    *fakePool
}

func Test_Storage(t *testing.T) {
	suite.Run(t, &StorageTestSuite{})
}

func (s *StorageTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.storage = NewRegistrationStorage(s.ctx, nil)
	s.pool = &fakePool{}
	s.storage.pool = s.pool
}

func (s *StorageTestSuite) Test_AddOrderError() {
	s.pool.execErr = fmt.Errorf("database error")
	err := s.storage.AddOrder("", "")
	s.ErrorContains(err, "unable to add order")
}

func (s *StorageTestSuite) Test_AddTokenError() {
	s.pool.execErr = fmt.Errorf("database error")
	err := s.storage.AddToken("", "")
	s.ErrorContains(err, "unable to add auth_token")
}

func (s *StorageTestSuite) Test_FindOrder_QueryError() {
	s.pool.queryErr = fmt.Errorf("database error")
	_, err := s.storage.FindOrder("")
	s.ErrorContains(err, "unable to find order")
}

func (s *StorageTestSuite) Test_FindOrder_RowsClosed() {
	rows := &spyRows{
		tag:     pgconn.CommandTag{},
		scanErr: fmt.Errorf("unable to scan"),
	}

	s.pool.rows = rows

	_, err := s.storage.FindOrder("")
	s.ErrorContains(err, "unable to find order")
	s.True(rows.IsClosed)
}

func (s *StorageTestSuite) Test_FindToken_QueryError() {
	s.pool.queryErr = fmt.Errorf("database error")
	_, err := s.storage.FindToken("")
	s.ErrorContains(err, "unable to find token")
}

func (s *StorageTestSuite) Test_FindToken_RowsClosed() {
	rows := &spyRows{
		tag:     pgconn.CommandTag{},
		scanErr: fmt.Errorf("unable to scan"),
	}

	s.pool.rows = rows

	_, err := s.storage.FindToken("")
	s.ErrorContains(err, "unable to find token")
	s.True(rows.IsClosed)
}

func (s *StorageTestSuite) Test_FindUser_QueryError() {
	s.pool.queryErr = fmt.Errorf("database error")
	_, err := s.storage.FindUser("")
	s.ErrorContains(err, "unable to find user")
}

func (s *StorageTestSuite) Test_FindUser_RowsClosed() {
	rows := &spyRows{
		tag:     pgconn.CommandTag{},
		scanErr: fmt.Errorf("unable to scan"),
	}

	s.pool.rows = rows

	_, err := s.storage.FindUser("")
	s.ErrorContains(err, "unable to find user")
	s.True(rows.IsClosed)
}

func (s *StorageTestSuite) Test_IsLoginExists_QueryError() {
	s.pool.queryErr = fmt.Errorf("database error")
	_, err := s.storage.IsLoginExists("")
	s.ErrorContains(err, "unable to find user")
}

func (s *StorageTestSuite) Test_IsLoginExists_RowsClosed() {
	rows := &spyRows{
		tag:     pgconn.CommandTag{},
		scanErr: fmt.Errorf("unable to scan"),
	}

	s.pool.rows = rows

	_, _ = s.storage.IsLoginExists("")
	s.True(rows.IsClosed)
}

func (s *StorageTestSuite) Test_ListOrders_QueryError() {
	s.pool.queryErr = fmt.Errorf("database error")
	_, err := s.storage.ListOrders("")
	s.ErrorContains(err, "unable to list orders")
}

func (s *StorageTestSuite) Test_ListOrders_RowsClosed() {
	rows := &spyRows{
		tag:     pgconn.CommandTag{},
		scanErr: fmt.Errorf("unable to scan"),
	}

	s.pool.rows = rows

	_, _ = s.storage.ListOrders("")
	s.True(rows.IsClosed)
}

func (s *StorageTestSuite) Test_SaveUserError() {
	s.pool.execErr = fmt.Errorf("database error")
	err := s.storage.SaveUser("", "")
	s.ErrorContains(err, "unable to save user")
}

func (s *StorageTestSuite) Test_updateOrder_CloseBatch() {
	batch := &spyBatchRes{IsClosed: false}
	s.pool.batchRes = batch
	_ = s.storage.UpdateOrdersState([]dto.OrderListItem{
		{},
	})
	s.True(batch.IsClosed)

}

type fakePool struct {
	execErr error
	tag     pgconn.CommandTag

	rows     *spyRows
	queryErr error

	batchRes *spyBatchRes
}

func (f *fakePool) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return f.tag, f.execErr
}

func (f *fakePool) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return f.rows, f.queryErr
}

func (f *fakePool) SendBatch(_ context.Context, _ *pgx.Batch) pgx.BatchResults {
	return f.batchRes
}

func (f *fakePool) Begin(_ context.Context) (pgx.Tx, error) {
	return nil, nil
}

type spyRows struct {
	IsClosed bool
	err      error

	tag     pgconn.CommandTag
	scanErr error
	noNext  bool
}

func (s *spyRows) Close() {
	s.IsClosed = true
}

func (s *spyRows) Err() error {
	return s.err
}

func (s *spyRows) CommandTag() pgconn.CommandTag {
	return s.tag
}

func (s *spyRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (s *spyRows) Next() bool {
	return !s.noNext
}

func (s *spyRows) Scan(_ ...any) error {
	return s.scanErr
}

func (s *spyRows) Values() ([]any, error) {
	return nil, nil
}

func (s *spyRows) RawValues() [][]byte {
	return nil
}

func (s *spyRows) Conn() *pgx.Conn {
	return nil
}

type spyBatchRes struct {
	IsClosed bool
	tag      pgconn.CommandTag
	rows     *spyRows
}

func (s *spyBatchRes) Exec() (pgconn.CommandTag, error) {
	panic("implement me")
}

func (s *spyBatchRes) Query() (pgx.Rows, error) {
	panic("implement me")
}

func (s *spyBatchRes) QueryRow() pgx.Row {
	panic("implement me")
}

func (s *spyBatchRes) Close() error {
	s.IsClosed = true
	return nil
}

type spyTx struct {
	rows     *spyRows
	queryErr error
}

func (s *spyTx) Begin(ctx context.Context) (pgx.Tx, error) {
	//TODO implement me
	panic("implement me")
}

func (s *spyTx) Commit(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s *spyTx) Rollback(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s *spyTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *spyTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	//TODO implement me
	panic("implement me")
}

func (s *spyTx) LargeObjects() pgx.LargeObjects {
	//TODO implement me
	panic("implement me")
}

func (s *spyTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	//TODO implement me
	panic("implement me")
}

func (s *spyTx) Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *spyTx) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return s.rows, s.queryErr
}

func (s *spyTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	//TODO implement me
	panic("implement me")
}

func (s *spyTx) Conn() *pgx.Conn {
	//TODO implement me
	panic("implement me")
}
