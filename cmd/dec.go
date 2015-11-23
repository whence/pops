package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/whence/pops/lib"
)

var flagSecret string

var decCmd = &cobra.Command{
	Use:   "dec",
	Short: "Decrypt Chef Data Bag",
	Long: `Decrypt Chef data bag using a secret file.
Outputs the result to STDOUT.
Currently only Ver.1 data bags are supported. `,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Please specify the file to decrypt.")
		} else if flagSecret == "" {
			return errors.New("Please specify the path of the secret file.")
		}

		fmt.Println(lib.Decrypt(args[0]))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(decCmd)
	decCmd.Flags().StringVarP(&flagSecret, "secret", "s", "", "Path to the secret file")

	decCmd.SetUsageTemplate(`Usage:
  pops dec [flags] file
{{ if .HasLocalFlags}}
Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}
`)
}
