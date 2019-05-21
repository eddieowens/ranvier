package poller

import (
	"errors"
	"fmt"
	"github.com/two-rabbits/ranvier/src/configuration"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"io"
	"os"
	"path"
	"regexp"
	"time"
)

const GitPollerKey = "GitPoller"

const remoteName = "origin"

type OnUpdateFunction func(directory string)

type GitPoller interface {
	Start(remote, branch string, OnUpdate OnUpdateFunction) error
	Stop()
}

type gitPollerImpl struct {
	Config      configuration.Config `inject:"Config"`
	quitChannel chan bool
	repo        *git.Repository
	branchName  string
	filters     []regexp.Regexp
}

func (g *gitPollerImpl) Stop() {
	close(g.quitChannel)
}

func (g *gitPollerImpl) Start(remote, branch string, onUpdate OnUpdateFunction, filters ...regexp.Regexp) error {
	repo, err := git.PlainClone(g.Config.CloneDirectory, false, &git.CloneOptions{
		URL:           remote,
		RemoteName:    remoteName,
		ReferenceName: plumbing.NewRemoteReferenceName(remote, branch),
	})

	if err != nil {
		return err
	}

	g.repo = repo
	g.branchName = branch
	g.filters = filters

	onUpdate(g.Config.CloneDirectory)

	ticker := time.NewTicker(time.Duration(g.Config.GitPollInterval) * time.Second)
	g.quitChannel = make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				changes := g.fetchUpdates()
				if len(changes) > 0 {
					for _, c := range changes {
						fp := path.Join(g.Config.CloneDirectory, c)
						onUpdate(fp)
					}
				}
			case <-g.quitChannel:
				ticker.Stop()
				return
			}
		}
	}()

	return nil
}

func (g *gitPollerImpl) isDirEmpty(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true
	}

	return false
}

func (g *gitPollerImpl) fetchUpdates() []string {
	_ = g.repo.Fetch(&git.FetchOptions{})

	h, _ := g.repo.Head()

	remCommit, _ := g.findLatestRemoteCommit()
	currentCommit, _ := g.repo.CommitObject(h.Hash())

	originTree, _ := remCommit.Tree()
	branchTree, _ := currentCommit.Tree()

	diffs, _ := branchTree.Diff(originTree)

	changes := make([]string, 0)
	for _, d := range diffs {
		skip := false
		fp := d.To.Name
		for _, f := range g.filters {
			if f.Match([]byte(fp)) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		changes = append(changes, fp)
	}

	wt, _ := g.repo.Worktree()

	_ = wt.Pull(&git.PullOptions{
		SingleBranch: true,
	})

	return changes
}

func (g *gitPollerImpl) findLatestRemoteCommit() (*object.Commit, error) {
	rem, _ := g.repo.Remote(remoteName)
	rfs, _ := rem.List(&git.ListOptions{})
	branchRef := fmt.Sprintf("refs/heads/%s", g.branchName)
	for _, v := range rfs {
		if v.Name().String() == branchRef {
			c, err := g.repo.CommitObject(v.Hash())
			if err != nil {
				return nil, err
			}
			return c, nil
		}
	}
	return nil, errors.New("commit for ref could not be found")
}
