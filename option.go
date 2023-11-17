package grpcp

type ServerOption struct {
	Port   int    `json:"port"`
	Listen string `json:"listen"`
}

type ClientOption struct {
	Port int `json:"port"`
	Quiet bool `json:"quiet"`
}
