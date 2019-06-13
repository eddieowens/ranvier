package tests

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/server/app"
	"github.com/eddieowens/ranvier/server/app/service"
	"github.com/eddieowens/ranvier/server/tests/mocks"
	"github.com/stretchr/testify/suite"
	"path"
	"runtime"
)

type TestSuite struct {
	suite.Suite
	Injector   axon.Injector
	GitService *mocks.GitServiceMock
}

func (t *TestSuite) SetupSuite() {
	t.Injector = app.CreateInjector()
	t.Injector.Add(service.GitServiceKey, axon.StructPtr(new(mocks.GitServiceMock)))
	t.GitService = t.Injector.GetStructPtr(service.GitServiceKey).(*mocks.GitServiceMock)
}

func (t *TestSuite) Resources() string {
	_, filename, _, _ := runtime.Caller(0)
	d, _ := path.Split(filename)
	return path.Join(d, "resources")
}
