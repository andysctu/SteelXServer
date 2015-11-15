package services

import (
	"database/sql"
)

var database *sql.DB

func InitDBSvc(db *sql.DB) {
	database = db
}

func GetDB() *sql.DB {
	return database
}
