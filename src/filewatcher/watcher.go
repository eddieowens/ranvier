package filewatcher

import (
	"fmt"
	"github.com/eddieowens/axon"
	"github.com/radovskyb/watcher"
	"github.com/two-rabbits/ranvier/src/configuration"
	"github.com/two-rabbits/ranvier/src/service"
	"github.com/two-rabbits/ranvier/src/state"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"
)

const fileWatcherKey = "FileWatcher"

type FileWatcher interface {
	Start() error
}

type fileWatcherImpl struct {
	LevelService       state.LevelService         `inject:"LevelService"`
	State              state.LevelConfigState     `inject:"LevelConfigState"`
	LevelConfigService service.LevelConfigService `inject:"LevelConfigService"`
	FileService        service.FileService        `inject:"FileService"`
	Config             configuration.Config       `inject:"Config"`
	Watcher            *watcher.Watcher
}

func (f *fileWatcherImpl) Start() error {
	f.Watcher.FilterOps(watcher.Create, watcher.Write)

	f.Watcher.AddFilterHook(watcher.RegexFilterHook(regexp.MustCompile(".+.json"), false))

	if err := f.Watcher.Add(f.Config.ConfigDirectory); err != nil {
		return err
	}

	files, err := ioutil.ReadDir(f.Config.ConfigDirectory)
	if err != nil {
		return err
	}
	for _, v := range files {
		if !v.IsDir() && f.isMetaFile(v.Name()) {
			f.hydrateStateFromMetaFile(v.Name())
		}
	}

	go func() {
		for {
			select {
			case event := <-f.Watcher.Event:
				if !event.IsDir() {
					f.updateStateFromFile(event.Name())
				}

			case err := <-f.Watcher.Error:
				fmt.Println(err)
			case <-f.Watcher.Closed:
				fmt.Println("closed")
				return
			}
		}
	}()

	go func() {
		if err := f.Watcher.Start(time.Second * 1); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

func (f *fileWatcherImpl) updateStateFromFile(filename string) {
	levelConfig := f.FileService.FromFile(filename)
	f.State.Set(levelConfig)
}

func (f *fileWatcherImpl) isMetaFile(filename string) bool {
	return strings.Contains(filename, service.FileNameSeparator+service.MetaFileMarker)
}

func (f *fileWatcherImpl) hydrateStateFromMetaFile(filename string) {
	meta := f.FileService.FromMetaFile(filename)
	for _, v := range meta.Versions {
		f.State.Set(v)
	}
}

func fileWatcherFactory(_ axon.Args) axon.Instance {
	return axon.StructPtr(&fileWatcherImpl{
		Watcher: watcher.New(),
	})
}
