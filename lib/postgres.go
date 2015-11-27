package lib

import (
	"database/sql"
	"errors"
	"fmt"

	// just load the pg driver
	_ "github.com/lib/pq"
)

// PostgresConnection represents information about a postgres connection
type PostgresConnection struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SslMode  string
}

func (c *PostgresConnection) getConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.SslMode,
	)
}

// TryPgConnection tests connection to a postgres database
func TryPgConnection(conn *PostgresConnection) error {
	db, err := sql.Open("postgres", conn.getConnectionString())
	defer db.Close()
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	result, errQuery := db.Exec("SELECT * FROM pg_database where datname = $1", conn.Database)
	if errQuery != nil {
		return errQuery
	}
	if rows, errResult := result.RowsAffected(); errResult != nil || rows == 0 {
		return errors.New("Failed to find database " + conn.Database)
	}

	return nil
}

func initialiseForApp(conn *PostgresConnection, appDatabase, appUsername, appPassword string) error {
	db, err := sql.Open("postgres", conn.getConnectionString())
	defer db.Close()
	if err != nil {
		return err
	}
	return nil
}
