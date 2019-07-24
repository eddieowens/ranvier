package loadtests

import (
	"fmt"
	ranvier "github.com/eddieowens/ranvier/client"
	"github.com/stretchr/testify/suite"
	"testing"
)

const TargetHost = "34.94.19.66"

type LoadTest struct {
	suite.Suite
}

func (l *LoadTest) SetupTest() {
}

func (l *LoadTest) TestWebsocketConnections() {
	// -- Given
	//
	numConns := 500
	c, err := ranvier.NewClient(&ranvier.ClientOptions{
		Hostname: TargetHost,
	})
	if !l.NoError(err) {
		l.FailNow(err.Error())
	}

	errChan := make(chan bool)

	// -- When
	//
	for i := 0; i < numConns; i++ {
		go func(i int) {
			fmt.Println(fmt.Sprintf("%d starting", i))
			conn, err := c.Connect(&ranvier.ConnOptions{
				Names: []string{"staging-users"},
			})
			if err != nil {
				fmt.Println("err!", err.Error())
				errChan <- true
			}
			for {
				select {
				case event := <-conn.OnUpdate:
					fmt.Println(fmt.Sprintf("%d event! %v", i, event))
				case <-errChan:
					fmt.Println(fmt.Sprintf("closing %d!", i))
				}
			}
		}(i)
	}

	// -- Then
	//
	for {
		<-errChan
		fmt.Println("killing everything!")
	}
}

func TestLoadTest(t *testing.T) {
	suite.Run(t, new(LoadTest))
}
