package build

import (
	"github.com/frigus02/website/generator/fs"
	"github.com/tdewolff/minify/v2"
)

type item interface {
	addItem(item)
	isDependentOn(item) bool
	update(*fs.File) error
	render(*renderContext) error
}

type settings struct {
	minifier *minify.M
}

type renderContext struct {
	out      string
	settings *settings
}

func addItemIfNotExists(items []item, newItem item) []item {
	found := false
	for _, existingItem := range items {
		if existingItem == newItem {
			found = true
		}
	}

	if !found {
		return append(items, newItem)
	}

	return items
}

func containsItem(items []item, testItem item) bool {
	for _, existingItem := range items {
		if existingItem == testItem {
			return true
		}
	}

	return false
}
