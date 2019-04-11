package scheduler

import (
	"fmt"

	"github.com/carminedamico/artemis/config"
)

// Scheduler represents a scheduling strategy, given a datacenter and a workload to allocate
type Scheduler struct {
	datacenter config.Datacenter
	workload   config.Workload
}

// NewScheduler create a new scheduler starting from the current datacenter status
func NewScheduler(datacenter config.Datacenter, workload config.Workload) *Scheduler {
	scheduler := &Scheduler{
		datacenter: datacenter,
		workload:   workload,
	}
	scheduler.getDeltas()

	return scheduler
}

// Run method starts the scheduling process
func (scheduler *Scheduler) Run() {
	fmt.Printf("INITIAL POWER CONSUMPTION -> %f\n", scheduler.GetPowerConsumption())

	optimizer := NewOptimizer(scheduler)
	optimizer.Run()
}

// GetPowerConsumption returns the currnt amount of DC power consumed by the datacenter
func (scheduler *Scheduler) GetPowerConsumption() float32 {
	powerConsumption := float32(0)

	for _, server := range scheduler.datacenter.Servers {
		powerConsumption += (server.PowerDC*(float32(1)-server.IdleConsumption))*(float32(server.CPU-server.FreeCPU)/float32(server.CPU)) + (server.PowerDC * server.IdleConsumption)
	}

	return powerConsumption
}

// getDeltas calculates the current amount of free CPU and free RAM of each server of the datacenter
func (scheduler *Scheduler) getDeltas() {
	for index, server := range scheduler.datacenter.Servers {
		scheduler.datacenter.Servers[index].FreeCPU = server.CPU 
		scheduler.datacenter.Servers[index].FreeRAM = server.RAM
	}

	for _, task := range scheduler.workload.Tasks {
		scheduler.datacenter.Servers[task.AllocatedOn].FreeCPU -= task.CPU
		scheduler.datacenter.Servers[task.AllocatedOn].FreeRAM -= task.RAM
	}
}
