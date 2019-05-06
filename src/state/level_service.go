package state

import (
	"config-manager/src/model"
	"strings"
)

const LevelServiceKey = "LevelService"

type LevelService interface {
	FromString(st string) model.Level
	ToString(level model.Level) string
}

type levelServiceImpl struct {
}

func (l *levelServiceImpl) ToString(level model.Level) string {
	switch level {
	case model.Global:
		return "global"
	case model.Cluster:
		return "cluster"
	case model.Namespace:
		return "namespace"
	case model.Application:
		return "application"
	default:
		return ""
	}
}

func (l *levelServiceImpl) FromString(st string) model.Level {
	switch strings.ToLower(st) {
	case "global":
		return model.Global
	case "cluster":
		return model.Cluster
	case "namespace":
		return model.Namespace
	case "application":
		return model.Application
	default:
		return -1
	}
}
