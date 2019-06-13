package services

import (
	"github.com/eddieowens/ranvier/commons"
	"github.com/eddieowens/ranvier/lang/domain"
	"path/filepath"
)

const FileFilterKey = "FileFilter"

type FileFilter interface {
	Filter(fp string) bool
}

type fileFilterImpl struct {
}

func (f *fileFilterImpl) Filter(fp string) bool {
	ext := filepath.Ext(fp)
	if ext != "" {
		// strip the dot
		if string(ext[0]) == "." {
			ext = ext[1:]
		}
		return commons.StringIncludes(ext, domain.SupportedFileTypes)
	}
	return false
}
