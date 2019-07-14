package repository

import (
	"database/sql"
	"testing"
)

func getDB() (*sql.DB, error) {
	connStr := "user=luaja password='password' dbname=luaja sslmode=disable"
	return sql.Open("postgres", connStr)
}
func TestDoMigrate(t *testing.T) {
	db, err := getDB()
	if err != nil {
		t.Error(err)
	}
	err = DoMigrate(db)
	if err != nil {
		t.Error(err)
	}
}

func TestGetCurrentVersion(t *testing.T) {
	db, err := getDB()
	if err != nil {
		t.Error(err)
	}
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		t.Error(err)
	}
	version, err := getCurrentVersion(tx)
	if err != nil {
		t.Error(err)
	}

	if version != 0 {
		t.Errorf("failed, expected 0, actual %d", version)
	}
}
