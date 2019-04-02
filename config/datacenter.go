package config

type Server struct {
	CPU     int `json:"CPU"`
	RAM     int `json:"RAM"`
	PowerDC int `json:"PowerDC"`
}

type Datacenter struct {
	Nodes []Server `json:"Nodes"`
}
