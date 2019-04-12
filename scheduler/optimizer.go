package scheduler

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/carminedamico/artemis/config"
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
	var bestAgent Scheduler
	var wg sync.WaitGroup

	for index := range agents {
		agents[index] = clone(optimizer.parent)
	}
	bestAgent = clone(optimizer.parent)

	log.Println("*** Initial power consumption -> ", bestAgent.powerConsumption, "W ***")
	log.Println("Agents created -- STARTING THE OPTIMIZATION PROCESS")

	// steadyState := 0

	for g := 0; g < optimizer.confs.NumberOfGenerations; g++ {

		// if steadyState == int(optimizer.confs.NumberOfGenerations/optimizer.confs.PopulationSize) {
		// 	steadyState = 0
		// 	for index := range agents {
		// 		agents[index] = clone(bestAgent)
		// 	}
		// }

		wg.Add(optimizer.confs.PopulationSize)

		for index := range agents {

			go func(agent *Scheduler) {
				defer wg.Done()

				rand.Seed(time.Now().UTC().UnixNano())

				op := randInt(0, 5)

				switch op {
				case 0:
					agent.TSWPMutation()

				case 1:
					agent.TFFCMutation()

				case 2:
					agent.TBFCMutation()

				case 3:
					agent.SCMutation()

				case 4:
					agent.SLRMutation()
				}

				agent.getPowerConsumptionAccountingMigration(optimizer.parent)
			}(&agents[index])

		}

		wg.Wait()

		for _, agent := range agents {
			if agent.powerConsumption < bestAgent.powerConsumption {
				bestAgent = clone(agent)
			}
		}

	}

	log.Println("New best power consumption -> ", bestAgent.powerConsumption, "W")
}

func clone(src Scheduler) Scheduler {
	var cpy Scheduler

	cpy.datacenter.Servers = make([]config.Server, len(src.datacenter.Servers))
	for i, server := range src.datacenter.Servers {
		cpy.datacenter.Servers[i] = server
	}
	cpy.workload.Tasks = make([]config.Task, len(src.workload.Tasks))
	for i, task := range src.workload.Tasks {
		cpy.workload.Tasks[i] = task
	}

	cpy.powerConsumption = src.powerConsumption

	return cpy
}

func (scheduler *Scheduler) getPowerConsumptionAccountingMigration(parent Scheduler) {
	migrationCost := float32(0)

	for index := range scheduler.workload.Tasks {
		if scheduler.workload.Tasks[index].AllocatedOn != parent.workload.Tasks[index].AllocatedOn {
			task := parent.workload.Tasks[index]
			server := parent.datacenter.Servers[parent.workload.Tasks[index].AllocatedOn]
			migrationCost += ((task.CPU / server.CPU) * server.PowerDC) * float32(0.10)
		}
	}

	scheduler.GetPowerConsumption()
	scheduler.powerConsumption += migrationCost

}
