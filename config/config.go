package config

// Configuration holds information about the current state of the data center (workload, server status, etc.)
type Configuration struct {
	NumberOfTasks       int     `json:"NumberOfTasks"`
	NumberOfServers     int     `json:"NumberOfServers"`
	NumberOfTaskTypes   int     `json:"NumberOfTaskTypes"`
	NumberOfServerTypes int     `json:"NumberOfServerTypes"`
	Capacity            float32 `json:"Capacity"`
}

// EvolutionaryAlgorithm holds the required parameters for the execution of the used evolutionary algorithm
type EvolutionaryAlgorithm struct {
	PopulationSize      int       `json:"PopulationSize"`
	NumberOfGenerations int       `json:"NumberOfGenerations"`
	MutationRate        []float32 `json:"MutationRate"`
	MaxTime             int       `json:"MaxTime"`
}
