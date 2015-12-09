package random

import "github.com/spf13/cobra"

// RandCmd represents the root command of the random package
var RandCmd = &cobra.Command{
	Use:   "rand",
	Short: "Generate random stuff",
	Long:  `Generate random stuff, iv, secret, etc`,
}
