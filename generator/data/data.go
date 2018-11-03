package data

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const dataDir = "../site/data"
const dataItemFile = "index.md"

type baseItem interface {
	setID(string)
	setContent(string)
}

func walkDataDir(itemDir string, walkFunc func(parentDir, itemDir string) error) error {
	dir := filepath.Join(dataDir, itemDir)
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking dir on %s: %v", path, err)
		}

		if dir == path || !info.IsDir() {
			return nil
		}

		err = walkFunc(dir, info.Name())
		if err != nil {
			return fmt.Errorf("error reading data item %s: %v", path, err)
		}

		return nil
	})
}

func readDataItem(parentDir, itemDir string, item baseItem) error {
	// Data items can have a number at the start to enable ordering on the file
	// system. The number is usually separated from the name with a dash.
	item.setID(strings.TrimLeft(itemDir, "0123456789-"))

	// Read file.
	filename := filepath.Join(parentDir, itemDir, dataItemFile)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	data := string(content)

	// Parse metadata.
	metadataStart := strings.Index(data, "---")
	metadataEnd := strings.LastIndex(data, "---")
	if metadataStart == -1 || metadataStart == metadataEnd {
		return fmt.Errorf("no metadata found (first and last index of --- were: %v, %v)", metadataStart, metadataEnd)
	}

	metadata := data[metadataStart+3 : metadataEnd]
	err = yaml.Unmarshal([]byte(metadata), item)
	if err != nil {
		return err
	}

	// Use rest as content.
	item.setContent(data[metadataEnd+3:])
	return nil
}
