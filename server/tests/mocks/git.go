package mocks

import (
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/stretchr/testify/mock"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type GitServiceMock struct {
	mock.Mock
}

func (g *GitServiceMock) Clone(remote, branch, directory string) (*git.Repository, error) {
	args := g.Called(remote, branch, directory)
	return Repository(args, 0), args.Error(1)
}

func (g *GitServiceMock) DiffRemote(repo *git.Repository, branch string) ([]model.GitChange, error) {
	args := g.Called(repo)
	return GitChangeSlice(args, 0), args.Error(1)
}

func (g *GitServiceMock) FetchLatestRemoteCommit(repo *git.Repository, branch string) (*object.Commit, error) {
	args := g.Called(repo)
	return Commit(args, 0), args.Error(1)
}
