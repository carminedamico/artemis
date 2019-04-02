package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "artemis",
	Short: "Artemis is a smart energy-efficient scheduler for cloud applications",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		// Print something about the usage (usage: artemis ...)
	},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
