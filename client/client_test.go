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
}

func (c *ClientTest) TestQuery() {
	// -- Given
	//
	client, err := NewClient(&ClientOptions{
		Hostname:        c.Hostname,
		ConfigDirectory: os.TempDir(),
	})

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

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
	conf, err := client.Query(&QueryOptions{
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
	// -- Given
	//
	client, err := NewClient(&ClientOptions{
		Hostname:        c.Hostname,
		ConfigDirectory: os.TempDir(),
	})

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

	// -- When
	//
	conn, err := client.Connect(&ConnOptions{
		Names: []string{"stagingusers"},
	})
	defer client.Disconnect(conn)

	// -- Then
	//
	if c.NoError(err) {
		c.NotNil(conn)
	}
}

func (c *ClientTest) TestDisconnect() {
	// -- Given
	//
	client, err := NewClient(&ClientOptions{
		Hostname:        c.Hostname,
		ConfigDirectory: os.TempDir(),
	})

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

	conn, err := client.Connect(&ConnOptions{
		Names: []string{"stagingusers"},
	})

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

	// -- When
	//
	client.Disconnect(conn)

	// -- Then
	//
}

func (c *ClientTest) TestConfigCreate() {
	// -- Given
	//
	client, err := NewClient(&ClientOptions{
		Hostname:        c.Hostname,
		ConfigDirectory: os.TempDir(),
	})

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

	expectedName := "expected"

	conn, err := client.Connect(&ConnOptions{
		Names: []string{expectedName},
	})
	defer client.Disconnect(conn)

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
		_, err := client.(*clientImpl).Update(&expected.Config)

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
	client, err := NewClient(&ClientOptions{
		Hostname:        c.Hostname,
		ConfigDirectory: os.TempDir(),
	})

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

	expectedName := "expected"

	_, err = client.(*clientImpl).Update(&model.Config{
		Name: expectedName,
		Data: "some data",
	})

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

	conn, err := client.Connect(&ConnOptions{
		Names: []string{expectedName},
	})
	defer client.Disconnect(conn)

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
		_, err := client.(*clientImpl).Update(&expected.Config)

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
	client, err := NewClient(&ClientOptions{
		Hostname:        c.Hostname,
		ConfigDirectory: os.TempDir(),
	})

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

	expectedName := "expected"

	_, err = client.(*clientImpl).Update(&model.Config{
		Name: expectedName,
		Data: "some data",
	})

	conn, err := client.Connect(&ConnOptions{
		Names: []string{expectedName},
	})
	defer client.Disconnect(conn)

	if !c.NoError(err) {
		c.FailNow(err.Error())
	}

	expected := model.ConfigEvent{
		Config: model.Config{
			Name: expectedName,
		},
		EventType: model.EventTypeDelete,
	}

	go func() {
		time.Sleep(time.Second * 1)
		_, err := client.(*clientImpl).Update(&expected.Config)

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
