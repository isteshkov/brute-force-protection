package repositories

import (
	"database/sql"
	"time"

	"gitlab.com/isteshkov/brute-force-protection/domain/database"
	"gitlab.com/isteshkov/brute-force-protection/domain/logging"
	"gitlab.com/isteshkov/brute-force-protection/domain/models"
)

type Subnets interface {
	WithLogger(l logging.Logger) Subnets

	Set(subnet models.Subnet, withTx database.Transaction) (tx database.Transaction, err error)
	SetDeleted(subnet models.Subnet, deletedAt time.Time, withTx database.Transaction) (tx database.Transaction, err error)

	ByAddress(address string) (subnet models.Subnet, err error)
	Whitelist() (whitelist []models.Subnet, err error)
	Blacklist() (blacklist []models.Subnet, err error)
}

const subnetFields = `
		version, created_at, updated_at, 
		address, is_blacklisted
`

type subnetListRepository struct {
	db database.Database
	l  logging.Logger
}

func NewSubnetListRepository(db database.Database, l logging.Logger) Subnets {
	return &subnetListRepository{
		db: db,
		l:  l,
	}
}

func (s subnetListRepository) WithLogger(l logging.Logger) Subnets {
	s.l = l
	s.db = s.db.WithLogger(l)
	return &s
}

func (s subnetListRepository) Set(
	subnet models.Subnet, withTx database.Transaction) (tx database.Transaction, err error) {
	defer processError(&err)

	tx, err = database.GetWithTx(s.db, withTx, s.l)
	if err != nil {
		return
	}

	query := `INSERT INTO subnets(` + subnetFields +
		`)VALUES(1,$2,$2,$3,$4) ON CONFLICT(address) DO UPDATE SET version = 
		subnets.version+1, updated_at=$2, is_blacklisted=$4
		WHERE subnets.version=$1;`

	var result sql.Result
	result, err = tx.Exec(
		query,
		subnet.Version,
		time.Now().UTC(),
		subnet.Address,
		subnet.IsBlacklisted,
	)
	if err != nil {
		return
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return
	}

	if ra != 1 {
		err = ErrorProducerInconsistent.New("wrong version")
		return
	}

	return
}

func (s subnetListRepository) SetDeleted(
	subnet models.Subnet,
	deletedAt time.Time,
	withTx database.Transaction) (tx database.Transaction, err error) {
	defer processError(&err)

	tx, err = database.GetWithTx(s.db, withTx, s.l)
	if err != nil {
		return
	}

	query := `UPDATE subnets SET version=subnets.version+1, deleted_at=$2 WHERE address=$1 AND subnets.version=$3;`

	result, err := tx.Exec(query, subnet.Address, deletedAt.UTC(), subnet.Version)
	if err != nil {
		return
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return
	}

	if ra != 1 {
		err = ErrorProducerInconsistent.New("wrong version")
		return
	}

	return
}

func (s subnetListRepository) ByAddress(address string) (subnet models.Subnet, err error) {
	defer processError(&err)

	stmt, err := s.db.Prepare(`SELECT ` + subnetFields +
		` FROM subnets WHERE address = $1 AND deleted_at IS NULL;`,
	)
	if err != nil {
		database.CloseStmt(stmt, &err)
		return
	}
	defer database.CloseStmt(stmt, &err)

	row := stmt.QueryRow(address)
	err = row.Scan(
		&subnet.Version,
		&subnet.CreatedAt,
		&subnet.UpdatedAt,
		&subnet.Address,
		&subnet.IsBlacklisted,
	)
	if err != nil {
		return
	}

	return
}

func (s subnetListRepository) Whitelist() (whitelist []models.Subnet, err error) {
	defer processError(&err)

	stmt, err := s.db.Prepare(`SELECT ` + subnetFields +
		` FROM subnets WHERE is_blacklisted = false AND deleted_at IS NULL;`)
	if err != nil {
		database.CloseStmt(stmt, &err)
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		database.CloseConnections(stmt, rows, &err)
		return
	}
	defer database.CloseConnections(stmt, rows, &err)

	for rows.Next() {
		var subnet models.Subnet
		err = rows.Scan(
			&subnet.Version,
			&subnet.CreatedAt,
			&subnet.UpdatedAt,
			&subnet.Address,
			&subnet.IsBlacklisted,
		)
		if err != nil {
			return
		}

		whitelist = append(whitelist, subnet)
	}

	return
}

func (s subnetListRepository) Blacklist() (blacklist []models.Subnet, err error) {
	defer processError(&err)

	stmt, err := s.db.Prepare(`SELECT ` + subnetFields +
		` FROM subnets WHERE is_blacklisted = true AND deleted_at IS NULL;`)
	if err != nil {
		database.CloseStmt(stmt, &err)
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		database.CloseConnections(stmt, rows, &err)
		return
	}
	defer database.CloseConnections(stmt, rows, &err)

	for rows.Next() {
		var subnet models.Subnet
		err = rows.Scan(
			&subnet.Version,
			&subnet.CreatedAt,
			&subnet.UpdatedAt,
			&subnet.Address,
			&subnet.IsBlacklisted,
		)
		if err != nil {
			return
		}

		blacklist = append(blacklist, subnet)
	}

	return
}
