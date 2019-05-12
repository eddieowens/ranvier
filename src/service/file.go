package service

import (
	"errors"
	"github.com/json-iterator/go"
	"github.com/two-rabbits/ranvier/src/collections"
	"github.com/two-rabbits/ranvier/src/configuration"
	"github.com/two-rabbits/ranvier/src/model"
	"github.com/two-rabbits/ranvier/src/state"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const FileServiceKey = "FileService"

const FileNameSeparator = "__"

const MetaFileMarker = "meta"

type FileService interface {
	Create(levelConfig model.LevelConfig) error
	Exists(filename string) bool
	Delete(levelConfig model.LevelConfig) error
	FromMetaFile(filename string) model.LevelConfigMeta
	FromConfigFile(filename string) model.LevelConfig
	FromFileName(filename string) model.LevelConfig
	IsMetaFile(filename string) bool
}

type fileServiceImpl struct {
	Config       configuration.Config `inject:"Config"`
	Json         jsoniter.API         `inject:"Json"`
	LevelService state.LevelService   `inject:"LevelService"`
	IdService    state.IdService      `inject:"IdService"`
}

func (f *fileServiceImpl) IsMetaFile(filename string) bool {
	n := strings.Split(filename, FileNameSeparator)
	return n[len(n)-1] == MetaFileMarker
}

func (f *fileServiceImpl) FromFileName(filename string) model.LevelConfig {
	var name = strings.TrimSuffix(filename, filepath.Ext(filename))
	meta := strings.Split(name, FileNameSeparator)
	return model.LevelConfig{
		Level: f.LevelService.FromString(meta[0]),
		Id:    model.Id(meta[1]),
	}
}

func (f *fileServiceImpl) Delete(levelConfig model.LevelConfig) error {
	name := f.toFilename(f.LevelService.ToString(levelConfig.Level), levelConfig.Id.String())
	metaFile := f.toMetaFileName(&levelConfig)

	err := os.Remove(path.Join(f.Config.ConfigDirectory, metaFile))
	if err != nil {
		return err
	}
	return os.Remove(path.Join(f.Config.ConfigDirectory, name))
}

func (f *fileServiceImpl) FromConfigFile(filename string) model.LevelConfig {
	fp := path.Join(f.Config.ConfigDirectory, filename)
	data, _ := ioutil.ReadFile(fp)

	var levelConfig model.LevelConfig

	err := f.Json.Unmarshal(data, &levelConfig)
	if err != nil {
		panic(err)
	}

	return levelConfig
}

func (f *fileServiceImpl) FromMetaFile(filename string) model.LevelConfigMeta {
	fp := path.Join(f.Config.ConfigDirectory, filename)
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		panic(err)
	}

	meta := model.LevelConfigMeta{}
	err = f.Json.Unmarshal(data, &meta)
	if err != nil {
		panic(err)
	}

	return meta
}

func (f *fileServiceImpl) Create(levelConfig model.LevelConfig) error {
	level := f.LevelService.ToString(levelConfig.Level)
	name := f.toFilename(level, levelConfig.Id.String())
	err := f.addVersion(levelConfig)
	if err != nil {
		return err
	}
	return f.create(name, &levelConfig)
}

func (f *fileServiceImpl) Exists(filename string) bool {
	_, err := os.Stat(path.Join(f.Config.ConfigDirectory, filename))
	return !os.IsNotExist(err)
}

func (f *fileServiceImpl) toFilename(names ...string) string {
	return strings.Join(names, FileNameSeparator) + ".json"
}

func (f *fileServiceImpl) create(filename string, levelConfig *model.LevelConfig) error {
	fp := path.Join(f.Config.ConfigDirectory, filename)

	bytes, _ := f.Json.Marshal(levelConfig)

	if err := ioutil.WriteFile(fp, bytes, 0644); err != nil {
		return errors.New("failed to write file")
	}
	return nil
}

func (f *fileServiceImpl) toMetaFileName(config *model.LevelConfig) string {
	level := f.LevelService.ToString(config.Level)
	return strings.Join([]string{level, config.Id.String(), MetaFileMarker}, FileNameSeparator)
}

func (f *fileServiceImpl) addVersion(levelConfig model.LevelConfig) error {
	metaFile := f.toMetaFileName(&levelConfig)
	metaFilepath := path.Join(f.Config.ConfigDirectory, metaFile)

	var jsonMap collections.JsonMap
	if !f.Exists(metaFile) {
		meta := model.LevelConfigMeta{
			Versions: []model.LevelConfig{},
		}
		d, _ := f.Json.Marshal(meta)
		jsonMap = collections.NewJsonMap(d)
	} else {
		data, _ := ioutil.ReadFile(metaFilepath)
		jsonMap = collections.NewJsonMap(data)
	}

	jsonMap.Set("versions.-1", levelConfig)

	return ioutil.WriteFile(metaFilepath, jsonMap.GetRaw(), 0644)
}
