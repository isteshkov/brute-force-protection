package database

import (
	"database/sql"
	"math"
	"time"

	"gitlab.com/isteshkov/brute-force-protection/domain/logging"
)

type Stmt struct {
	l     logging.Logger
	stmt  *sql.Stmt
	query string
}

func (s *Stmt) Query(args ...interface{}) (rows *Rows, err error) {
	defer processError(&err)

	tsBeforeRequest := time.Now().UTC()
	sqlRows, err := s.stmt.Query(args...)
	latency := math.Floor(time.Now().UTC().Sub(tsBeforeRequest).Seconds()*1000) / 1000
	rows = &Rows{
		rows: sqlRows,
	}

	s.l.WithField("query", s.query).WithField("args", args).WithField("db_latency", latency).Debug("db profiling")

	return
}

func (s *Stmt) QueryRow(args ...interface{}) (row *Row) {
	tsBeforeRequest := time.Now().UTC()
	sqlRow := s.stmt.QueryRow(args...)
	latency := math.Floor(time.Now().UTC().Sub(tsBeforeRequest).Seconds()*1000) / 1000
	row = &Row{
		row: sqlRow,
	}

	s.l.WithField("query", s.query).WithField("args", args).WithField("db_latency", latency).Debug("db profiling")

	return row
}

func (s *Stmt) Close() (err error) {
	defer processError(&err)

	if s.stmt != nil {
		err = s.stmt.Close()
	} else {
		return
	}

	return
}
