package commands

import (
	"fmt"

	"github.com/alireza-ahmadi/hoor/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information of the Hoor",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Hoor %s", version.Version)
	},
}
