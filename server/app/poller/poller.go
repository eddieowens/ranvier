package poller

import (
	"fmt"
	"github.com/eddieowens/ranvier/server/app/configuration"
	"github.com/eddieowens/ranvier/server/app/service"
	"gopkg.in/src-d/go-git.v4"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"
)

const GitPollerKey = "GitPoller"

type OnUpdateFunction func(filepath string)

type OnStartFunc OnUpdateFunction

type GitPoller interface {
	Start(onUpdate OnUpdateFunction, onStart OnStartFunc, filters ...regexp.Regexp) error
	Stop()
}

type gitPollerImpl struct {
	Config      configuration.Config `inject:"Config"`
	GitService  service.GitService   `inject:"GitService"`
	quitChannel chan bool
	repo        *git.Repository
	branchName  string
	filters     []regexp.Regexp
}

func (g *gitPollerImpl) Stop() {
	close(g.quitChannel)
}

func (g *gitPollerImpl) Start(onUpdate OnUpdateFunction, onStart OnStartFunc, filters ...regexp.Regexp) error {
	repo, err := g.GitService.Clone(g.Config.Git.Remote, g.Config.Git.Branch, g.Config.Git.Directory)
	if err != nil {
		return err
	}

	g.repo = repo
	g.branchName = g.Config.Git.Branch
	g.filters = filters

	err = g.initializeConfig(onStart)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Duration(g.Config.Git.PollingInterval) * time.Second)
	g.quitChannel = make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				changes, err := g.GitService.DiffRemote(g.repo, g.branchName)
				if err != nil {
					fmt.Println(err)
					continue
				}
				changes = g.filter(changes)
				if len(changes) > 0 {
					for _, c := range changes {
						fp := path.Join(g.Config.Git.Directory, c)
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

func (g *gitPollerImpl) filterFile(file string) bool {
	for _, f := range g.filters {
		if !f.Match([]byte(file)) {
			return false
		}
	}
	return true
}

func (g *gitPollerImpl) filter(files []string) []string {
	changes := make([]string, 0)
	for _, f := range files {
		if g.filterFile(f) {
			changes = append(changes, f)
		}
	}
	return changes
}

func (g *gitPollerImpl) initializeConfig(onStart OnStartFunc) error {
	return filepath.Walk(g.Config.Git.Directory, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir
		}
		if g.filterFile(path) {
			onStart(path)
		}
		return nil
	})
}
