package build

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"

	"github.com/frigus02/website/generator/fs"
)

// Build builds the website and is able to continuously update the output.
type Build struct {
	In     string
	Out    string
	Minify bool

	renderCtx  *renderContext
	items      map[string]item
	dirtyItems []item
}

// Build builds the site once and returns.
func (b *Build) Build() {
	files, _, err := fs.ListRecursive(b.In)
	if err != nil {
		log.Fatal(err)
	}

	b.renderCtx = newRenderContext(b.Out)
	b.items = make(map[string]item)
	if b.Minify {
		b.renderCtx.settings.minifier = newMinifier()
	}

	b.items["#feed"] = newFeedItem()

	log.Printf("Reading files...\n")
	for _, file := range files {
		err = b.handleFile(file)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Rendering site...\n")
	err = b.render()
	if err != nil {
		log.Fatal(err)
	}
}

func (b *Build) handleFile(name string) error {
	file, err := fs.ReadFile(b.In, name)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", name, err)
	}

	item, ok := b.items[file.Name]
	if ok {
		err = item.update(file)
		if err != nil {
			return fmt.Errorf("error handling update to file %s: %v", file.Name, err)
		}
	} else {
		topDir := file.Name[0:strings.Index(file.Name, "/")]
		if topDir == "data" {
			item, err = b.createDataItem(file)
		} else if topDir == "pages" {
			item, err = b.createPageItem(file)
		} else if topDir == "static" {
			item, err = b.createStaticFileItem(file)
		} else {
			err = fmt.Errorf("unknown file")
		}

		if err != nil {
			return fmt.Errorf("error handling file %s: %v", file.Name, err)
		}

		for _, existingItem := range b.items {
			item.addItem(existingItem)
			existingItem.addItem(item)
		}

		b.items[file.Name] = item
	}

	b.dirtyItems = addItemIfNotExists(b.dirtyItems, item)
	for _, existingItem := range b.items {
		if existingItem.isDependentOn(item) {
			b.dirtyItems = addItemIfNotExists(b.dirtyItems, existingItem)
		}
	}

	return nil
}

func (b *Build) createDataItem(file *fs.File) (item, error) {
	fileName := filepath.Base(file.Name)
	if fileName != "index.md" {
		typeDir, id, _, err := getMetadataFromFileName(file.Name)
		if err != nil {
			return nil, err
		}

		newName := fmt.Sprintf("static/data/%s/%s/%s", typeDir, id, fileName)

		return b.createStaticFileItem(&fs.File{
			Name:    newName,
			Content: file.Content,
		})
	}

	return newDataItem(file)
}

func (b *Build) createPageItem(file *fs.File) (item, error) {
	if file.Name == "pages/_layout.html" {
		return newLayoutItem(file)
	} else if filepath.Base(file.Name) == "_details.html" {
		return newDataPageItem(file)
	} else {
		return newPageItem(file)
	}
}

func (b *Build) createStaticFileItem(file *fs.File) (item, error) {
	return newStaticFileItem(file, b.renderCtx.settings)
}

func (b *Build) render() error {
	for _, item := range b.dirtyItems {
		if err := item.render(b.renderCtx); err != nil {
			return err
		}
	}

	b.dirtyItems = nil
	return nil
}

func newRenderContext(out string) *renderContext {
	return &renderContext{
		out:      out,
		settings: &settings{},
	}
}

func newMinifier() *minify.M {
	minifier := minify.New()
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("text/html", html.Minify)
	return minifier
}
