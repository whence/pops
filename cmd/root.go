package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/whence/pops/cmd/db"
	"github.com/whence/pops/cmd/random"
	"github.com/whence/pops/lib"
	"github.com/hashicorp/go-version"
	"github.com/olebedev/config"
	"github.com/spf13/cobra"
)

var flagVersion bool
var flagVerbose bool

var rootCmd = &cobra.Command{
	Use:   "pops",
	Short: "Everything ops",
	Long:  `The CLI for Ops.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !flagVersion {
			return cmd.Usage()
		}

		fmt.Println("Pops version " + lib.VersionNumber)

		if configPath, err := findConfigFile(); err == nil {
			cfg, err := config.ParseYamlFile(configPath)
			if err != nil {
				return err
			}
			if suggestedVersion, err := cfg.String("suggested_version"); err == nil {
				ok, err := checkVersion(suggestedVersion)
				if err != nil {
					return err
				}
				if !ok {
					fmt.Println("Your version is not compliant with " + suggestedVersion + " as specified in " + configPath)
					fmt.Println("Upgrade your version at https://github.com/whence/pops/releases")
				}
			}
		}

		return nil
	},
}

func findConfigFile() (string, error) {
	searchPath, err := filepath.Abs("./.pops.yml")
	if err != nil {
		return "", err
	}

	for i := 0; i < 100; i++ {
		if flagVerbose {
			fmt.Println("searching config at " + searchPath)
		}
		_, err := os.Stat(searchPath)
		if err == nil {
			return searchPath, nil // found
		}

		dir := filepath.Dir(searchPath)
		// at filesystem root...not found
		if strings.HasSuffix(dir, "/") || strings.HasSuffix(dir, "\\") {
			return "", errors.New("No config found")
		}

		if os.IsNotExist(err) {
			searchPath = filepath.Clean(dir + "/../.pops.yml")
			continue
		}

		return "", err
	}

	return "", errors.New("No config found")
}

func checkVersion(suggestedVersion string) (bool, error) {
	constraints, err := version.NewConstraint(suggestedVersion)
	if err != nil {
		return false, err
	}
	version, err := version.NewVersion(lib.VersionNumber)
	if err != nil {
		return false, err
	}
	return constraints.Check(version), nil
}

func init() {
	rootCmd.Flags().BoolVar(&flagVersion, "version", false, "Print version number")
	rootCmd.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "Log more for what is happening")
	rootCmd.AddCommand(db.DbCmd)
	rootCmd.AddCommand(random.RandCmd)
}

// Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
