package config

// Configuration holds information about the current state of the data center (workload, server status, etc.)
type Configuration struct {
	NumberOfTasks       int     `json:"numberOfTasks"`
	NumberOfServers     int     `json:"numberOfServers"`
	NumberOfTaskTypes   int     `json:"numberOfTaskTypes"`
	NumberOfServerTypes int     `json:"numberOfServerTypes"`
	Capped              float32 `json:"capped"`
}

// EvolutionaryAlgorithm holds the required parameters for the execution of the used evolutionary algorithm
type EvolutionaryAlgorithm struct {
	PopulationSize      int       `json:"populationSize"`
	NumberOfGenerations int       `json:"numberOfGenerations"`
	MutationRate        []float32 `json:"mutationRate"`
	MaxTime             int       `json:"maxTime"`
}
