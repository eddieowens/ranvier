package poller

import (
	"fmt"
	"github.com/src-d/go-git/plumbing/object"
	"github.com/two-rabbits/ranvier/src/configuration"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"io"
	"os"
	"path"
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
}

func (g *gitPollerImpl) Stop() {
	close(g.quitChannel)
}

func (g *gitPollerImpl) Start(remote, branch string, onUpdate OnUpdateFunction) error {
	repo, err := git.PlainClone(g.Config.CloneDirectory, false, &git.CloneOptions{
		URL:           remote,
		RemoteName:    remoteName,
		ReferenceName: plumbing.NewRemoteReferenceName(remote, branch),
	})

	if err != nil {
		return err
	}

	g.repo = repo

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
	//t, _ := g.repo.Worktree()
	h, _ := g.repo.Head()
	_ = g.repo.Fetch(&git.FetchOptions{})
	rem, _ := g.repo.Remote(remoteName)
	rfs, _ := rem.List(&git.ListOptions{})
	latest := rfs[0].Hash()
	originTree, _ := object.GetTree(g.repo.Storer, latest)
	branchTree, _ := object.GetTree(g.repo.Storer, h.Hash())

	c, _ := branchTree.Diff(originTree)
	if c.Len() <= 0 {
		return nil
	}

	fmt.Println(c)

	//_ := t.Pull(&git.PullOptions{
	//	SingleBranch: true,
	//})

	return nil
}
