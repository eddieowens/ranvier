package ranvier

import (
	"fmt"
	"github.com/ory/dockertest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"os"
	"path"
)

type IntegrationSuite struct {
	suite.Suite
	Hostname string
	Port     string
	Host     string
	Resource *dockertest.Resource
	Pool     *dockertest.Pool
}

func (i *IntegrationSuite) TearDownSuite() {
	_ = i.Pool.Purge(i.Resource)
}

func (i *IntegrationSuite) SetupSuite() {
	var err error
	i.Pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Could not find the user's home: %s", err)
	}

	i.Resource, err = i.Pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "edwardrowens/ranvier-server",
		Tag:        "latest",
		Env:        []string{"RANVIER_GIT_SSHKEY=/.ssh/id_rsa"},
		Mounts:     []string{fmt.Sprintf("%s:/.ssh/id_rsa", path.Join(home, ".ssh", "id_rsa"))},
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	defer i.Resource.Expire(180)

	i.Port = i.Resource.GetPort("8080/tcp")
	i.Host = "localhost"
	i.Hostname = fmt.Sprintf("%s:%s", i.Host, i.Port)

	if err := i.Pool.Retry(func() error {
		resp, err := http.Get("http://" + i.Hostname + "/api/health")
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return errors.New("failed to start container")
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
}
