package mocks

import (
	"github.com/stretchr/testify/mock"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func Repository(args mock.Arguments, i int) *git.Repository {
	var r *git.Repository
	v := args.Get(i)
	if v != nil {
		r = v.(*git.Repository)
	}
	return r
}

func Commit(args mock.Arguments, i int) *object.Commit {
	var r *object.Commit
	v := args.Get(i)
	if v != nil {
		r = v.(*object.Commit)
	}
	return r
}

func StringSlice(args mock.Arguments, i int) []string {
	var r []string
	v := args.Get(i)
	if v != nil {
		r = v.([]string)
	}
	return r
}
