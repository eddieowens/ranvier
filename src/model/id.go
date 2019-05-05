package model

type Id string

func (i Id) String() string {
	return string(i)
}
