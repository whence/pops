package db

import (
	"errors"
	"fmt"

	"github.com/MYOB-Technology/pops/lib"
	"github.com/spf13/cobra"
)

var flagMasterUsername string
var flagMasterPassword string
var flagAppDatabase string
var flagAppUsername string
var flagAppPassword string
var flagDbHost string
var flagDbPort int
var flagImageName string
var flagPollDbAttempt int

var dbUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the database",
	Long:  `Create a database ready to be used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagDriver == "" {
			return errors.New("Please specify the driver to use.")
		} else if flagDriver == "local-docker-pg" {
			if flagImageName == "" {
				return errors.New("Please specify the image to use.")
			}
			return upLocalDockerPg()
		} else {
			return errors.New("Unknown driver.")
		}
	},
	SilenceErrors: true,
}

func upLocalDockerPg() error {
	var dbPort int
	if flagDbPort == -1 {
		dbPort = 5432
	} else {
		dbPort = flagDbPort
	}

	if err := lib.EnsureDockerWorking(); err != nil {
		return err
	}

	if !lib.IsContainerExist(flagContainerName) {
		if err := lib.RunContainer(flagContainerName, []string{
			"-e", fmt.Sprintf("POSTGRES_USER=%s", flagMasterUsername),
			"-e", fmt.Sprintf("POSTGRES_PASSWORD=%s", flagMasterPassword),
			"-p", fmt.Sprintf("%d:5432", dbPort),
			"-d",
		}, flagImageName); err != nil {
			return err
		}
		fmt.Println("Running container " + flagContainerName)
	} else {
		fmt.Println("Container " + flagContainerName + " is already running.")
	}

	conn := &lib.PostgresConnection{
		Username: flagMasterUsername,
		Password: flagMasterPassword,
		Host:     flagDbHost,
		Port:     dbPort,
		Database: "postgres",
		SslMode:  "disable",
	}

	if err := lib.TryPgConnection(conn, flagPollDbAttempt); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("%s:%d is ready to use!", conn.Host, conn.Port))

	return nil
}

func init() {
	DbCmd.AddCommand(dbUpCmd)
	dbUpCmd.Flags().StringVar(&flagMasterUsername, "master-username", "postgres", "The master username of database server to create.")
	dbUpCmd.Flags().StringVar(&flagMasterPassword, "master-password", "mysecretpassword", "The master password of database server to create.")
	dbUpCmd.Flags().StringVar(&flagAppDatabase, "app-database", "", "The application database to create.")
	dbUpCmd.Flags().StringVar(&flagAppUsername, "app-username", "app", "The application username of application database to create.")
	dbUpCmd.Flags().StringVar(&flagAppPassword, "app-password", "mysecretpassword", "The application password of application database to create.")
	dbUpCmd.Flags().StringVar(&flagDbHost, "host", "localhost", "The database host")
	dbUpCmd.Flags().IntVarP(&flagDbPort, "port", "p", -1, "The database port to run the datbase. Defaults to the database default port. e.g. Postgres is 5432")
	dbUpCmd.Flags().StringVarP(&flagImageName, "image", "i", "", "The docker image (can append tag) to use for the datbase. Applicable to docker drivers only.")
	dbUpCmd.Flags().IntVar(&flagPollDbAttempt, "attempt", 60, "The number of attempt for trying to connect to the database while starting up.")
}
