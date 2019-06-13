package controller

import (
	"github.com/eddieowens/ranvier/server/app/collections"
	"github.com/eddieowens/ranvier/server/app/controller"
	"github.com/eddieowens/ranvier/server/app/exchange/response"
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/eddieowens/ranvier/server/app/state"
	"github.com/eddieowens/ranvier/server/tests/integration"
	json "github.com/json-iterator/go"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type ConfigControllerTest struct {
	integration.Integration
	ConfigController controller.ConfigController
	ConfigMap        collections.ConfigMap
}

func (s *ConfigControllerTest) SetupTest() {
	s.ConfigController = s.Injector.GetStructPtr(controller.ConfigControllerKey).(controller.ConfigController)
	s.ConfigMap = s.Injector.GetStructPtr(state.ConfigMapKey).(collections.ConfigMap)
}

func (s *ConfigControllerTest) TestQuery() {
	// -- Given
	//
	name := "StagingUsers"
	ctx, rec := s.Request(integration.Request{
		PathParams: map[string]string{
			"name": name,
		},
		QueryParams: map[string]string{
			"query": "$.db",
		},
	})

	s.ConfigMap.Set(model.Config{
		Name: name,
		Data: map[string]interface{}{"db": "pg"},
	})

	expected := response.Config{
		Data: &response.ConfigData{
			Name:   name,
			Config: "pg",
		},
	}

	// -- When
	//
	err := s.ConfigController.Query(ctx)

	// -- Then
	//
	if s.NoError(err) {
		s.Equal(http.StatusOK, rec.Code)
		var actual response.Config
		err := json.Unmarshal(rec.Body.Bytes(), &actual)
		if s.NoError(err) {
			s.Equal(expected, actual)
		}
	}
}

func TestServerTest(t *testing.T) {
	suite.Run(t, new(ConfigControllerTest))
}
