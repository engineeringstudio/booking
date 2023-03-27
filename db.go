package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type dbOperator struct {
	insertCmd *sql.Stmt
	queryCmd  *sql.Stmt
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

	return &dbOperator{
		insertCmd: insertCmd,
		queryCmd:  queryCmd,
	}, nil
}

func (d *dbOperator) insert(r *request) error {
	_, err := d.insertCmd.Exec(r.name, r.sno, r.pn, r.date, r.info)
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
		err = result.Scan(&req.name, &req.sno, &req.pn, &req.date, &req.info)
		if err != nil {
			return nil, err
		}

		tmp = append(tmp, req)
	}

	return &tmp, nil
}
