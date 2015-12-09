package random

import (
	"fmt"

	"github.com/MYOB-Technology/pops/lib"
	"github.com/spf13/cobra"
)

var flagRandSecretSize int
var flagRandSecretBase64 bool
var flagRandSecretNewLine string

var randSecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Create a random secret",
	Long: `Create a random secret. Print to STDOUT.
  If you pipe to a file, we recommend you to chmod the file to 400`,
	Run: func(cmd *cobra.Command, args []string) {
		b := lib.RandomBytes(flagRandSecretSize)
		if flagRandSecretBase64 {
			fmt.Print(lib.EncodeBase64(b))
		} else {
			fmt.Print(string(b))
		}
		fmt.Print(flagRandSecretNewLine)
	},
}

func init() {
	RandCmd.AddCommand(randSecretCmd)
	randSecretCmd.Flags().IntVar(&flagRandSecretSize, "size", 512, "Size of the secret")
	randSecretCmd.Flags().BoolVar(&flagRandSecretBase64, "base64", true, "Whether to encode output in base64")
	randSecretCmd.Flags().StringVar(&flagRandSecretNewLine, "newline", "", "Optionally append a newline to the end of the output")
}
