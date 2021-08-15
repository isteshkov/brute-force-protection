package database

import "database/sql"

type Row struct {
	row *sql.Row
}

func (r *Row) Scan(dest ...interface{}) (err error) {
	defer processError(&err)

	if r.row != nil {
		err = r.row.Scan(dest...)
	} else {
		err = ErrorProducerGeneral.New("try to scan nil row")
	}

	return
}

type Rows struct {
	rows *sql.Rows
}

func (r *Rows) Scan(dest ...interface{}) (err error) {
	defer processError(&err)

	if r.rows != nil {
		err = r.rows.Scan(dest...)
	} else {
		return
	}

	return
}

func (r *Rows) Close() (err error) {
	defer processError(&err)

	if r.rows != nil {
		err = r.rows.Close()
	} else {
		return
	}

	return
}

func (r *Rows) Next() bool {
	if r.rows != nil {
		return r.rows.Next()
	}
	return false
}
