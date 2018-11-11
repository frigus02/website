package build

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/frigus02/website/generator/data"
	"github.com/frigus02/website/generator/fs"
)

type pageMetadata struct {
	Title string `yaml:"title"`
}

type pageItem struct {
	id             string
	metadata       *pageMetadata
	tmpl           *template.Template
	dependsOnItems []item
}

func newPageItem(file *fs.File) (*pageItem, error) {
	item := &pageItem{}
	err := item.update(file)
	return item, err
}

func (i *pageItem) addItem(item item) {
	switch item := item.(type) {
	case *layoutItem:
		i.dependsOnItems = addItemIfNotExists(i.dependsOnItems, item)
	case *pageItem:
		if strings.HasPrefix(i.id, item.id) {
			i.dependsOnItems = addItemIfNotExists(i.dependsOnItems, item)
		}
	case *staticFileItem:
		i.dependsOnItems = addItemIfNotExists(i.dependsOnItems, item)
	case *dataItem:
		i.dependsOnItems = addItemIfNotExists(i.dependsOnItems, item)
	}
}

func (i *pageItem) isDependentOn(item item) bool {
	return containsItems(i.dependsOnItems, item)
}

func (i *pageItem) update(file *fs.File) error {
	metadata, tmpl, err := loadPageFile(file)
	if err != nil {
		return fmt.Errorf("error loading page file: %v", err)
	}

	id := strings.TrimSuffix(file.Name[6:], filepath.Ext(file.Name))
	if id == "index" {
		id = ""
	} else {
		id = strings.TrimSuffix(id, "/index")
	}

	i.id = id
	i.metadata = metadata
	i.tmpl = tmpl
	return nil
}

func (i *pageItem) render(ctx *renderContext) error {
	var layout *layoutItem
	var parentPage *pageItem
	var staticFiles []*staticFileItem
	var posts []*data.Post
	var projects []*data.Project
	for _, item := range i.dependsOnItems {
		switch item := item.(type) {
		case *layoutItem:
			layout = item
		case *pageItem:
			parentPage = item
		case *staticFileItem:
			staticFiles = append(staticFiles, item)
		case *dataItem:
			switch dataItem := item.item.(type) {
			case *data.Post:
				posts = append(posts, dataItem)
			case *data.Project:
				projects = append(projects, dataItem)
			}
		}
	}

	if layout == nil {
		return fmt.Errorf("no layout found")
	}

	var outfile string
	if i.id == "" {
		outfile = filepath.Join(ctx.out, "index.html")
	} else {
		outfile = filepath.Join(ctx.out, i.id+".html")
	}

	staticFileNames := getStaticFileMap(staticFiles)
	pageCtx := &pageContext{
		Posts:       posts,
		Projects:    projects,
		StaticFiles: staticFileNames,
	}

	err := renderPageToFile(i.id, outfile, i.metadata, i.tmpl, ctx.settings, pageCtx, layout, parentPage, staticFileNames)
	if err != nil {
		return fmt.Errorf("error rendering page: %v", err)
	}

	return nil
}

func loadPageFile(file *fs.File) (*pageMetadata, *template.Template, error) {
	metadata := pageMetadata{}
	content, err := file.SplitMetadataAndContent(&metadata)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading file with metadata %s: %v", file, err)
	}

	tmpl, err := parseTemplate(content)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing page template %s: %v", file, err)
	}

	return &metadata, tmpl, nil
}

func getStaticFileMap(staticFiles []*staticFileItem) map[string]string {
	staticFileNames := make(map[string]string)
	for _, staticFile := range staticFiles {
		staticFileNames[strings.TrimPrefix(staticFile.origName, "static/")] = staticFile.name
	}

	return staticFileNames
}

func renderPageToFile(
	id string,
	outfile string,
	metadata *pageMetadata,
	tmpl *template.Template,
	settings *settings,
	pageCtx interface{},
	layout *layoutItem,
	parentPage *pageItem,
	staticFileNames map[string]string,
) error {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, pageCtx)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", id, err)
	}

	err = os.MkdirAll(filepath.Dir(outfile), 0644)
	if err != nil {
		return fmt.Errorf("error creating destination folder %s for page %s: %v", outfile, id, err)
	}

	var destination io.WriteCloser
	destination, err = os.Create(outfile)
	if err != nil {
		return fmt.Errorf("error creating destination %s for page %s: %v", outfile, id, err)
	}

	if settings.minifier != nil {
		destination = settings.minifier.Writer("text/html", destination)
	}

	layoutCtx := layoutContext{
		ID:          id,
		Title:       metadata.Title,
		Content:     template.HTML(buf.String()),
		StaticFiles: staticFileNames,
	}
	if parentPage != nil {
		layoutCtx.ParentID = parentPage.id
		layoutCtx.ParentTitle = parentPage.metadata.Title
	}

	err = layout.tmpl.Execute(destination, &layoutCtx)
	if err != nil {
		destination.Close()
		return fmt.Errorf("error executing layout template for %s: %v", id, err)
	}

	err = destination.Close()
	if err != nil {
		return fmt.Errorf("error closing (minify) writer for %s: %v", id, err)
	}

	return nil
}
