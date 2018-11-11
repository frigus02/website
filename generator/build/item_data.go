package build

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/frigus02/website/generator/data"
	"github.com/frigus02/website/generator/fs"
	"github.com/russross/blackfriday"
)

type dataItem struct {
	id       string
	order    int
	content  template.HTML
	metadata interface{}
}

func newDataItem(file *fs.File) (*dataItem, error) {
	item := &dataItem{}
	err := item.update(file)
	return item, err
}

func (i *dataItem) addItem(item item) {
}

func (i *dataItem) isDependentOn(item item) bool {
	return false
}

func (i *dataItem) update(file *fs.File) error {
	dataType, id, order, err := getMetadataFromFileName(file.Name)
	if err != nil {
		return fmt.Errorf("error extracting metadata from file name: %v", err)
	}

	var metadata interface{}
	switch dataType {
	case data.PostsDir:
		metadata = &data.PostMetadata{}
	case data.ProjectsDir:
		metadata = &data.ProjectMetadata{}
	default:
		return fmt.Errorf("unknown type: %s", dataType)
	}

	content, err := file.SplitMetadataAndContent(metadata)
	if err != nil {
		return err
	}

	htmlContent := blackfriday.Run([]byte(content))

	i.id = id
	i.order = order
	i.content = template.HTML(htmlContent)
	i.metadata = metadata
	return nil
}

func (i *dataItem) render(ctx *renderContext) error {
	return nil
}

func getMetadataFromFileName(name string) (dataType, id string, order int, err error) {
	itemPath := filepath.Dir(name)
	itemDir := filepath.Base(itemPath)
	typePath := filepath.Dir(itemPath)
	dataType = filepath.Base(typePath)

	// Data items have a number at the start to enable ordering on the file
	// system. The number is separated from the name with a dash.
	dashIndex := strings.Index(itemDir, "-")
	order, err = strconv.Atoi(itemDir[:dashIndex])
	id = itemDir[dashIndex+1:]
	return
}
