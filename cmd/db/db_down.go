package db

import (
	"errors"
	"fmt"

	"github.com/MYOB-Technology/pops/lib"
	"github.com/spf13/cobra"
)

var dbDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Destroy the database",
	Long:  `Stop and destroy the database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagDriver == "" {
			return errors.New("Please specify the driver to use.")
		} else if flagDriver == "local-docker-pg" {
			return downLocalDockerPg()
		} else {
			return errors.New("Unknown driver.")
		}
	},
	SilenceErrors: true,
	Aliases:       []string{"destroy"},
}

func downLocalDockerPg() error {
	if err := lib.EnsureDockerWorking(); err != nil {
		return err
	}

	if lib.IsContainerExist(flagContainerName) {
		if err := lib.RemoveContainer(flagContainerName); err != nil {
			return err
		}
		fmt.Println("Removed container " + flagContainerName)
	} else {
		fmt.Println("Container " + flagContainerName + " is not running. No action required.")
	}

	return nil
}

func init() {
	DbCmd.AddCommand(dbDownCmd)
}
