package service

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/imdario/mergo"
	"reflect"
)

const MergeServiceKey = "MergeService"

type MergeService interface {
	MergeJson(dest, src []byte) ([]byte, error)
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

func (m *mergeServiceImpl) MergeJson(dest, src []byte) ([]byte, error) {
	destGabs, _ := gabs.ParseJSON(dest)
	srcGabs, _ := gabs.ParseJSON(src)
	err := destGabs.MergeFn(srcGabs, overrideMerge)
	if err != nil {
		return nil, err
	}

	return destGabs.Bytes(), nil
}
