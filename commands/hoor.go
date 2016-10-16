// Package commands define command-line interface for Hoor, current implementation
// is based on the Cobra package.
package commands

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/hugo/hugolib"
	jww "github.com/spf13/jwalterweatherman"
)

// HoorCmd represents the base command when called without any subcommands.
var HoorCmd = &cobra.Command{
	Use:   "hoor",
	Short: "Add Shamsi date to Hugo based websites with ease",
	Run: func(cmd *cobra.Command, args []string) {
		jww.INFO.Println("Hello")
	},
}

// Available command line flags
var (
	cfgFile    string
	source     string
	contentDir string
)

// init initilizes required flags for Hoor.
func init() {
	// Persistent flags
	HoorCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is path/config.yaml|json|toml)")

	validConfigFilenames := []string{"json", "js", "yaml", "yml", "toml", "tml"}
	HoorCmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)

	// Common flags
	HoorCmd.Flags().StringVarP(&source, "source", "s", "", "filesystem path to read files relative from")
	HoorCmd.Flags().StringVarP(&contentDir, "contentDir", "c", "", "filesystem path to content directory")
}

// Setup adds all child commands to the root command and sets flags appropriately.
func Setup() {
	// Load configuration options into Viper
	if err := hugolib.LoadGlobalConfig(source, cfgFile); err != nil {
		jww.ERROR.Println("Cannot find configurations file")
		return
	}

	HoorCmd.AddCommand(versionCmd)

	if err := HoorCmd.Execute(); err != nil {
		jww.ERROR.Println(err)
		os.Exit(-1)
	}
}
