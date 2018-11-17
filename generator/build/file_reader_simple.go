package build

import (
	"github.com/frigus02/website/generator/fs"
)

type simpleFileReader struct {
}

func newSimpleFileReader() *simpleFileReader {
	return &simpleFileReader{}
}

func (fr *simpleFileReader) readFile(path, name string) (*fs.File, error) {
	return fs.ReadFile(path, name)
}
