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

// DatabaseExists checks if a database exists
func DatabaseExists(db *sql.DB, databaseName string) (bool, error) {
	rows, err := db.Query("SELECT * FROM pg_database where datname = $1", databaseName)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

// CreateDatabase creates a database
func CreateDatabase(db *sql.DB, databaseName string) error {
	_, err := db.Exec("CREATE DATABASE " + databaseName)
	if err != nil {
		return err
	}
	return nil
}

// UserExists checks if a user exists
func UserExists(db *sql.DB, username string) (bool, error) {
	rows, err := db.Query("SELECT * FROM pg_user WHERE usename=$1", username)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

// CreateUser creates a user
func CreateUser(db *sql.DB, username, password string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE USER \"%s\" WITH NOCREATEDB NOCREATEROLE LOGIN NOREPLICATION NOSUPERUSER PASSWORD '%s'", username, password))
	if err != nil {
		return err
	}
	return nil
}

// GrantFullPrivileges grants full previledges for a user to a database schema
func GrantFullPrivileges(db *sql.DB, username, databaseName, schema string) error {
	queries := []string{
		fmt.Sprintf("GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA %s TO %s", schema, username),
		fmt.Sprintf("GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA %s TO %s", schema, username),
		fmt.Sprintf("GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA %s TO %s", schema, username),
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA %s GRANT ALL PRIVILEGES ON TABLES TO %s", schema, username),
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA %s GRANT ALL PRIVILEGES ON SEQUENCES TO %s", schema, username),
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA %s GRANT ALL PRIVILEGES ON FUNCTIONS TO %s", schema, username),
	}
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}
