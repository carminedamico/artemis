package config

type Server struct {
	CPU             int     `json:"CPU"`
	RAM             int     `json:"RAM"`
	PowerDC         float32 `json:"powerDC"`
	IdleConsumption float32 `json:"idleConsumption"`
	Capping         float32 `json:"capping"`
}

type Datacenter struct {
	Servers []Server `json:"servers"`
}
