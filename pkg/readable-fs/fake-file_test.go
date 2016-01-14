package rfs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewFakeReadableFile_Name(t *testing.T) {
	f := FakeReadableFile{
		basename: "foo",
	}
	assert.Equal(t, "foo", f.Name())
}

func Test_NewFakeReadableFile_Size_0(t *testing.T) {
	f := FakeReadableFile{}
	assert.Equal(t, int64(0), f.Size())
}

func Test_NewFakeReadableFile_Size_5(t *testing.T) {
	f := FakeReadableFile{
		bytes: FromString("bytes"),
	}
	assert.Equal(t, int64(5), f.Size())
}

func Test_NewFakeReadableFile_Mode(t *testing.T) {
	f := FakeReadableFile{
		mode: 0644,
	}
	assert.Equal(t, os.FileMode(0644), f.Mode())
}
