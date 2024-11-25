package db

import (
	"app/utils"
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

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		utils.DB_USER, utils.DB_PASSWORD, utils.DB_HOST, utils.DB_PORT, utils.DB_NAME,
	)

	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbURL), pgdriver.WithInsecure(true)))

	dbConn := bun.NewDB(sqlDb, pgdialect.New())
	db = dbConn
	return nil
}
