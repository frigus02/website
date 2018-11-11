package build

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	txtTemplate "text/template"

	"github.com/frigus02/website/generator/data"
	"github.com/frigus02/website/generator/fs"
	yaml "gopkg.in/yaml.v2"
)

type dataPageItem struct {
	dataType     string
	metadata     *pageMetadata
	tmpl         *template.Template
	dependencies []item
}

func newDataPageItem(file *fs.File) (*dataPageItem, error) {
	item := &dataPageItem{}
	err := item.update(file)
	return item, err
}

func (i *dataPageItem) addItem(item item) {
	switch item := item.(type) {
	case *layoutItem:
		i.dependencies = addItemIfNotExists(i.dependencies, item)
	case *pageItem:
		if item.id == i.dataType {
			i.dependencies = addItemIfNotExists(i.dependencies, item)
		}
	case *staticFileItem:
		i.dependencies = addItemIfNotExists(i.dependencies, item)
	case *dataItem:
		i.dependencies = addItemIfNotExists(i.dependencies, item)
	}
}

func (i *dataPageItem) isDependentOn(item item) bool {
	return containsItem(i.dependencies, item)
}

func (i *dataPageItem) update(file *fs.File) error {
	metadata, tmpl, err := loadPageFile(file)
	if err != nil {
		return fmt.Errorf("error loading page file: %v", err)
	}

	i.dataType = filepath.Base(filepath.Dir(file.Name))
	i.metadata = metadata
	i.tmpl = tmpl
	return nil
}

func (i *dataPageItem) render(ctx *renderContext) error {
	var layout *layoutItem
	var parentPage *pageItem
	var staticFiles []*staticFileItem
	var posts []*dataItem
	var projects []*dataItem
	for _, item := range i.dependencies {
		switch item := item.(type) {
		case *layoutItem:
			layout = item
		case *pageItem:
			parentPage = item
		case *staticFileItem:
			staticFiles = append(staticFiles, item)
		case *dataItem:
			switch item.metadata.(type) {
			case *data.PostMetadata:
				posts = append(posts, item)
			case *data.ProjectMetadata:
				projects = append(projects, item)
			}
		}
	}

	if layout == nil {
		return fmt.Errorf("no layout found")
	}

	staticFileNames := getStaticFileMap(staticFiles)
	render := func(id string, item *dataItem) error {
		outfile := filepath.Join(ctx.out, id+".html")
		pageCtx := dataPageContext{
			ID:          item.id,
			Order:       item.order,
			Metadata:    item.metadata,
			Content:     item.content,
			StaticFiles: staticFileNames,
		}

		renderedMetadata, err := renderPageMetadata(i.metadata, pageCtx)
		if err != nil {
			return fmt.Errorf("error rendering data page %s metadata: %v", id, err)
		}

		err = renderPageToFile(id, outfile, renderedMetadata, i.tmpl, ctx.settings, pageCtx, layout, parentPage, staticFileNames)
		if err != nil {
			return fmt.Errorf("error rendering data page %s: %v", id, err)
		}

		return nil
	}

	switch i.dataType {
	case "posts":
		for _, post := range posts {
			id := i.dataType + "/" + post.id
			if err := render(id, post); err != nil {
				return err
			}
		}
	case "projects":
		for _, project := range projects {
			id := i.dataType + "/" + project.id
			if err := render(id, project); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unknown data type %s", i.dataType)
	}

	return nil
}

func renderPageMetadata(metadata *pageMetadata, context interface{}) (*pageMetadata, error) {
	tmplBytes, err := yaml.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	tmpl, err := txtTemplate.New("").Parse(string(tmplBytes))
	if err != nil {
		return nil, err
	}

	var renderedBytes bytes.Buffer
	err = tmpl.Execute(&renderedBytes, context)
	if err != nil {
		return nil, err
	}

	outMetadata := pageMetadata{}
	err = yaml.Unmarshal(renderedBytes.Bytes(), &outMetadata)
	if err != nil {
		return nil, err
	}

	return &outMetadata, nil
}
