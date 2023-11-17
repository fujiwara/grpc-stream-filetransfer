package filetransfer

type Option struct {
	Port int `json:"port"`
}

func NewDefaultOption() *Option {
	return &Option{
		Port: 5000,
	}
}
