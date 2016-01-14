package rfs

import (
	"path/filepath"
	"testing"

	"github.com/kardianos/osext"
	"github.com/stretchr/testify/assert"
)

func Test_RealReadableFileSystem_ReadDir(t *testing.T) {
	rootpath, _ := RootPath()
	fs := Reality()
	files, err := fs.ReadDir(rootpath)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, len(files))
}

func Test_RealReadableFileSystem_ReadDir_Failure(t *testing.T) {
	rootpath, _ := RootPath()
	fakepath := filepath.Join(rootpath, "PLEASEDONOTCREATEMEEVENASAJOKE")
	fs := Reality()
	_, err := fs.ReadDir(fakepath)
	assert.NotEqual(t, nil, err)
}

func Test_RealReadableFileSystem_ReadFile(t *testing.T) {
	self, _ := osext.Executable()
	fs := Reality()
	bytes, err := fs.ReadFile(self)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, len(bytes))
}

func Test_RealReadableFileSystem_Stat(t *testing.T) {
	self, _ := osext.Executable()
	fs := Reality()
	_, err := fs.Stat(self)
	assert.Equal(t, nil, err)
}

func Test_RealReadableFileSystem_Stat_Failure(t *testing.T) {
	rootpath, _ := RootPath()
	fakepath := filepath.Join(rootpath, "PLEASEDONOTCREATEMEEVENASAJOKE")
	fs := Reality()
	_, err := fs.Stat(fakepath)
	assert.NotEqual(t, nil, err)
}
