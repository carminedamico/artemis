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

	for index := randInt(0, len(scheduler.workload.Tasks)); index < len(scheduler.workload.Tasks); index = ((index + 1) % len(scheduler.workload.Tasks)) {
		targetTask := scheduler.workload.Tasks[index]

		indexToServer := targetTask.AllocatedOn
		toServer := scheduler.datacenter.Servers[indexToServer]

		if index != indexTaskToSwap && indexFromServer != indexToServer {
			if (toServer.FreeCPU-taskToSwap.CPU+targetTask.CPU) >= 0 && (fromServer.FreeCPU-targetTask.CPU+taskToSwap.CPU) >= 0 && (toServer.FreeRAM-taskToSwap.RAM+targetTask.RAM) >= 0 && (fromServer.FreeRAM-targetTask.RAM+taskToSwap.RAM) >= 0 {
				scheduler.migrateTask(indexTaskToSwap, indexToServer)
				scheduler.migrateTask(index, indexFromServer)

				fmt.Println(indexTaskToSwap, " -><- ", index)
				break
			}
		}
	}
}

// TFFCMutation performs a mutation
func (scheduler *Scheduler) TFFCMutation() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// TBFCMutation performs a mutation
func (scheduler *Scheduler) TBFCMutation() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// SCMutation performs a mutation
func (scheduler *Scheduler) SCMutation() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// SLRMutation performs a mutation
func (scheduler *Scheduler) SLRMutation() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func (scheduler *Scheduler) migrateTask(targetTask int, targetServer int) {
	scheduler.removeTaskFromServer(targetTask)

	scheduler.addTaskToServer(targetTask, targetServer)
}

func (scheduler *Scheduler) removeTaskFromServer(targetTask int) {
	fromServer := scheduler.workload.Tasks[targetTask].AllocatedOn

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
