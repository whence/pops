package lib

import (
	"database/sql"
	"errors"
	"fmt"

	// just load the pg driver
	_ "github.com/lib/pq"
)

// TryPgConnection tests connection to a postgres database
func TryPgConnection(username, password, host string, port int, database, sslmode string) error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", username, password, host, port, database, sslmode)
	db, err := sql.Open("postgres", connStr)
	defer db.Close()

	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	result, errQuery := db.Exec("SELECT * FROM pg_database where datname = $1", database)
	if errQuery != nil {
		return errQuery
	}
	if rows, errResult := result.RowsAffected(); errResult != nil || rows == 0 {
		return errors.New("Failed to find database " + database)
	}

	return nil
}
