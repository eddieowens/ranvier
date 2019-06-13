package tests

import (
	"encoding/json"
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/lang/injector"
	"github.com/stretchr/testify/suite"
	"path"
	"runtime"
)

type TestSuite struct {
	suite.Suite
	Injector axon.Injector
}

func (t *TestSuite) SetupSuite() {
	t.Injector = injector.CreateInjector()
}

func (t *TestSuite) Resources() string {
	_, filename, _, _ := runtime.Caller(0)
	d, _ := path.Split(filename)
	return path.Join(d, "resources")
}

func (t *TestSuite) EqualConfig(expected, actual json.RawMessage) {

}
