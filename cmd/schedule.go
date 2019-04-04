package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/carminedamico/artemis/config"
	"github.com/carminedamico/artemis/scheduler"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.Flags().StringP("datacenter", "d", "", "Provide the json describing your datacenter's nodes")
	scheduleCmd.Flags().StringP("workload", "w", "", "Provide the json describing your workload's tasks")
}

var scheduleCmd = &cobra.Command{
	Use:   "schedule [-d datacenter.json] [-w workload.json] ",
	Short: "Schedule the workload passed as argument",
	Run: func(cmd *cobra.Command, args []string) {
		datacenterInfoFilename, err := cmd.Flags().GetString("datacenter")
		workloadInfoFilename, err := cmd.Flags().GetString("workload")
		if err != nil {
			log.Fatal(err)
		}
		var datacenter config.Datacenter
		var workload config.Workload
		datacenterFile, err := os.Open(datacenterInfoFilename)
		defer datacenterFile.Close()
		workloadFile, err := os.Open(workloadInfoFilename)
		defer workloadFile.Close()
		if err != nil {
			fmt.Println(err)
		}
		jsonParser := json.NewDecoder(datacenterFile)
		jsonParser.Decode(&datacenter)
		jsonParser = json.NewDecoder(workloadFile)
		jsonParser.Decode(&workload)

		scheduler := scheduler.NewScheduler(datacenter, workload)
		scheduler.Run()
	},
}
