package model

type Config struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

func (l Config) Copy() interface{} {
	return Config{
		Data: l.Data,
		Name: l.Name,
	}
}
