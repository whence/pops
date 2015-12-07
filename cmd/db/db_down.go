package db

import (
	"errors"
	"fmt"

	"github.com/MYOB-Technology/pops/lib"
	"github.com/spf13/cobra"
)

var flagDownDriver string
var flagDownContainerName string

var dbDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Destroy the database",
	Long:  `Stop and destroy the database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagDownDriver == "" {
			return errors.New("Please specify the driver to use.")
		} else if flagDownDriver == "local-docker-pg" {
			return downLocalDockerPg()
		} else {
			return errors.New("Unknown driver.")
		}
	},
	SilenceErrors: true,
	Aliases:       []string{"destroy"},
}

func downLocalDockerPg() error {
	containerName := flagDownContainerName

	if err := lib.EnsureDockerWorking(); err != nil {
		return err
	}

	if lib.IsContainerExist(containerName) {
		if err := lib.RemoveContainer(containerName); err != nil {
			return err
		}
		fmt.Println("Removed container " + containerName)
	} else {
		fmt.Println("Container " + containerName + " is not running. No action required.")
	}

	return nil
}

func init() {
	DbCmd.AddCommand(dbDownCmd)
	dbDownCmd.Flags().StringVarP(&flagDownDriver, "driver", "d", "", "The driver to use to control the database. Currently only local-docker-pg is supported.")
	dbDownCmd.Flags().StringVar(&flagDownContainerName, "container", "pops-db", "The name of container to run. Applicable to docker drivers only.")
}
