package fs

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// File is an abstraction of a file on the file system with direct access to
// the content in bytes.
type File struct {
	Name    string
	Content []byte
}

// ReadFile reads the specified file from disk and returns a new instance of
// File.
func ReadFile(path, name string) (*File, error) {
	bytes, err := ioutil.ReadFile(filepath.Join(path, name))
	if err != nil {
		return nil, err
	}

	return &File{
		Name:    strings.Replace(name, string(filepath.Separator), "/", -1),
		Content: bytes,
	}, nil
}

// SplitMetadataAndContent reads a file, which contains YAML metadata at the
// start.
func (f *File) SplitMetadataAndContent(outMetadata interface{}) (content string, err error) {
	data := string(f.Content)

	// Parse metadata.
	metadataStart := strings.Index(data, "---")
	metadataEnd := strings.LastIndex(data, "---")
	if metadataStart == -1 || metadataStart == metadataEnd {
		return "", fmt.Errorf("no metadata found (first and last index of --- were: %v, %v)", metadataStart, metadataEnd)
	}

	metadata := data[metadataStart+3 : metadataEnd]
	err = yaml.Unmarshal([]byte(metadata), outMetadata)
	if err != nil {
		return "", err
	}

	return data[metadataEnd+3:], nil
}
