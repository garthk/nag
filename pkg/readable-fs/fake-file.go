package rfs // import "github.com/garthk/nag/pkg/readable-fs"

import (
	"bytes"
	"os"
)

type FakeReadableFile struct {
	basename string
	bytes    []byte
	members  []*FakeReadableFile
	mode     os.FileMode
}

func (f FakeReadableFile) Name() string {
	return f.basename
}

func (f FakeReadableFile) Size() int64 {
	return int64(len(f.bytes))
}

func (f FakeReadableFile) Mode() os.FileMode {
	return f.mode
}

func FromString(str string) []byte {
	return bytes.NewBufferString(str).Bytes()
}
