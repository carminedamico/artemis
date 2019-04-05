package scheduler

import (
	"github.com/carminedamico/artemis/config"
)

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
	return scheduler
}

// Run method starts the scheduling process
func (scheduler *Scheduler) Run() {
	optimizer := NewOptimizer(scheduler)
	optimizer.Run()
}

func (scheduler *Scheduler) GetPowerConsumption() float32 {
	var freeCPU = make([]int, len(scheduler.datacenter.Servers))

	for _, task := range scheduler.workload.Tasks {
		freeCPU[task.AllocatedOn] += task.CPU
	}

	powerConsumption := float32(0)

	for index, server := range scheduler.datacenter.Servers {
		freeCPU[index] = server.CPU - freeCPU[index]
		powerConsumption += (server.PowerDC*(float32(1)-server.IdleConsumption))*(float32(server.CPU-freeCPU[index])/float32(server.CPU)) + (server.PowerDC * server.IdleConsumption)
	}

	return powerConsumption
}
