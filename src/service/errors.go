package service

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

func NewKeyNotFoundError(key string) error {
	return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("the %s key could not be found", key))
}

type FileExistsError struct {
	filename string
}

func (f *FileExistsError) Error() string {
	return fmt.Sprintf("%s already exists", f.filename)
}

func NewFileExistsError(filename string) error {
	return &FileExistsError{filename: filename}
}
