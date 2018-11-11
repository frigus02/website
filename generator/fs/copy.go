package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies the specified file to the specified location, creating the
// necessary destination folders and overwriting the existing destination file.
func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file %s: %v", src, err)
	}

	defer source.Close()

	dstDir := filepath.Dir(dst)
	err = os.MkdirAll(dstDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating destination folder %s: %v", dstDir, err)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating destination file %s: %v", dst, err)
	}

	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("error copying file %s to %s: %v", src, dst, err)
	}

	return nil
}
