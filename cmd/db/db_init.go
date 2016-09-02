package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/whence/pops/lib"
	"github.com/spf13/cobra"
)

var flagInitMasterUsername string
var flagInitMasterPassword string
var flagInitDbHost string
var flagInitDbPort int
var flagAppDatabase string
var flagAppSchema string
var flagAppUsername string
var flagAppPassword string
var flagInitDbSslMode string

var dbInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the database",
	Long:  `Initialize the database so it is ready to be used by the app`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagAppDatabase == "" {
			return errors.New("Please specify app database to initialize.")
		}
		return initPg()
	},
	SilenceErrors: true,
}

func initPg() error {
	var dbPort int
	if flagInitDbPort == -1 {
		dbPort = 5432
	} else {
		dbPort = flagInitDbPort
	}

	conn := &lib.PostgresConnection{
		Username: flagInitMasterUsername,
		Password: flagInitMasterPassword,
		Host:     flagInitDbHost,
		Port:     dbPort,
		Database: "postgres",
		SslMode:  flagInitDbSslMode,
	}

	db, err := lib.PgConnection(conn)
	defer db.Close()
	if err != nil {
		return err
	}

	if err := createAppDatabase(db, flagAppDatabase); err != nil {
		return err
	}

	if err := createAppUser(db, flagAppUsername, flagAppPassword, flagAppDatabase, flagAppSchema); err != nil {
		return err
	}

	return nil
}

func createAppDatabase(db *sql.DB, databaseName string) error {
	exists, err := lib.DatabaseExists(db, databaseName)
	if err != nil {
		return err
	}
	if exists {
		fmt.Println("Database " + databaseName + " already exists")
	} else {
		if err := lib.CreateDatabase(db, databaseName); err != nil {
			return err
		}
		fmt.Println("Database " + databaseName + " created")
	}
	return nil
}

func createAppUser(db *sql.DB, username, password, databaseName, schema string) error {
	exists, err := lib.UserExists(db, username)
	if err != nil {
		return err
	}
	if exists {
		fmt.Println("User " + username + " already exists")
	} else {
		if err := lib.CreateUser(db, username, password); err != nil {
			return err
		}
		fmt.Println("User " + username + " created")

		if err := lib.GrantFullPrivileges(db, username, databaseName, schema); err != nil {
			return err
		}
		fmt.Printf("Granted full privileges to user %s on database %s schema %s", username, databaseName, schema)
		fmt.Println()
	}
	return nil
}

func init() {
	DbCmd.AddCommand(dbInitCmd)
	dbInitCmd.Flags().StringVar(&flagInitMasterUsername, "master-username", "postgres", "The master username of database server.")
	dbInitCmd.Flags().StringVar(&flagInitMasterPassword, "master-password", "mysecretpassword", "The master password of database server.")
	dbInitCmd.Flags().StringVar(&flagInitDbHost, "host", "localhost", "The database host")
	dbInitCmd.Flags().IntVarP(&flagInitDbPort, "port", "p", -1, "The database port to run the database. Defaults to the database default port. e.g. Postgres is 5432")
	dbInitCmd.Flags().StringVar(&flagAppDatabase, "app-database", "", "The application database to create.")
	dbInitCmd.Flags().StringVar(&flagAppSchema, "app-schema", "public", "The application schema to create.")
	dbInitCmd.Flags().StringVar(&flagAppUsername, "app-username", "app", "The application username of application database to create.")
	dbInitCmd.Flags().StringVar(&flagAppPassword, "app-password", "mysecretpassword", "The application password of application database to create.")
	dbInitCmd.Flags().StringVar(&flagInitDbSslMode, "ssl-mode", "require", "SSL mode for some drivers, such as Postgres.")
}
