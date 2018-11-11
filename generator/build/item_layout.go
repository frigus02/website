package build

import (
	"fmt"
	"html/template"

	"github.com/frigus02/website/generator/fs"
)

type layoutItem struct {
	tmpl *template.Template
}

func newLayoutItem(file *fs.File) (*layoutItem, error) {
	item := &layoutItem{}
	err := item.update(file)
	return item, err
}

func (i *layoutItem) addItem(item item) {
}

func (i *layoutItem) isDependentOn(item item) bool {
	return false
}

func (i *layoutItem) update(file *fs.File) error {
	tmpl, err := parseTemplate(string(file.Content))
	if err != nil {
		return fmt.Errorf("error parsing _layout.html template: %v", err)
	}

	i.tmpl = tmpl
	return nil
}

func (i *layoutItem) render(ctx *renderContext) error {
	return nil
}
