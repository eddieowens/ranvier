package poller

import (
	"github.com/two-rabbits/ranvier/src/configuration"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"testing"
)

func TestName(t *testing.T) {
	repo, err := git.PlainClone(g.Config.CloneDirectory, false, &git.CloneOptions{
		URL:           "git@github.com:two-rabbits/ranvier.git",
		RemoteName:    remoteName,
		ReferenceName: plumbing.NewRemoteReferenceName("git@github.com:two-rabbits/ranvier.git", "testing_polling"),
	})
	gp := gitPollerImpl{
		Config: configuration.Config{
			CloneDirectory:  "something",
			GitPollInterval: 10,
		},
		repo: repo,
	}

	gp.fetchUpdates()

}
