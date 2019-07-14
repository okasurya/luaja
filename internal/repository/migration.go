package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func DoMigrate(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Exec("SELECT 1"); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func getCurrentVersion(tx *sql.Tx) (int, error) {
	_, err := tx.Exec(`
    CREATE TABLE IF NOT EXISTS meta (
      key VARCHAR(50) PRIMARY KEY,
      value integer NOT NULL DEFAULT 0
    );
	`)
	if err != nil {
		return -1, err
	}

	_, err = tx.Exec(`
    INSERT INTO meta (key, value)
                VALUES ('db-version', 0)
                ON CONFLICT DO NOTHING
	`)
	if err != nil {
		return -1, err
	}
	row, err := tx.Query(`SELECT value FROM meta WHERE key = 'db-version'`)
	var version int
	for row.Next() {
		err = row.Scan(&version)
		if err != nil {
			return -1, err
		}
	}
	return version, nil
}
