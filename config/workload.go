package config

type Task struct {
	CPU         int `json:"CPU"`
	RAM         int `json:"RAM"`
	AllocatedOn int `json:"AllocatedOn"`
}

type Workload struct {
	Tasks []Task `json:"Tasks"`
}
