package service

import (
	"config-manager/src/collections"
	"config-manager/src/configuration"
	"config-manager/src/model"
	"config-manager/src/state"
	"errors"
	"github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const FileServiceKey = "FileService"

const FileNameSeparator = "__"

const MetaFileMarker = "meta"

type FileService interface {
	Create(levelConfig model.LevelConfig) error
	Exists(filename string) bool
	FromMetaFile(filename string) model.LevelConfigMeta
	FromFile(filename string) model.LevelConfig
}

type fileServiceImpl struct {
	Config       configuration.Config `inject:"Config"`
	Json         jsoniter.API         `inject:"Json"`
	LevelService state.LevelService   `inject:"LevelService"`
	IdService    state.IdService      `inject:"IdService"`
}

func (f *fileServiceImpl) FromFile(filename string) model.LevelConfig {
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
	filepath := path.Join(f.Config.ConfigDirectory, filename)
	data, err := ioutil.ReadFile(filepath)
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
	filepath := path.Join(f.Config.ConfigDirectory, filename)

	bytes, _ := f.Json.Marshal(levelConfig)

	if err := ioutil.WriteFile(filepath, bytes, 0644); err != nil {
		return errors.New("failed to write file")
	}
	return nil
}

func (f *fileServiceImpl) addVersion(levelConfig model.LevelConfig) error {
	level := f.LevelService.ToString(levelConfig.Level)
	metaFile := strings.Join([]string{level, levelConfig.Id.String(), MetaFileMarker}, FileNameSeparator)
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
