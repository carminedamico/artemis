package scheduler

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// TSWPMutation randomly selects a task from the workload and searches among the others one task
// that can be exchanged with.
func (scheduler *Scheduler) TSWPMutation() {
	rand.Seed(time.Now().UTC().UnixNano())

	indexTaskToSwap := randInt(0, len(scheduler.workload.Tasks))
	taskToSwap := scheduler.workload.Tasks[indexTaskToSwap]

	indexFromServer := taskToSwap.AllocatedOn
	fromServer := scheduler.datacenter.Servers[indexFromServer]

	startingIndex := randInt(0, len(scheduler.workload.Tasks))

	for i := range scheduler.workload.Tasks {
		index := (i + startingIndex) % len(scheduler.workload.Tasks)

		targetTask := scheduler.workload.Tasks[index]

		indexToServer := targetTask.AllocatedOn
		toServer := scheduler.datacenter.Servers[indexToServer]

		if index != indexTaskToSwap && indexFromServer != indexToServer {

			firstServerCPUok := (fromServer.FreeCPU - targetTask.CPU + taskToSwap.CPU) >= 0
			firstServerRAMok := (fromServer.FreeRAM - targetTask.RAM + taskToSwap.RAM) >= 0
			secondServerCPUok := (toServer.FreeCPU - taskToSwap.CPU + targetTask.CPU) >= 0
			secondServerRAMok := (toServer.FreeRAM - taskToSwap.RAM + targetTask.CPU) >= 0

			if firstServerCPUok && firstServerRAMok && secondServerCPUok && secondServerRAMok {

				scheduler.migrateTask(indexTaskToSwap, indexToServer)
				scheduler.migrateTask(index, indexFromServer)

				break
			}
		}

	}
}

// TFFCMutation randomly selects one task that can be moved on another server
func (scheduler *Scheduler) TFFCMutation() {
	rand.Seed(time.Now().UTC().UnixNano())

	indexTaskToSwap := randInt(0, len(scheduler.workload.Tasks))
	taskToSwap := scheduler.workload.Tasks[indexTaskToSwap]

	startingIndex := randInt(0, len(scheduler.datacenter.Servers))

	for i := range scheduler.datacenter.Servers {
		serverIndex := (i + startingIndex) % len(scheduler.datacenter.Servers)

		if serverIndex != taskToSwap.AllocatedOn {
			toServer := scheduler.datacenter.Servers[serverIndex]
			if toServer.FreeCPU >= taskToSwap.CPU && toServer.FreeRAM >= taskToSwap.RAM {
				scheduler.migrateTask(indexTaskToSwap, serverIndex)

				break
			}
		}
	}
}

// TBFCMutation randomly selects one task that can be moved on the best possible node,
// that is represented by the one with the largest unused amount of resources.
func (scheduler *Scheduler) TBFCMutation() {
	rand.Seed(time.Now().UTC().UnixNano())

	indexTaskToSwap := randInt(0, len(scheduler.workload.Tasks))
	taskToSwap := scheduler.workload.Tasks[indexTaskToSwap]

	bestServerIndex := -1
	maxFreeCPU := float32(0)
	maxFreeRAM := float32(0)

	startingIndex := randInt(0, len(scheduler.datacenter.Servers))

	for i := range scheduler.datacenter.Servers {
		serverIndex := (i + startingIndex) % len(scheduler.datacenter.Servers)
		if serverIndex != taskToSwap.AllocatedOn {
			toServer := scheduler.datacenter.Servers[serverIndex]
			if toServer.FreeCPU >= taskToSwap.CPU && toServer.FreeRAM >= taskToSwap.RAM && toServer.FreeCPU > maxFreeCPU && toServer.FreeRAM >= maxFreeRAM {
				bestServerIndex = serverIndex
			}
		}
	}

	if bestServerIndex != -1 {
		scheduler.migrateTask(indexTaskToSwap, bestServerIndex)
	}

}

