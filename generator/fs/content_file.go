package fs

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// ReadFileWithMetadata reads a file, which contains YAML metadata at the
// start.
func ReadFileWithMetadata(file string, outMetadata interface{}) (content string, err error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	data := string(bytes)

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
