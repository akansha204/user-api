package repository

import (
	"database/sql"
)

func OpenDB(driverName, dsn string) (*sql.DB, error) {
	return sql.Open(driverName, dsn)
}
