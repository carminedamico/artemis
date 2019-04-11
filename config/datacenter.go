package config

type Server struct {
	CPU             float32 `json:"CPU"`
	RAM             float32 `json:"RAM"`
	PowerDC         float32 `json:"powerDC"`
	IdleConsumption float32 `json:"idleConsumption"`
	Capping         float32 `json:"capping"`
	FreeCPU         float32
	FreeRAM         float32
}

type Datacenter struct {
	Servers []Server `json:"servers"`
}
