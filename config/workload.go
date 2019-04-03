package config

type Task struct {
	CPU         int `json:"CPU"`
	RAM         int `json:"RAM"`
	AllocatedOn int `json:"allocatedOn"`
}

type Workload struct {
	Tasks []Task `json:"tasks"`
}
