package rfs

import (
	"path/filepath"

	"github.com/kardianos/osext"
)

func RootPath() (string, error) {
	path, err := osext.Executable()
	if err != nil {
		return "", err
	}

	for {
		parent := filepath.Dir(path)
		if parent == path {
			break
		} else {
			path = parent
		}
	}

	return path, nil
}
