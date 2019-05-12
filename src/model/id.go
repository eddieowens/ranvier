package model

type Id string

func (i Id) String() string {
	return string(i)
}

type IdNames struct {
	Cluster     string
	Namespace   string
	Application string
}
