package commons

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/imdario/mergo"
	"reflect"
)

type JsonMerger interface {
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

type jsonMergerImpl struct {
}

func NewJsonMerger() JsonMerger {
	return &jsonMergerImpl{}
}

func (m *jsonMergerImpl) MergeJson(dest, src []byte) ([]byte, error) {
	destGabs, _ := gabs.ParseJSON(dest)
	srcGabs, _ := gabs.ParseJSON(src)
	err := destGabs.MergeFn(srcGabs, overrideMerge)
	if err != nil {
		return nil, err
	}

	return destGabs.Bytes(), nil
}
