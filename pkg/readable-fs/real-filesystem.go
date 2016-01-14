package rfs

import (
	"io/ioutil"
	"os"
)

type RealReadableFileSystem struct{}

func Reality() RealReadableFileSystem {
	return RealReadableFileSystem{}
}

func (fs RealReadableFileSystem) ReadDir(dirname string) ([]MinimalFileInfo, error) {
	members, err := ioutil.ReadDir(dirname)

	if err != nil {
		return nil, err
	}

	minimals := make([]MinimalFileInfo, len(members))
	for i := range members {
		minimals[i] = members[i]
	}
	return minimals, nil
}

func (fs RealReadableFileSystem) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (fs RealReadableFileSystem) Stat(filename string) (MinimalFileInfo, error) {
	return os.Stat(filename)
}
