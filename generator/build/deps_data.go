package build

import (
	"fmt"

	"github.com/frigus02/website/generator/data"
	"github.com/frigus02/website/generator/fs"
)

type dataItem struct {
	item interface{}
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
	item, err := data.GetItem(file)
	if err != nil {
		return fmt.Errorf("error getting data item: %v", err)
	}

	i.item = item
	return nil
}

func (i *dataItem) render(ctx *renderContext) error {
	return nil
}
