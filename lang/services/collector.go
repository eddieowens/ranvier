package services

import (
	"os"
	"path/filepath"
)

const FileCollectorKey = "FileCollector"

type FileCollector interface {
	Collect(path string) []string
}

type fileCollectorImpl struct {
	FileFilter FileFilter `inject:"FileFilter"`
}

func (f *fileCollectorImpl) Collect(path string) []string {
	gathered := make([]string, 0)
	_ = filepath.Walk(path, func(fp string, info os.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		if f.FileFilter.Filter(fp) {
			gathered = append(gathered, fp)
		}

		return nil
	})

	return gathered
}
