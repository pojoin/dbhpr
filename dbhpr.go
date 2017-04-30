package dbhpr

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var dbHive map[string]*sql.DB = make(map[string]*sql.DB)

var logger *log.Logger

func NewDB(dbname, driverName, url string) error {
	db, err := sql.Open(driverName, url)
	if err != nil {
		fmt.Errorf("error: %v\n", err)
		return err
	}
	err = db.Ping()
	if err != nil {
		fmt.Errorf("error: %v\n", err)
		return err
	}
	dbHive[dbname] = db
	return nil
}

func GetDB(dbname string) (*sql.DB, error) {
	if db, ok := dbHive[dbname]; ok {
		return db, nil
	}
	return nil, errors.New(dbname + " not found!")
}

func NewHelper(dbname string) Helper {
	return &DBHelper{
		dbname: dbname,
	}
}

func Get(sql string, args ...interface{}) (Row, error) {
	h := NewHelper("default")
	row, err := h.Get(sql, args...)
	return row, err
}

func Query(sql string, args ...interface{}) ([]Row, error) {
	h := NewHelper("default")
	rows, err := h.Query(sql, args...)
	return rows, err
}
