package model

type Level int

const (
	Global Level = iota
	Cluster
	Namespace
	Application
)