// SCMutation randomly selects one server from the server list and tries to saturate its
// available resources, moving tasks to it.
func (scheduler *Scheduler) SCMutation() {
	rand.Seed(time.Now().UTC().UnixNano())

	indexServerToConsolidate := randInt(0, len(scheduler.datacenter.Servers))
	

	startingIndex := randInt(0, len(scheduler.workload.Tasks))

	for i := range scheduler.workload.Tasks {
		taskIndex := (i + startingIndex) % len(scheduler.workload.Tasks)

		task := scheduler.workload.Tasks[taskIndex]
		serverToConsolidate := scheduler.datacenter.Servers[indexServerToConsolidate]

		if task.AllocatedOn != indexServerToConsolidate && serverToConsolidate.FreeCPU >= task.CPU && serverToConsolidate.FreeRAM >= task.RAM {
			scheduler.migrateTask(taskIndex, indexServerToConsolidate)
		}

	}

}

// SLRMutation randomly selects one server from the server list and tries to redistribute
// its whole load on other servers.
func (scheduler *Scheduler) SLRMutation() {
	rand.Seed(time.Now().UTC().UnixNano())

	indexServerToEmpty := randInt(0, len(scheduler.datacenter.Servers))

	startingIndex := randInt(0, len(scheduler.workload.Tasks))

	for i := range scheduler.workload.Tasks {
		indexTask := (i + startingIndex) % len(scheduler.workload.Tasks)
		task := scheduler.workload.Tasks[indexTask]

		if task.AllocatedOn == indexServerToEmpty {
			for indexServer, server := range scheduler.datacenter.Servers {
				if server.FreeCPU >= task.CPU && server.FreeRAM >= task.RAM {

					scheduler.migrateTask(indexTask, indexServer)
				}
			}
		}
	}
}

func (scheduler *Scheduler) greedyMove(taskIndex int) {
	rand.Seed(time.Now().UTC().UnixNano())

	task := scheduler.workload.Tasks[taskIndex]

	startingIndex := randInt(0, len(scheduler.datacenter.Servers))

	bestServer := -1
	bestAttractiveness := float32(math.MaxFloat32)

	for i := range scheduler.datacenter.Servers {
		serverIndex := (i + startingIndex) % len(scheduler.datacenter.Servers)
		server := scheduler.datacenter.Servers[serverIndex]

		if server.FreeCPU >= task.CPU && server.FreeRAM >= task.RAM {
			serverPowerConsumption := (server.PowerDC*(float32(1)-server.IdleConsumption))*((server.CPU-server.FreeCPU-task.CPU)/server.CPU) + (server.PowerDC * server.IdleConsumption)
			if serverPowerConsumption < bestAttractiveness {
				bestServer = serverIndex
				bestAttractiveness = serverPowerConsumption
			}
		}

	}

	scheduler.addTaskToServer(taskIndex, bestServer)
}

func (scheduler *Scheduler) migrateTask(targetTask int, targetServer int) {
	scheduler.removeTaskFromServer(targetTask)

	scheduler.addTaskToServer(targetTask, targetServer)

}

func (scheduler *Scheduler) removeTaskFromServer(targetTask int) {
	fromServer := scheduler.workload.Tasks[targetTask].AllocatedOn

	if targetTask >= len(scheduler.workload.Tasks) || fromServer >= len(scheduler.datacenter.Servers) {
		fmt.Println(targetTask, " ", fromServer)
	}

	scheduler.datacenter.Servers[fromServer].FreeCPU += scheduler.workload.Tasks[targetTask].CPU
	scheduler.datacenter.Servers[fromServer].FreeRAM += scheduler.workload.Tasks[targetTask].RAM

	scheduler.workload.Tasks[targetTask].AllocatedOn = -1
}

func (scheduler *Scheduler) addTaskToServer(targetTask int, targetServer int) {
	scheduler.datacenter.Servers[targetServer].FreeCPU -= scheduler.workload.Tasks[targetTask].CPU
	scheduler.datacenter.Servers[targetServer].FreeRAM -= scheduler.workload.Tasks[targetTask].RAM

	scheduler.workload.Tasks[targetTask].AllocatedOn = targetServer
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
