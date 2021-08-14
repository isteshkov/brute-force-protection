package database

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
	"gitlab.com/isteshkov/brute-force-protection/domain/logging"
)

type TransactionMock struct {
	mock.Mock
}

func (m *TransactionMock) WithLogger(l logging.Logger) Transaction {
	panic("implement me")
}

func (m *TransactionMock) Query(query string, args ...interface{}) (rows *Rows, err error) {
	resp := m.Called(query, args)
	return resp.Get(0).(*Rows), resp.Error(1)
}

func (m *TransactionMock) Commit() (err error) {
	resp := m.Called()
	return resp.Error(0)
}

func (m *TransactionMock) MustRollBack(entailed string) {
	_ = m.Called(entailed)
}

func (m *TransactionMock) Exec(query string, args ...interface{}) (result sql.Result, err error) {
	resp := m.Called(query, args)
	return resp.Get(0).(sql.Result), resp.Error(1)
}

func (m *TransactionMock) Prepare(query string) (result *Stmt, err error) {
	resp := m.Called(query)
	return resp.Get(0).(*Stmt), resp.Error(1)
}

func (m *TransactionMock) IsInitialized() bool {
	resp := m.Called()
	return resp.Get(0).(bool)
}
