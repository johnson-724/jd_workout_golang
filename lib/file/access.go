package file

import (
	"os"
	"path/filepath"
)

func AccessFromCurrentDir(path string) string {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	absPath, err := filepath.Abs(filepath.Join(wd, path))

	if err != nil {
		panic(err)
	}

	return absPath
}
