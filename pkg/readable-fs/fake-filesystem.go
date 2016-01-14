package rfs

import (
	"errors"
	"os"
	"path/filepath"
)

type FakeReadableFileSystem struct {
	members map[string]*FakeReadableFile
}

func NewFake() FakeReadableFileSystem {
	return FakeReadableFileSystem{
		members: make(map[string]*FakeReadableFile),
	}
}

func (fs FakeReadableFileSystem) add(dirname string, fakeFile *FakeReadableFile) {
	filename := filepath.Join(dirname, fakeFile.basename)
	fs.members[filename] = fakeFile

	dir, present := fs.members[dirname]

	if !present {
		dir = &FakeReadableFile{
			basename: filepath.Base(dirname),
			mode:     os.ModeDir,
		}

		nextdirname := filepath.Dir(dirname)

		if nextdirname == dirname { // root
			fs.members[dirname] = dir
		} else {
			fs.add(nextdirname, dir)
		}
	}

	dir.members = append(dir.members, fakeFile)
}

func (fs FakeReadableFileSystem) AddFile(filename string, bytes []byte, mode os.FileMode) {
	dirname := filepath.Dir(filename)
	basename := filepath.Base(filename)

	fs.add(dirname, &FakeReadableFile{
		basename: basename,
		bytes:    bytes,
		mode:     mode,
	})
}

func (fs FakeReadableFileSystem) ReadDir(dirname string) ([]MinimalFileInfo, error) {
	dirnamelen := len(dirname)
	shortestlen := len(filepath.VolumeName(dirname)) + 1 // C:\ or /

	if dirnamelen > shortestlen && dirname[dirnamelen-1] == os.PathSeparator {
		dirname = dirname[:dirnamelen-1]
	}

	dir, present := fs.members[dirname]
	if present {
		members := make([]MinimalFileInfo, 0, len(dir.members))
		for i := range dir.members {
			members = append(members, dir.members[i])
		}
		return members, nil
	} else {
		return nil, notFoundError(dirname, "open")
	}
}

func (fs FakeReadableFileSystem) ReadFile(filename string) ([]byte, error) {
	file, present := fs.members[filename]
	if present {
		return file.bytes, nil
	} else {
		return nil, notFoundError(filename, "open")
	}
}

func (fs FakeReadableFileSystem) Stat(filename string) (MinimalFileInfo, error) {
	file, present := fs.members[filename]
	if present {
		return file, nil
	} else {
		return nil, notFoundError(filename, "stat")
	}
}

func notFoundError(filename string, op string) error {
	return &os.PathError{
		Op:   op,
		Path: filename,
		Err:  errors.New("No such file or directory"),
	}
}
