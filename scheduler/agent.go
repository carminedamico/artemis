package scheduler

import (
	"math/rand"
	"time"
)

// Agent represents an individual belonging to the population of an optimization algorithm
type Agent struct {
	scheduler Scheduler
	tabuList  TabuList
}

func newRandomAgent(parent Scheduler) *Agent {
	agent := &Agent{
		scheduler: clone(parent),
		tabuList:  *newTabuList(3),
	}

	for agent.scheduler.randomizer() == false {
	}

	agent.scheduler.getPowerConsumptionAccountingMigration(parent)

	return agent
}

// Run method allows an agent to start its optimization process, performing mutations based on the status of its tabu list
func (agent *Agent) Run(scheduler Scheduler) {
	op := agent.tabuList.get(5)

	switch op {
	case 0:
		oldPowerConsumption := agent.scheduler.powerConsumption
		agent.scheduler.TSWPMutation()
		agent.scheduler.getPowerConsumptionAccountingMigration(scheduler)
		newPowerConsumption := agent.scheduler.powerConsumption
		if newPowerConsumption > oldPowerConsumption {
			agent.tabuList.update(op)
		}

	case 1:
		oldPowerConsumption := agent.scheduler.powerConsumption
		agent.scheduler.TFFCMutation()
		agent.scheduler.getPowerConsumptionAccountingMigration(scheduler)
		newPowerConsumption := agent.scheduler.powerConsumption
		if newPowerConsumption > oldPowerConsumption {
			agent.tabuList.update(op)
		}

	case 2:
		oldPowerConsumption := agent.scheduler.powerConsumption
		agent.scheduler.TBFCMutation()
		agent.scheduler.getPowerConsumptionAccountingMigration(scheduler)
		newPowerConsumption := agent.scheduler.powerConsumption
		if newPowerConsumption > oldPowerConsumption {
			agent.tabuList.update(op)
		}

	case 3:
		oldPowerConsumption := agent.scheduler.powerConsumption
		agent.scheduler.SCMutation()
		agent.scheduler.getPowerConsumptionAccountingMigration(scheduler)
		newPowerConsumption := agent.scheduler.powerConsumption
		if newPowerConsumption > oldPowerConsumption {
			agent.tabuList.update(op)
		}

	case 4:
		oldPowerConsumption := agent.scheduler.powerConsumption
		agent.scheduler.SLRMutation()
		agent.scheduler.getPowerConsumptionAccountingMigration(scheduler)
		newPowerConsumption := agent.scheduler.powerConsumption
		if newPowerConsumption > oldPowerConsumption {
			agent.tabuList.update(op)
		}
	}
}

func (scheduler *Scheduler) randomizer() bool {
	rand.Seed(time.Now().UTC().UnixNano())

	for index := range scheduler.workload.Tasks {
		scheduler.workload.Tasks[index].AllocatedOn = -1
	}

	for index := range scheduler.datacenter.Servers {
		scheduler.datacenter.Servers[index].FreeCPU = scheduler.datacenter.Servers[index].CPU
		scheduler.datacenter.Servers[index].FreeRAM = scheduler.datacenter.Servers[index].RAM
	}

	rndTaskIndex := randInt(0, len(scheduler.workload.Tasks))

	for t := 0; t < len(scheduler.workload.Tasks); t++ {
		indexTask := (t + rndTaskIndex) % len(scheduler.workload.Tasks)

		scheduler.greedyMove(indexTask)
	}

	for _, task := range scheduler.workload.Tasks {
		if task.AllocatedOn == -1 {
			return false
		}
	}

	return true
}
