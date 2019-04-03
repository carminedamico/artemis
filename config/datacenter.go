package config

type Server struct {
	CPU     int `json:"CPU"`
	RAM     int `json:"RAM"`
	PowerDC int `json:"powerDC"`
}

type Datacenter struct {
	Servers []Server `json:"servers"`
}
