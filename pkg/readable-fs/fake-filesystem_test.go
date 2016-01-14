package rfs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func Test_FakeReadableFileSystem_ReadDir_Failure(t *testing.T) {
	fs := NewFake()
	_, err := fs.ReadDir("/not/present")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "open /not/present: No such file or directory", err.Error())
}

func Test_FakeReadableFileSystem_ReadFile_Failure(t *testing.T) {
	fs := NewFake()
	_, err := fs.ReadFile("/not/present")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "open /not/present: No such file or directory", err.Error())
}

func Test_FakeReadableFileSystem_Stat_Failure(t *testing.T) {
	fs := NewFake()
	_, err := fs.Stat("/not/present")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "stat /not/present: No such file or directory", err.Error())
}

type FakeReadableFileSystem_Suite struct {
	suite.Suite
	RootPath     string
	DirPath      string
	FilePath     string
	FileContents []byte
	FileSystem   ReadableFileSystem
	FileMode     os.FileMode
}

func (suite *FakeReadableFileSystem_Suite) SetupTest() {
	rootpath, err := RootPath()
	if err != nil {
		panic("can't find root")
	}

	suite.RootPath = rootpath
	suite.DirPath = filepath.Join(suite.RootPath, "bin")
	suite.FilePath = filepath.Join(suite.DirPath, "isgood")
	suite.FileContents = FromString("#!/bin/sh\ntrue\n")
	suite.FileMode = 0755

	fs := NewFake()
	fs.AddFile(suite.FilePath, suite.FileContents, suite.FileMode)

	suite.FileSystem = fs
}

func Test_FakeReadableFileSystem(t *testing.T) {
	suite.Run(t, new(FakeReadableFileSystem_Suite))
}

func (suite *FakeReadableFileSystem_Suite) Test_FakeReadableFileSystem_ReadDir_Top() {
	t := suite.T()

	files, err := suite.FileSystem.ReadDir(suite.RootPath)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(files))

	dir := files[0]
	assert.Equal(t, filepath.Base(suite.DirPath), dir.Name())
	assert.Equal(t, os.ModeDir, dir.Mode())
}

func (suite *FakeReadableFileSystem_Suite) Test_FakeReadableFileSystem_ReadDir_Middle() {
	t := suite.T()

	files, err := suite.FileSystem.ReadDir(suite.DirPath)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(files))

	file := files[0]
	assert.Equal(t, filepath.Base(suite.FilePath), file.Name())
	assert.Equal(t, int64(len(suite.FileContents)), file.Size())
	assert.Equal(t, suite.FileMode, file.Mode())
}

func (suite *FakeReadableFileSystem_Suite) Test_FakeReadableFileSystem_ReadDir_TrailingSlash() {
	t := suite.T()

	files, err := suite.FileSystem.ReadDir(suite.DirPath + "/")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(files))

	file := files[0]
	assert.Equal(t, filepath.Base(suite.FilePath), file.Name())
	assert.Equal(t, int64(len(suite.FileContents)), file.Size())
	assert.Equal(t, suite.FileMode, file.Mode())
}

func (suite *FakeReadableFileSystem_Suite) Test_FakeReadableFileSystem_ReadFile() {
	t := suite.T()

	bytes, err := suite.FileSystem.ReadFile(suite.FilePath)
	assert.Equal(t, nil, err)
	assert.Equal(t, suite.FileContents, bytes)
}

func (suite *FakeReadableFileSystem_Suite) Test_FakeReadableFileSystem_Stat() {
	t := suite.T()

	file, err := suite.FileSystem.Stat(suite.FilePath)
	assert.Equal(t, nil, err)
	assert.Equal(t, filepath.Base(suite.FilePath), file.Name())
	assert.Equal(t, int64(len(suite.FileContents)), file.Size())
	assert.Equal(t, suite.FileMode, file.Mode())
}
