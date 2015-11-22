package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/whence/pops/lib"
)

var flagVersion bool

var RootCmd = &cobra.Command{
	Use:   "pops",
	Short: "Everything ops",
	Long: `The CLI for Ops.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagVersion {
			fmt.Printf("Pops version %s", lib.VersionNumber)
			fmt.Println()
			return nil
		} else {
			return cmd.Usage()
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.Flags().BoolVar(&flagVersion, "version", false, "Print version number")
}
