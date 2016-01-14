package rfs

import (
	"os"
)

type MinimalFileInfo interface {
	Name() string
	Size() int64
	Mode() os.FileMode
}

type ReadableFileSystem interface {
	ReadDir(dirname string) ([]MinimalFileInfo, error)
	ReadFile(filename string) ([]byte, error)
	Stat(filename string) (MinimalFileInfo, error)
}
