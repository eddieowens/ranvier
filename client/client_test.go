package ranvier

import (
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type ClientTest struct {
	*IntegrationSuite
	client Client
}

func (c *ClientTest) SetupTest() {
	client, err := NewClient(&ClientOptions{
		Hostname:        "localhost:8080",
		ConfigDirectory: os.TempDir(),
	})

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

	c.client = client
}

func (c *ClientTest) TestQuery() {
	// -- Given
	//
	expectedData := map[string]interface{}{
		"db": map[string]interface{}{
			"conns": float64(10),
			"pool":  float64(13),
			"url":   "username@pg.staging.mycompany.com:5432",
		},
	}

	expected := &model.Config{
		Name: "users-staging",
		Data: expectedData,
	}

	// -- When
	//
	conf, err := c.client.Query(&QueryOptions{
		IgnoreCache: false,
		Name:        expected.Name,
		Query:       "$",
	})

	// -- Then
	//
	if c.NoError(err) {
		c.Equal(expected, conf)
	}
}

func (c *ClientTest) TestConnect() {
	// -- When
	//
	conn, err := c.client.Connect(&ConnOptions{
		Names: []string{"stagingusers"},
	})
	defer c.client.Disconnect(conn)

	// -- Then
	//
	if c.NoError(err) {
		c.NotNil(conn)
	}
}

func (c *ClientTest) TestDisconnect() {
	// -- Given
	conn, err := c.client.Connect(&ConnOptions{
		Names: []string{"stagingusers"},
	})

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

	// -- When
	//
	c.client.Disconnect(conn)

	// -- Then
	//
}

func (c *ClientTest) TestConfigCreate() {
	// -- Given
	//
	expectedName := "expectedCreate"

	conn, err := c.client.Connect(&ConnOptions{
		Names: []string{expectedName},
	})
	defer c.client.Disconnect(conn)

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

	expected := model.ConfigEvent{
		Config: model.Config{
			Name: expectedName,
			Data: "some data",
		},
		EventType: model.EventTypeCreate,
	}

	go func() {
		time.Sleep(time.Second * 1)
		_, err := c.client.(*clientImpl).Update(&expected.Config)

		if !c.NoError(err) {
			panic(err)
		}
	}()

	// -- When
	//
	actual := <-conn.OnUpdate

	// -- Then
	//
	c.Equal(expected, actual)
}

func (c *ClientTest) TestConfigUpdate() {
	// -- Given
	//
	expectedName := "expectedUpdate"

	_, err := c.client.(*clientImpl).Update(&model.Config{
		Name: expectedName,
		Data: "some data",
	})

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

	conn, err := c.client.Connect(&ConnOptions{
		Names: []string{expectedName},
	})
	defer c.client.Disconnect(conn)

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

	expected := model.ConfigEvent{
		Config: model.Config{
			Name: expectedName,
			Data: "new data",
		},
		EventType: model.EventTypeUpdate,
	}

	go func() {
		time.Sleep(time.Second * 1)
		_, err := c.client.(*clientImpl).Update(&expected.Config)

		if !c.NoError(err) {
			panic(err)
		}
	}()

	// -- When
	//
	actual := <-conn.OnUpdate

	// -- Then
	//
	c.Equal(expected, actual)
}

func (c *ClientTest) TestConfigDelete() {
	// -- Given
	//
	expectedName := "expectedDelete"

	_, err := c.client.(*clientImpl).Update(&model.Config{
		Name: expectedName,
		Data: "some data",
	})

	conn, err := c.client.Connect(&ConnOptions{
		Names: []string{expectedName},
	})
	defer c.client.Disconnect(conn)

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

	expected := model.ConfigEvent{
		Config: model.Config{
			Name: expectedName,
			Data: "some data",
		},
		EventType: model.EventTypeDelete,
	}

	go func() {
		time.Sleep(time.Second * 1)
		_, err := c.client.(*clientImpl).Delete(expectedName)

		if !c.NoError(err) {
			panic(err)
		}
	}()

	// -- When
	//
	actual := <-conn.OnUpdate

	// -- Then
	//
	c.Equal(expected, actual)
}

func TestClientTest(t *testing.T) {
	c := new(ClientTest)
	c.IntegrationSuite = new(IntegrationSuite)
	suite.Run(t, c)
}
