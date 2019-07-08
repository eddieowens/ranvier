package services

import (
	"github.com/eddieowens/ranvier/lang/services"
	"github.com/eddieowens/ranvier/lang/tests/unit"
	"github.com/stretchr/testify/suite"
	"path"
	"testing"
)

type FileTest struct {
	unit.Unit
	fileSvc services.FileService
}

func (f *FileTest) SetupTest() {
	f.fileSvc = f.Injector.GetStructPtr(services.FileServiceKey).(services.FileService)
}

func (f *FileTest) TestSubtractPath() {
	// -- Given
	//
	root := "/path/to/file"
	fp := path.Join(root, "some", "sub", "path")

	expected := "some/sub/path"

	// -- When
	//
	actual := f.fileSvc.SubtractPath(root, fp)

	// -- Then
	//
	f.Equal(expected, actual)
}

func (f *FileTest) TestSubtractPathTrailingSlash() {
	// -- Given
	//
	root := "/path/to/file/"
	fp := path.Join(root, "some", "sub", "path")

	expected := "some/sub/path"

	// -- When
	//
	actual := f.fileSvc.SubtractPath(root, fp)

	// -- Then
	//
	f.Equal(expected, actual)
}

func (f *FileTest) TestSubtractPathRootNotPresent() {
	// -- Given
	//
	root := "/path/to/file/"
	fp := path.Join("/some", "sub", "path")

	expected := "/some/sub/path"

	// -- When
	//
	actual := f.fileSvc.SubtractPath(root, fp)

	// -- Then
	//
	f.Equal(expected, actual)
}

func (f *FileTest) TestSubtractPaths() {
	// -- Given
	//
	root := "/path/to/file"
	paths := []string{
		path.Join(root, "some", "sub", "path"),
		path.Join(root, "some", "other", "path"),
		path.Join(root, "some", "sub", "bleh"),
		path.Join(root, "wow", "sub", "bleh"),
	}

	expected := []string{
		"some/sub/path",
		"some/other/path",
		"some/sub/bleh",
		"wow/sub/bleh",
	}

	// -- When
	//
	actual := f.fileSvc.SubtractPaths(root, paths)

	// -- Then
	//
	f.Equal(expected, actual)
}

func TestFileTest(t *testing.T) {
	suite.Run(t, new(FileTest))
}
