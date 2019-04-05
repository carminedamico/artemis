package config

// EvolutionaryAlgorithmConfs holds the required parameters for the execution of the used evolutionary algorithm
type EvolutionaryAlgorithmConfs struct {
	PopulationSize      int       `json:"populationSize"`
	NumberOfGenerations int       `json:"numberOfGenerations"`
	MutationRate        []float32 `json:"mutationRate"`
	MaxTime             int       `json:"maxTime"`
}
