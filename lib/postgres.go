package lib

import (
	"database/sql"
	"fmt"
	"time"

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

// PgConnection connects to a postgres DB
func PgConnection(conn *PostgresConnection) (*sql.DB, error) {
	return sql.Open("postgres", conn.getConnectionString())
}

// TryPgConnection tests connection to a postgres database
func TryPgConnection(conn *PostgresConnection, attempt int) error {
	db, err := PgConnection(conn)
	defer db.Close()
	if err != nil {
		return err
	}

	for i := 0; i < attempt; i++ {
		err = db.Ping()
		if err == nil {
			return nil
		}
		fmt.Println(fmt.Sprintf("Try connecting to %s:%d", conn.Host, conn.Port))
		time.Sleep(1 * time.Second)
	}

	return err
}

// DatabaseExists check if a database exists
func DatabaseExists(db *sql.DB, name string) (bool, error) {
	rows, err := db.Query("SELECT * FROM pg_database where datname = $1", name)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

// CreateDatabase create a database
func CreateDatabase(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}
