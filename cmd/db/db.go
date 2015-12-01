package db

import "github.com/spf13/cobra"

// DbCmd represents the root command of the db package
var DbCmd = &cobra.Command{
	Use:   "db",
	Short: "Control database",
	Long:  `Control database`,
}
