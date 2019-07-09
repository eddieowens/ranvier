package services

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/eddieowens/ranvier/lang/domain"
	json "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

const FileServiceKey = "FileService"

type FileService interface {
	ToFile(directory string, config *domain.Schema) error

	// Removes the root path from the filepath (fp) if it's present. If the root path cannot be found in the fp, the fp
	// is returned.
	SubtractPath(root, fp string) string

	// Same as SubtractPath but for a slice of filepaths (fps).
	SubtractPaths(root string, fps []string) []string
}

type fileServiceImpl struct {
}

func (f *fileServiceImpl) SubtractPaths(root string, fps []string) []string {
	sl := make([]string, len(fps))
	for i, v := range fps {
		sl[i] = f.SubtractPath(root, v)
	}
	return sl
}

func (f *fileServiceImpl) SubtractPath(root, fp string) string {
	p, err := filepath.Rel(root, fp)
	if err != nil {
		return ""
	}
	return p
}

func (f *fileServiceImpl) ToFile(directory string, config *domain.Schema) error {
	if config == nil {
		return errors.New("schema cannot be nil")
	}

	var jsonConfig interface{}
	err := json.Unmarshal(config.Config, &jsonConfig)
	if err != nil {
		return err
	}

	var data []byte
	switch config.Type {
	case "toml":
		var buf bytes.Buffer
		err = toml.NewEncoder(&buf).Encode(jsonConfig)
		data = buf.Bytes()
	case "yaml", "yml":
		data, err = yaml.Marshal(jsonConfig)
	default:
		data = config.Config
	}

	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s.%s", config.Name, config.Type)
	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(directory, filename), data, 0644)
}
