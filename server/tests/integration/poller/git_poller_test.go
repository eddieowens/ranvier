package poller

import (
	"github.com/eddieowens/ranvier/server/tests/integration"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GitPollerTest struct {
	integration.Integration
}

func (g *GitPollerTest) SetupTest() {
}

func (g *GitPollerTest) TestPoll() {
	// -- Given
	//

	// -- When
	//

	// -- Then
	//
}

func TestGitPollerTest(t *testing.T) {
	suite.Run(t, new(GitPollerTest))
}
