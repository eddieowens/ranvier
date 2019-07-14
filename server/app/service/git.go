package service

import (
	"errors"
	"fmt"
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
)

const GitServiceKey = "GitService"

const remoteName = "origin"

type GitService interface {
	Clone(remote, branch, directory string) (*git.Repository, error)
	DiffRemote(repo *git.Repository, branch string) ([]model.GitChange, error)
	FetchLatestRemoteCommit(repo *git.Repository, branch string) (*object.Commit, error)
}

type gitServiceImpl struct {
	AuthMethod transport.AuthMethod `inject:"AuthMethod"`
}

func (g *gitServiceImpl) DiffRemote(repo *git.Repository, branch string) ([]model.GitChange, error) {
	err := repo.Fetch(&git.FetchOptions{
		Auth: g.AuthMethod,
	})
	if err != nil {
		return nil, err
	}

	h, err := repo.Head()
	if err != nil {
		return nil, err
	}

	remCommit, err := g.FetchLatestRemoteCommit(repo, branch)
	if err != nil {
		return nil, err
	}

	currentCommit, err := repo.CommitObject(h.Hash())
	if err != nil {
		return nil, err
	}

	originTree, err := remCommit.Tree()
	if err != nil {
		return nil, err
	}
	branchTree, err := currentCommit.Tree()
	if err != nil {
		return nil, err
	}

	diffs, err := branchTree.Diff(originTree)
	if err != nil {
		return nil, err
	}

	changes := make([]model.GitChange, 0)
	for _, d := range diffs {
		a, err := d.Action()
		if err != nil {
			return nil, err
		}

		var gitChange model.GitChange
		switch a {
		case merkletrie.Modify:
			gitChange = model.GitChange{
				Filename:  d.To.Name,
				EventType: model.EventTypeUpdate,
			}
		case merkletrie.Delete:
			gitChange = model.GitChange{
				Filename:  d.From.Name,
				EventType: model.EventTypeDelete,
			}
		case merkletrie.Insert:
			gitChange = model.GitChange{
				Filename:  d.To.Name,
				EventType: model.EventTypeCreate,
			}
		}

		changes = append(changes, gitChange)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	err = wt.Pull(&git.PullOptions{
		SingleBranch: true,
		Auth:         g.AuthMethod,
	})

	if err != nil {
		return nil, err
	}

	return changes, nil
}

func (g *gitServiceImpl) Clone(remote, branch, directory string) (*git.Repository, error) {
	logrus.WithField("remote", remote).
		WithField("branch", branch).
		WithField("clone_directory", directory).
		Debug("Cloning repo")
	repo, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:           remote,
		RemoteName:    remoteName,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		Auth:          g.AuthMethod,
	})

	if err == git.ErrRepositoryAlreadyExists {
		logrus.Debug("Failed to clone as repo is already present")
		return git.PlainOpen(directory)
	} else if err != nil {
		return nil, err
	}

	return repo, nil
}

func (g *gitServiceImpl) FetchLatestRemoteCommit(repo *git.Repository, branch string) (*object.Commit, error) {
	rem, err := repo.Remote(remoteName)
	if err != nil {
		return nil, err
	}

	rfs, err := rem.List(&git.ListOptions{
		Auth: g.AuthMethod,
	})
	if err != nil {
		return nil, err
	}

	branchRef := fmt.Sprintf("refs/heads/%s", branch)
	for _, v := range rfs {
		if v.Name().String() == branchRef {
			c, err := repo.CommitObject(v.Hash())
			if err != nil {
				return nil, err
			}
			return c, nil
		}
	}
	return nil, errors.New("commit for ref could not be found")
}
