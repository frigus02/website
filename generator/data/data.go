package data

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/frigus02/website/generator/fs"
	"github.com/russross/blackfriday"
)

const dataDir = "data"
const dataItemFile = "index.md"

type baseItem interface {
	setID(string)
	setOrder(int)
	setContent(template.HTML)
}

// GetItem reads a data item from the file system based on any file in the
// items directory, deciding the item type from the folder name.
func GetItem(file *fs.File) (interface{}, error) {
	typeDir, _, id, order, err := ExtractMetadataFromFilePath(file.Name)
	if err != nil {
		return nil, err
	}

	var item baseItem
	switch typeDir {
	case postsDir:
		item = &Post{}
	case projectsDir:
		item = &Project{}
	default:
		return nil, fmt.Errorf("unknown type: %s", typeDir)
	}

	content, err := file.SplitMetadataAndContent(item)
	if err != nil {
		return nil, err
	}

	htmlContent := blackfriday.Run([]byte(content))

	item.setID(id)
	item.setOrder(order)
	item.setContent(template.HTML(htmlContent))

	return item, nil
}

// ExtractMetadataFromFilePath extracts typePath, typeDir, itemDir and id from
// a file in the data directory.
func ExtractMetadataFromFilePath(file string) (typeDir, fileName, id string, order int, err error) {
	itemPath := filepath.Dir(file)
	itemDir := filepath.Base(itemPath)
	typePath := filepath.Dir(itemPath)
	typeDir = filepath.Base(typePath)
	fileName = filepath.Base(file)

	// Data items have a number at the start to enable ordering on the file
	// system. The number is separated from the name with a dash.
	dashIndex := strings.Index(itemDir, "-")
	order, err = strconv.Atoi(itemDir[:dashIndex])
	id = itemDir[dashIndex+1:]

	return
}
