package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/carminedamico/artemis/config"
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
		filename, _ := cmd.Flags().GetString("datacenter")
		var datacenter config.Datacenter
		datacenterFile, err := os.Open(filename)
		defer datacenterFile.Close()
		if err != nil {
			fmt.Println(err)
		}
		jsonParser := json.NewDecoder(datacenterFile)
		jsonParser.Decode(&datacenter)
		fmt.Println(datacenter.Servers[1])
	},
}
