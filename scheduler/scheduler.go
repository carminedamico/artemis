package scheduler

import (
	"github.com/carminedamico/artemis/config"
)

// Scheduler represents a scheduling strategy, given a datacenter and a workload to allocate
type Scheduler struct {
	datacenter       config.Datacenter
	workload         config.Workload
	powerConsumption float32
}

// NewScheduler create a new scheduler starting from the current datacenter status
func NewScheduler(datacenter config.Datacenter, workload config.Workload) *Scheduler {
	scheduler := &Scheduler{
		datacenter: datacenter,
		workload:   workload,
	}
	scheduler.getDeltas()
	scheduler.GetPowerConsumption()

	return scheduler
}

// Run method starts the scheduling process
func (scheduler *Scheduler) Run() {
	optimizer := NewOptimizer(scheduler)
	optimizer.Run()
}

// GetPowerConsumption returns the currnt amount of DC power consumed by the datacenter
func (scheduler *Scheduler) GetPowerConsumption() {
	powerConsumption := float32(0)

	for _, server := range scheduler.datacenter.Servers {
		if server.CPU != server.FreeCPU {
			powerConsumption += (server.PowerDC*(float32(1)-server.IdleConsumption))*(float32(server.CPU-server.FreeCPU)/float32(server.CPU)) + (server.PowerDC * server.IdleConsumption)
		}
	}

	scheduler.powerConsumption = powerConsumption
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
