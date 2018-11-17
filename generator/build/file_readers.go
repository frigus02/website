package build

import (
	"github.com/frigus02/website/generator/fs"
)

type fileReader interface {
	readFile(path, name string) (*fs.File, error)
}
