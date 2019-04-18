package config

// EvolutionaryAlgorithmConfs holds the required parameters for the execution of the used evolutionary algorithm
type EvolutionaryAlgorithmConfs struct {
	PopulationSize      int    `json:"populationSize"`
	EliteSize           int    `json:"eliteSize"`
	NumberOfGenerations int    `json:"numberOfGenerations"`
	MaxTime             int    `json:"maxTime"`
	LogFile             string `json:"logFile"`
}
