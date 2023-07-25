package database

import (
	"database/sql"

	"github.com/Ed-cred/SolarPal/config"
)

type SQLiteRepo struct {
	Cfg *config.AppConfig
	DB *sql.DB
}


func NewSQLiteRepo(conn *sql.DB, a *config.AppConfig) *SQLiteRepo {
	return &SQLiteRepo{
		Cfg: a,
		DB: conn,
	}
}	