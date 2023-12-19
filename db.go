package main

import (
	"github.com/jmoiron/sqlx"
)

func sqlxConnect() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", "root:123@tcp/students_management")

	if err != nil {
		return nil, err
	}
	return db, nil
}
