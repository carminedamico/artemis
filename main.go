package main

import (
	"log"

	"github.com/carminedamico/artemis/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
	log.Print("Scheduling completed")
}
