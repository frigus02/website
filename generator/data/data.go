package data

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/frigus02/website/generator/fs"
	"github.com/russross/blackfriday"
)

const dataDir = "data"
const dataItemFile = "index.md"

type baseItem interface {
	setID(string)
	setContent(template.HTML)
}

func walkDataDir(itemDir string, walkFunc func(path, itemDir string) error) error {
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

func getIDFromItemDir(itemDir string) string {
	// Data items can have a number at the start to enable ordering on the file
	// system. The number is usually separated from the name with a dash.
	return strings.TrimLeft(itemDir, "0123456789-")
}

func readDataItem(path, itemDir string, item baseItem) error {
	item.setID(getIDFromItemDir(itemDir))

	filename := filepath.Join(path, itemDir, dataItemFile)
	content, err := fs.ReadFileWithMetadata(filename, item)
	if err != nil {
		return err
	}

	htmlContent := blackfriday.Run([]byte(content))

	item.setContent(template.HTML(htmlContent))
	return nil
}

// GetItem reads a data item from the file system based on any file in the
// items directory, deciding the item type from the folder name.
func GetItem(file string) (interface{}, error) {
	typePath, typeDir, itemDir, _, _ := ExtractMetadataFromFilePath(file)

	var item baseItem
	switch typeDir {
	case postsDir:
		item = &Post{}
	case projectsDir:
		item = &Project{}
	default:
		return nil, fmt.Errorf("unknown type: %s", typeDir)
	}

	err := readDataItem(typePath, itemDir, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// ExtractMetadataFromFilePath extracts typePath, typeDir, itemDir and id from
// a file in the data directory.
func ExtractMetadataFromFilePath(file string) (typePath, typeDir, itemDir, fileName, id string) {
	itemPath := filepath.Dir(file)
	itemDir = filepath.Base(itemPath)
	typePath = filepath.Dir(itemPath)
	typeDir = filepath.Base(typePath)
	fileName = filepath.Base(file)
	id = getIDFromItemDir(itemDir)
	return
}
