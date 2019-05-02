package scheduler

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"sync"

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

	agents := make([]Agent, optimizer.confs.PopulationSize)

	var bestScheduler Scheduler
	var wg sync.WaitGroup

	f, err := os.OpenFile(optimizer.confs.LogFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	for index := range agents {
		agents[index] = *newRandomAgent(optimizer.parent)
	}

	sort.Slice(agents, func(i, j int) bool {
		return agents[i].scheduler.powerConsumption < agents[j].scheduler.powerConsumption
	})

	bestScheduler = clone(agents[0].scheduler)
	steadyState := 0

	log.Println("Agents created -- STARTING THE OPTIMIZATION PROCESS")

	for g := 0; g < optimizer.confs.NumberOfGenerations; g++ {

		log.Println(g, " -- ", bestScheduler.powerConsumption)

		if steadyState >= 30 && g > 0 {
			log.Println("RESTART")
			for index := range agents {
				agents[index].scheduler.randomizer()
			}
			steadyState = 0
		}

		sort.Slice(agents, func(i, j int) bool {
			return agents[i].scheduler.powerConsumption < agents[j].scheduler.powerConsumption
		})

		agents[3].scheduler = clone(agents[0].scheduler)
		agents[4].scheduler = clone(agents[0].scheduler)
		agents[5].scheduler = clone(agents[1].scheduler)
		agents[6].scheduler.randomizer()
		agents[6].scheduler.getPowerConsumptionAccountingMigration(optimizer.parent)
		agents[7].scheduler.randomizer()
		agents[7].scheduler.getPowerConsumptionAccountingMigration(optimizer.parent)

		wg.Add(optimizer.confs.PopulationSize)

		for index := range agents {
			go func(agent *Agent) {
				defer wg.Done()
				agent.Run(optimizer.parent)
			}(&agents[index])
		}

		wg.Wait()

		foundBetter := false

		for _, agent := range agents {
			if agent.scheduler.powerConsumption < bestScheduler.powerConsumption {
				bestScheduler = clone(agent.scheduler)
				foundBetter = true
			}
		}

		if foundBetter {
			steadyState = 0
		} else {
			steadyState++
		}

	}

	log.Println("New best power consumption -> ", bestScheduler.powerConsumption, "W")

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
