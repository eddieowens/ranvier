package state

import (
	"github.com/two-rabbits/ranvier/src/model"
	"strconv"
	"strings"
)

const IdServiceKey = "IdService"

const IdSeparator = "~"

type IdService interface {
	Name(id model.Id) string
	Names(id model.Id) []string
	Id(names ...string) model.Id
	GlobalId() model.Id
	ClusterId(name string) model.Id
	NamespaceId(name, cluster string) model.Id
	ApplicationId(name, namespace, cluster string) model.Id
	VersionedId(id model.Id, version int) model.Id
	IsVersionedId(id model.Id) bool
}

type idServiceImpl struct {
}

func (i *idServiceImpl) IsVersionedId(id model.Id) bool {
	n := strings.Split(id.String(), IdSeparator)
	if len(n) <= 0 {
		return false
	}
	if _, err := strconv.Atoi(n[len(n)-1]); err != nil {
		return false
	}
	return true
}

func (i *idServiceImpl) Id(names ...string) model.Id {
	return model.Id(strings.Join(names, IdSeparator))
}

func (i *idServiceImpl) Names(id model.Id) []string {
	return strings.Split(id.String(), IdSeparator)
}

func (i *idServiceImpl) Name(id model.Id) string {
	meta := strings.Split(id.String(), IdSeparator)
	return meta[len(meta)-1]
}

func (i *idServiceImpl) GlobalId() model.Id {
	return GlobalId
}

func (i *idServiceImpl) ClusterId(name string) model.Id {
	return i.Id(name)
}

func (i *idServiceImpl) NamespaceId(name, cluster string) model.Id {
	return i.Id(cluster, name)
}

func (i *idServiceImpl) ApplicationId(name, namespace, cluster string) model.Id {
	return i.Id(cluster, namespace, name)
}

func (i *idServiceImpl) VersionedId(id model.Id, version int) model.Id {
	ver := strconv.Itoa(version)
	idSlice := append([]string{id.String()}, ver)
	return model.Id(strings.Join(idSlice, IdSeparator))
}

func (i *idServiceImpl) ToName(id model.Id) string {
	names := strings.Split(string(id), IdSeparator)
	if _, err := strconv.Atoi(names[len(names)-1]); err != nil {
		return names[len(names)-1]
	} else {
		return names[len(names)-2]
	}
}
