package scheduler

import (
	"fmt"
	"math/rand"
	"time"
)

// TSWPMutation performs a mutation
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

			if (toServer.FreeCPU-taskToSwap.CPU+targetTask.CPU) >= 0 && (fromServer.FreeCPU-targetTask.CPU+taskToSwap.CPU) >= 0 && (toServer.FreeRAM-taskToSwap.RAM+targetTask.RAM) >= 0 && (fromServer.FreeRAM-targetTask.RAM+taskToSwap.RAM) >= 0 {

				scheduler.migrateTask(indexTaskToSwap, indexToServer)
				scheduler.migrateTask(index, indexFromServer)

				break
			}
		}

	}
}

// TFFCMutation performs a mutation
func (scheduler *Scheduler) TFFCMutation() {
	rand.Seed(time.Now().UTC().UnixNano())

	indexTaskToSwap := randInt(0, len(scheduler.workload.Tasks))
	taskToSwap := scheduler.workload.Tasks[indexTaskToSwap]

	for serverIndex := range scheduler.datacenter.Servers {
		if serverIndex != taskToSwap.AllocatedOn {

			toServer := scheduler.datacenter.Servers[serverIndex]
			if (toServer.FreeCPU-taskToSwap.CPU) >= 0 && (toServer.FreeRAM-taskToSwap.RAM) >= 0 {
				scheduler.migrateTask(indexTaskToSwap, serverIndex)

				break
			}
		}
	}
}

// TBFCMutation performs a mutation
func (scheduler *Scheduler) TBFCMutation() {
	rand.Seed(time.Now().UTC().UnixNano())

	indexTaskToSwap := randInt(0, len(scheduler.workload.Tasks))
	taskToSwap := scheduler.workload.Tasks[indexTaskToSwap]
	
	bestServerIndex := -1
	maxFreeCPU := float32(0)
	maxFreeRAM := float32(0)

	for serverIndex := range scheduler.datacenter.Servers {
		if serverIndex != taskToSwap.AllocatedOn {

			toServer := scheduler.datacenter.Servers[serverIndex]

			if (toServer.FreeCPU-taskToSwap.CPU) >= 0 && (toServer.FreeRAM-taskToSwap.RAM) >= 0 && (toServer.FreeCPU > maxFreeCPU) && (toServer.FreeRAM >= maxFreeRAM) {
				bestServerIndex = serverIndex
			}
		}
	}

	if bestServerIndex != -1 {
		scheduler.migrateTask(indexTaskToSwap, bestServerIndex)
	}

}

// SCMutation performs a mutation
func (scheduler *Scheduler) SCMutation() {
	rand.Seed(time.Now().UTC().UnixNano())

	indexServerToConsolidate := randInt(0, len(scheduler.datacenter.Servers))
	serverToConsolidate := scheduler.datacenter.Servers[indexServerToConsolidate]

	startingIndex := randInt(0, len(scheduler.workload.Tasks))

	for i, task := range scheduler.workload.Tasks {
		taskIndex := (i + startingIndex) % len(scheduler.workload.Tasks)

		if task.AllocatedOn != indexServerToConsolidate {
			if (serverToConsolidate.FreeCPU-task.CPU) >= 0 && (serverToConsolidate.FreeRAM-task.RAM) >= 0 {

				scheduler.migrateTask(taskIndex, indexServerToConsolidate)
			}
		}
	}

}

// SLRMutation performs a mutation
func (scheduler *Scheduler) SLRMutation() {
	rand.Seed(time.Now().UTC().UnixNano())

	indexServerToEmpty := randInt(0, len(scheduler.datacenter.Servers))

	for indexTask, task := range scheduler.workload.Tasks {
		if task.AllocatedOn == indexServerToEmpty {
			for indexServer, server := range scheduler.datacenter.Servers {
				if server.FreeCPU >= task.CPU && server.FreeRAM >= task.RAM {

					scheduler.migrateTask(indexTask, indexServer)
				}
			}
		}
	}

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
