package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type DB struct {
	SQL *sql.DB
}
var dbConn = &DB{}

const maxOpenConns = 10
const maxIdleConns = 5
const maxDBConnLifetime = 3 * time.Minute

func ConnectDb(path string) (*DB, error){
	db, err := NewDB(path)
	if err != nil {
		log.Fatal( "Failed to open sqlite database: ", err )
		panic(err)
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(maxDBConnLifetime)
	dbConn.SQL = db
	err = TestDb(db)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

func TestDb(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

func NewDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}