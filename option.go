package grpcp

type ServerOption struct {
	Port     int    `json:"port"`
	Listen   string `json:"listen"`
	TLS      bool   `json:"tls"`
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

type ClientOption struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Quiet      bool   `json:"quiet"`
	TLS        bool   `json:"tls"`
	SkipVerify bool   `json:"skip_verify"`
}
