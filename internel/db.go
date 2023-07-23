package internel

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type dbOperator struct {
	insertCmd *sql.Stmt
	queryCmd  *sql.Stmt
	countCmd  *sql.Stmt
}

func newDbOperator(db *sql.DB, table string) (*dbOperator, error) {
	insertCmd, err := db.Prepare(
		fmt.Sprintf("INSERT INTO %s(NAME,SNO,PN,DATE,INFO)  values(?,?,?,?,?)",
			table))
	if err != nil {
		return nil, err
	}

	queryCmd, err := db.Prepare(
		fmt.Sprintf("SELECT NAME,SNO,PN,DATE,INFO FROM %s WHERE DATE=?",
			table))
	if err != nil {
		return nil, err
	}

	countCmd, err := db.Prepare(
		fmt.Sprintf("SELECT COUNT(DATE) AS COUNT FROM %s WHERE DATE=?",
			table))
	if err != nil {
		return nil, err
	}

	return &dbOperator{
		insertCmd: insertCmd,
		queryCmd:  queryCmd,
		countCmd:  countCmd,
	}, nil
}

func (d *dbOperator) insert(r *request) error {
	_, err := d.insertCmd.Exec(r.Name, r.Sno, r.Pn, r.Date, r.Info)
	return err
}

func (d *dbOperator) query(date string) (*[]request, error) {
	var req request
	tmp := make([]request, 0, 10)

	result, err := d.queryCmd.Query(date)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		err = result.Scan(&req.Name, &req.Sno, &req.Pn, &req.Date, &req.Info)
		if err != nil {
			return nil, err
		}

		tmp = append(tmp, req)
	}

	return &tmp, nil
}

func (d *dbOperator) count(date string) (int, error) {
	result, err := d.queryCmd.Query(date)
	if err != nil {
		return 0, err
	}

	var count int
	result.Next()
	result.Scan(&count)

	return count, nil
}
