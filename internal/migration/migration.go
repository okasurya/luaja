package migration

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var migration = [...]string{Query01}

// DoMigrate ...
func DoMigrate(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	currentVersion, err := getCurrentVersion(tx)
	if err != nil {
		return err
	}
	codeVersion := getCodeVersion()

	for version := currentVersion; version < codeVersion; version++ {
		_, err = tx.Exec(migration[version])
		if err != nil {
			return err
		}
		_, err = tx.Exec(`
    		UPDATE meta SET VALUE = $1 WHERE key = 'db-version'
		`, version+1)
		if err != nil {
			return err
		}
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

func getCodeVersion() int {
	return len(migration)
}
