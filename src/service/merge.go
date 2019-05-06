package service

import (
	"config-manager/src/collections"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/imdario/mergo"
	"reflect"
)

const MergeServiceKey = "MergeService"

type MergeService interface {
	MergeJsonMaps(dest, src *collections.JsonMap) collections.JsonMap
}

func overrideMerge(destination, source interface{}) interface{} {
	destType := reflect.TypeOf(destination)
	srcType := reflect.TypeOf(source)

	if destType.Kind() != srcType.Kind() {
		return source
	}

	if srcType.Kind() == reflect.Map || srcType.Kind() == reflect.Struct {
		if err := mergo.Merge(&destination, source, mergo.WithOverride); err != nil {
			fmt.Println(err)
			return destination
		}
		return destination
	} else {
		return source
	}
}

type mergeServiceImpl struct {
}

func (m *mergeServiceImpl) MergeJsonMaps(dest, src *collections.JsonMap) collections.JsonMap {
	destGabs, _ := gabs.ParseJSON(dest.GetRaw())
	srcGabs, _ := gabs.ParseJSON(src.GetRaw())

	err := destGabs.MergeFn(srcGabs, overrideMerge)
	if err != nil {
		fmt.Println(err)
	}

	return collections.NewJsonMap(destGabs.Bytes())
}
