package database

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"gitlab.com/isteshkov/brute-force-protection/domain/logging"
)

type Database interface {
	WithLogger(l logging.Logger) Database

	Ping() (err error)
	Begin() (tx *sql.Tx, err error)
	Prepare(query string) (stmt *Stmt, err error)

	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

type sqlDatabase struct {
	client    *sql.DB
	newClient *sqlx.DB
	logger    logging.Logger
}

func (s sqlDatabase) WithLogger(l logging.Logger) Database {
	s.logger = l
	return &s
}

func (s *sqlDatabase) Client() *sql.DB {
	return s.client
}

func (s *sqlDatabase) Begin() (tx *sql.Tx, err error) {
	defer processError(&err)

	tx, err = s.client.Begin()
	if err != nil {
		return
	}

	return
}

//nolint:sqlclosecheck
func (s *sqlDatabase) Prepare(query string) (stmt *Stmt, err error) {
	defer processError(&err)

	sqlStmt, err := s.client.Prepare(query)
	if err != nil {
		return
	}

	stmt = &Stmt{
		l:     s.logger,
		stmt:  sqlStmt,
		query: query,
	}

	return
}

func (s *sqlDatabase) Ping() (err error) {
	defer processError(&err)

	pong := "pong"
	stmt, err := s.client.Prepare("SELECT $1;")
	if err != nil {
		if stmt != nil {
			err = stmt.Close()
		}
		return
	}
	defer func() {
		if stmt != nil {
			err = stmt.Close()
			return
		}
	}()

	row := stmt.QueryRow(pong)
	var selectedPong string
	err = row.Scan(&selectedPong)
	if err != nil {
		return
	}

	if selectedPong != pong {
		err = ErrorProducerGeneral.New("something wrong with db")
		return
	}

	return
}

func (s *sqlDatabase) Get(dest interface{}, query string, args ...interface{}) (err error) {
	defer processError(&err)
	err = s.newClient.Get(dest, query, args...)

	return err
}

func (s *sqlDatabase) Select(dest interface{}, query string, args ...interface{}) (err error) {
	defer processError(&err)
	err = s.newClient.Select(dest, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	return err
}
