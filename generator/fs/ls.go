package fs

import (
	"os"
	"path/filepath"
	"strings"
)

// ListRecursive lists all files and folders recursively in the specified path.
// Returned paths are relative to the specified path.
func ListRecursive(path string) (files, folders []string, err error) {
	path = filepath.Clean(path) + string(filepath.Separator)
	err = filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		pathRelativeToInput := strings.TrimPrefix(newPath, path)

		if info.IsDir() {
			folders = append(folders, pathRelativeToInput)
		} else {
			files = append(files, pathRelativeToInput)
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return
}
