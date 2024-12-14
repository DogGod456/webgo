package handlers

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}
