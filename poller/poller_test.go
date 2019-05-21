package poller

import (
	"github.com/two-rabbits/ranvier/src/configuration"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"testing"
)

func TestName(t *testing.T) {
	repo, err := git.PlainClone("./something", false, &git.CloneOptions{
		URL:           "git@github.com:two-rabbits/ranvier.git",
		RemoteName:    remoteName,
		ReferenceName: plumbing.NewBranchReferenceName("testing_polling"),
	})

	if err == git.ErrRepositoryAlreadyExists {
		repo, err = git.PlainOpen("./something")
		if err != nil {
			t.Fatal(err)
		}
	} else if err != nil {
		t.Fatal(err)
	}

	gp := gitPollerImpl{
		Config: configuration.Config{
			CloneDirectory:  "something",
			GitPollInterval: 10,
		},
		repo: repo,
	}

	gp.fetchUpdates()

}
