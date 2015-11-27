package db

import "github.com/spf13/cobra"

var flagDriver string
var flagContainerName string

// DbCmd represents the root command of the db package
var DbCmd = &cobra.Command{
	Use:   "db",
	Short: "Control database",
	Long:  `Control database`,
}

func init() {
	DbCmd.PersistentFlags().StringVarP(&flagDriver, "driver", "d", "", "The driver to use to control the datbase. Currently only local-docker-pg is supported.")
	DbCmd.PersistentFlags().StringVar(&flagContainerName, "container", "pops-db", "The name of container to run. Applicable to docker drivers only.")
}
