package commands

import (
	"github.com/alireza-ahmadi/hoor/version"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

// versionCmd returns the current application version.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information of the Hoor",
	Run: func(cmd *cobra.Command, args []string) {
		jww.FEEDBACK.Println("Hoor", version.Version)
	},
}
