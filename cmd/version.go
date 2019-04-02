package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "Artemis 0.1 (pre-alpha)"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Artemis' version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
