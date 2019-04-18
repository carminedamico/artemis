package main

import (
	"log"

	"github.com/carminedamico/artemis/cmd"
)

func main() {

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
	log.Print("Scheduling completed")
}
