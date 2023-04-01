package controller

import "database/sql"

var (
	pool    *sql.DB
	connStr = "postgres://postgres:zkw030813@101.43.168.188:5433/postgres?sslmode=disable"
)

func Init() error {
	var err error
	pool, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return nil
}
