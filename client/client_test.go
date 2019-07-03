package ranvier

import (
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
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
		Name: "stagingusers",
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

	// -- Then
	//
	if c.NoError(err) {
		c.NotNil(conn)
	}
}

func TestClientTest(t *testing.T) {
	c := new(ClientTest)
	c.IntegrationSuite = new(IntegrationSuite)
	suite.Run(t, c)
}
