package service

import "fmt"

type FileExistsError struct {
	filename string
}

func (f *FileExistsError) Error() string {
	return fmt.Sprintf("%s already exists", f.filename)
}

func NewFileExistsError(filename string) *FileExistsError {
	return &FileExistsError{filename: filename}
}
