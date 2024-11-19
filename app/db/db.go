package db

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var db *bun.DB = nil

func GetConn() *bun.DB {
	if db == nil {
		err := connect()
		if err != nil {
			fmt.Println("[DB] DB Connection Failed:", err.Error())
		}
		fmt.Println("[DB] DB Connection Established")
	}
	return db
}

func connect() error {
	// TODO: Env vars
	dbURL := "postgres://postgres:pass@localhost:5432/pets-go"
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbURL), pgdriver.WithInsecure(true)))

	dbConn := bun.NewDB(sqlDb, pgdialect.New())
	db = dbConn
	return nil
}
