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
		fmt.Sprintf("SELECT NAME,SNO,PN,DATE,INFO FROM %s ORDER BY ID DESC",
			table))
	if err != nil {
		return nil, err
	}

	return &dbOperator{
		insertCmd: insertCmd,
		queryCmd:  queryCmd,
	}, nil
}

func (d *dbOperator) insert(date, baseUrl string) error {
	_, err := d.insertCmd.Exec(date, baseUrl)
	return err
}

func (d *dbOperator) query(num int) (*[]list, error) {
	return nil, nil
}
