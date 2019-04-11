package scheduler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/carminedamico/artemis/config"
	"github.com/jinzhu/copier"
)

//Optimizer optimizes
type Optimizer struct {
	confs  config.EvolutionaryAlgorithmConfs
	parent Scheduler
}

// NewOptimizer creates a nre optimizer
func NewOptimizer(scheduler *Scheduler) *Optimizer {
	var confs config.EvolutionaryAlgorithmConfs
	confsFile, err := os.Open("confs.json")
	defer confsFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(confsFile)
	jsonParser.Decode(&confs)

	optimizer := &Optimizer{
		confs:  confs,
		parent: *scheduler,
	}

	return optimizer
}

// Run starts the optimization process
func (optimizer *Optimizer) Run() {
	agents := make([]Scheduler, optimizer.confs.PopulationSize)
	var wg sync.WaitGroup

	for index := range agents {
		agents[index] = Scheduler{}
		copier.Copy(&agents[index], &optimizer.parent)
	}
	log.Println("Agents created -- START WITH GENERATIONS")

	for g := 0; g < optimizer.confs.NumberOfGenerations; g++ {
		wg.Add(optimizer.confs.PopulationSize)

		for index := range agents {
			go func(agent *Scheduler) {
				defer wg.Done()
				agent.TSWPMutation()
				for index, task := range agent.workload.Tasks {
					if task.AllocatedOn > 999 {
						fmt.Println("Error on ", index)
					}
				}
			}(&agents[index])
		}

		wg.Wait()

		for index := range agents {
			fmt.Println("\tNEW POWER CONSUMPTION -> ", agents[index].GetPowerConsumption())
		}

	}

	// for g := 0; g <= optimizer.confs.NumberOfGenerations; g++ {
	// 	wg.Add(1)

	// 	for index := range agents {
	// 		// Apply mutation based on probability distribution
	// 		go func(agent *Scheduler) {
	// 			defer wg.Done()
	// 			agent.TSWPMutation()
	// 		}(&agents[index])
	// 	}

	// 	wg.Wait()

	// 	for index, agent := range agents {
	// 		fmt.Printf("\tAGENT #%d POWER CONSUMPTION -> %f\n", index, agent.GetPowerConsumption())
	// 	}

	// 	// calculate fitness
	// 	// check best fitness
	// 	// update for next generation
	// }
}
