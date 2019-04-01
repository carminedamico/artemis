package commands

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "artemis",
	Short: "Artemis is a smart energy-efficient scheduler for cloud applications",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		//
	},
}
