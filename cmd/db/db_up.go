package db

import (
	"errors"
	"fmt"

	"github.com/MYOB-Technology/pops/lib"
	"github.com/spf13/cobra"
)

var flagUpDriver string
var flagUpContainerName string
var flagUpMasterUsername string
var flagUpMasterPassword string
var flagUpDbHost string
var flagUpDbPort int
var flagImageName string
var flagPollDbAttempt int
var flagUpDbSslMode string

var dbUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the database",
	Long:  `Create a database ready to be used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagUpDriver == "" {
			return errors.New("Please specify the driver to use.")
		}

		switch flagUpDriver {
		case "local-docker-pg":
			if flagImageName == "" {
				return errors.New("Please specify the image to use.")
			}
			return upLocalDockerPg()
		default:
			return errors.New("Unknown driver.")
		}
	},
	SilenceErrors: true,
}

func upLocalDockerPg() error {
	var dbPort int
	if flagUpDbPort == -1 {
		dbPort = 5432
	} else {
		dbPort = flagUpDbPort
	}

	conn := &lib.PostgresConnection{
		Username: flagUpMasterUsername,
		Password: flagUpMasterPassword,
		Host:     flagUpDbHost,
		Port:     dbPort,
		Database: "postgres",
		SslMode:  flagUpDbSslMode,
	}

	containerName := flagUpContainerName

	if err := lib.EnsureDockerWorking(); err != nil {
		return err
	}

	if !lib.IsContainerExist(containerName) {
		if err := lib.RunContainer(containerName, []string{
			"-e", fmt.Sprintf("POSTGRES_USER=%s", conn.Username),
			"-e", fmt.Sprintf("POSTGRES_PASSWORD=%s", conn.Password),
			"-p", fmt.Sprintf("%d:5432", conn.Port),
			"-d",
		}, flagImageName); err != nil {
			return err
		}
		fmt.Println("Running container " + containerName)
	} else {
		fmt.Println("Container " + containerName + " is already running.")
	}

	if err := lib.TryPgConnection(conn, flagPollDbAttempt); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("%s:%d is ready to use!", conn.Host, conn.Port))

	return nil
}

func init() {
	DbCmd.AddCommand(dbUpCmd)
	dbUpCmd.Flags().StringVarP(&flagUpDriver, "driver", "d", "", "The driver to use to control the datbase. Currently only local-docker-pg is supported.")
	dbUpCmd.Flags().StringVar(&flagUpContainerName, "container", "pops-db", "The name of container to run. Applicable to docker drivers only.")
	dbUpCmd.Flags().StringVar(&flagUpMasterUsername, "master-username", "postgres", "The master username of database server.")
	dbUpCmd.Flags().StringVar(&flagUpMasterPassword, "master-password", "mysecretpassword", "The master password of database server.")
	dbUpCmd.Flags().StringVar(&flagUpDbHost, "host", "localhost", "The database host")
	dbUpCmd.Flags().IntVarP(&flagUpDbPort, "port", "p", -1, "The database port to run the datbase. Defaults to the database default port. e.g. Postgres is 5432")
	dbUpCmd.Flags().StringVarP(&flagImageName, "image", "i", "", "The docker image (can append tag) to use for the datbase. Applicable to docker drivers only.")
	dbUpCmd.Flags().IntVar(&flagPollDbAttempt, "attempt", 60, "The number of attempt for trying to connect to the database while starting up.")
	dbUpCmd.Flags().StringVar(&flagUpDbSslMode, "ssl-mode", "", "SSL mode for some drivers, such as Postgres.")
}
