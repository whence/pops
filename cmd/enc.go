package cmd

import (
	"errors"
	"fmt"

	"github.com/MYOB-Technology/pops/lib"
	"github.com/spf13/cobra"
)

var flagEncSecret string

var encCmd = &cobra.Command{
	Use:   "enc",
	Short: "Encrypt Chef Data Bag",
	Long: `Encrypt Chef data bag using a secret file.
Outputs the result to STDOUT.
Currently only Ver.1 data bags are supported. `,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Please specify the file to encrypt.")
		} else if flagEncSecret == "" {
			return errors.New("Please specify the path of the secret file.")
		}

		fmt.Println(lib.Encrypt(args[0], flagEncSecret))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(encCmd)
	encCmd.Flags().StringVarP(&flagEncSecret, "secret", "s", "", "Path to the secret file")

	encCmd.SetUsageTemplate(`Usage:
  pops enc [flags] file
{{ if .HasLocalFlags}}
Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}
`)
}
