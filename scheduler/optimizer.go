package scheduler

import (
	"encoding/json"
	"log"
	"os"

	"github.com/carminedamico/artemis/config"
)

//Optimizer optimizes
type Optimizer struct {
	confs     config.EvolutionaryAlgorithmConfs
	scheduler Scheduler
}

// NewOptimizer creates a nre optimizer
func NewOptimizer(scheduler *Scheduler) *Optimizer {
	var confs config.EvolutionaryAlgorithmConfs
	confsFile, err := os.Open("example/confs.json")
	defer confsFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(confsFile)
	jsonParser.Decode(&confs)

	optimizer := &Optimizer{
		confs : confs,
		scheduler : *scheduler,
	}

	return optimizer
}

// Run starts the optimization process
func (optimizer *Optimizer) Run() {
	for g := 0; g <= optimizer.confs.NumberOfGenerations; g++ {
		for a := 0; a <= optimizer.confs.PopulationSize; a++ {
			
		}
	}
}
