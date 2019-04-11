package config

type Task struct {
	CPU         float32 `json:"CPU"`
	RAM         float32 `json:"RAM"`
	AllocatedOn int `json:"allocatedOn"`
}

type Workload struct {
	Tasks []Task `json:"tasks"`
}
