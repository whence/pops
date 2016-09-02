package cmd

import (
	"errors"
	"fmt"

	"github.com/whence/pops/lib"
	"github.com/spf13/cobra"
)

var flagDecSecret string

var decCmd = &cobra.Command{
	Use:   "dec",
	Short: "Decrypt Chef Data Bag",
	Long: `Decrypt Chef data bag using a secret file.
Outputs the result to STDOUT.
Currently only Ver.1 data bags are supported. `,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Please specify the file to decrypt.")
		} else if flagDecSecret == "" {
			return errors.New("Please specify the path of the secret file.")
		}

		fmt.Println(lib.Decrypt(args[0], flagDecSecret))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(decCmd)
	decCmd.Flags().StringVarP(&flagDecSecret, "secret", "s", "", "Path to the secret file")

	decCmd.SetUsageTemplate(`Usage:
  pops dec [flags] file
{{ if .HasLocalFlags}}
Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}
`)
}
