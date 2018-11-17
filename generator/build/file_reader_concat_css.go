package build

import (
	"fmt"
	"path/filepath"

	"github.com/frigus02/website/generator/fs"
)

type cssConcatFileReader struct {
	cssFiles map[string]*fs.File
}

func newCSSConcatFileReader() *cssConcatFileReader {
	return &cssConcatFileReader{
		cssFiles: make(map[string]*fs.File),
	}
}

func (fr *cssConcatFileReader) readFile(path, name string) (*fs.File, error) {
	file, err := fs.ReadFile(path, name)
	if err != nil {
		return nil, err
	}

	if filepath.Ext(file.Name) == ".css" {
		fr.cssFiles[file.Name] = file

		file = &fs.File{
			Name: "static/styles.css",
		}
		for _, cssFile := range fr.cssFiles {
			header := []byte(fmt.Sprintf("\n/******************************\n * %s\n */\n\n", cssFile.Name))
			file.Content = append(file.Content, header...)
			file.Content = append(file.Content, cssFile.Content...)
		}
	}

	return file, nil
}
