package random

import (
	"fmt"

	"github.com/MYOB-Technology/pops/lib"
	"github.com/spf13/cobra"
)

var flagRandIvSize int
var flagRandIvBase64 bool

var randIvCmd = &cobra.Command{
	Use:   "iv",
	Short: "Create a random iv",
	Long:  `Create a random iv. Print to STDOUT`,
	Run: func(cmd *cobra.Command, args []string) {
		b := lib.RandomBytes(flagRandIvSize)
		if flagRandIvBase64 {
			fmt.Println(lib.EncodeBase64(b))
		} else {
			fmt.Println(string(b))
		}
	},
}

func init() {
	RandCmd.AddCommand(randIvCmd)
	randIvCmd.Flags().IntVar(&flagRandIvSize, "size", 16, "Size of the iv")
	randIvCmd.Flags().BoolVar(&flagRandIvBase64, "base64", true, "Whether to encode output in base64")
}
