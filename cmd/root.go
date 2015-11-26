package cmd

import (
	"fmt"
	"os"

	"github.com/MYOB-Technology/pops/lib"
	"github.com/spf13/cobra"
)

var flagVersion bool

var rootCmd = &cobra.Command{
	Use:   "pops",
	Short: "Everything ops",
	Long:  `The CLI for Ops.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !flagVersion {
			return cmd.Usage()
		}

		fmt.Printf("Pops version %s", lib.VersionNumber)
		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.Flags().BoolVar(&flagVersion, "version", false, "Print version number")
}

// Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
