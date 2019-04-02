package cmd

import (
	"fmt"

	"github.com/carminedamico/artemis/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scheduleCmd)
}

var scheduleCmd = &cobra.Command{
	Use:   "schedule [workload to schedule]",
	Short: "Schedule the workload passed as argument",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var tmp = config.Server{
			8,
			32,
			1000,
		}
		fmt.Println(tmp)
	},
}
