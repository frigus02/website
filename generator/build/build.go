package build

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	txtTemplate "text/template"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	yaml "gopkg.in/yaml.v2"

	"github.com/frigus02/website/generator/data"
	"github.com/frigus02/website/generator/fs"
)

type pageMetadata struct {
	Title string `yaml:"title"`
}

// Build builds the website and is able to continuously update the output.
type Build struct {
	In             string
	Out            string
	Minify         bool
	minifier       *minify.M
	stylesheets    []string
	stylesheetName string
	layoutTemplate *template.Template
	seenPageFiles  []string
	pageContext    pageContext
}

// Build builds the site once and returns.
func (b *Build) Build() {
	files, _, err := fs.ListRecursive(b.In)
	if err != nil {
		log.Fatal(err)
	}

	b.initMinifier()

	for _, file := range files {
		b.handleFile(file)
	}
}

func (b *Build) initMinifier() {
	if b.Minify {
		b.minifier = minify.New()
		b.minifier.AddFunc("text/css", css.Minify)
		b.minifier.AddFunc("text/html", html.Minify)
	}
}

func (b *Build) handleFile(file string) {
	topDir := file[0:strings.Index(file, string(filepath.Separator))]
	if topDir == "data" {
		b.updateDataItem(file)
	} else if topDir == "pages" {
		b.updatePageFile(file)
	} else if topDir == "static" {
		b.updateStaticFile(file)
	} else {
		log.Printf("Unknown file: %s\n", file)
	}
}

func (b *Build) trackSeenPageFiles(file string) {
	found := false
	for _, seenPageFile := range b.seenPageFiles {
		if seenPageFile == file {
			found = true
		}
	}

	if !found {
		b.seenPageFiles = append(b.seenPageFiles, file)
	}
}

func (b *Build) invalidateSeenPageFiles() {
	for _, file := range b.seenPageFiles {
		b.handleFile(file)
	}
}

func (b *Build) updateDataItem(file string) {
	if filepath.Ext(file) == ".png" {
		_, typeDir, _, fileName, id := data.ExtractMetadataFromFilePath(file)

		src := filepath.Join(b.In, file)
		dst := filepath.Join(b.Out, "images", typeDir, id, fileName)

		err := fs.CopyFile(src, dst)
		if err != nil {
			log.Printf("Error copying static data file %s: %v\n", file, err)
			return
		}
	} else {
		item, err := data.GetItem(filepath.Join(b.In, file))
		if err != nil {
			log.Printf("Error getting data item for file %s: %v\n", file, err)
			return
		}

		switch item := item.(type) {
		case *data.Post:
			for i, post := range b.pageContext.Posts {
				if post.ID == item.ID {
					b.pageContext.Posts[i] = item
					return
				}
			}

			b.pageContext.Posts = append(b.pageContext.Posts, item)
		case *data.Project:
			for i, project := range b.pageContext.Projects {
				if project.ID == item.ID {
					b.pageContext.Projects[i] = item
					return
				}
			}

			b.pageContext.Projects = append(b.pageContext.Projects, item)
		default:
			log.Fatal("Unexpected data item type")
		}

		b.invalidateSeenPageFiles()
	}
}

func (b *Build) updatePageFile(file string) {
	if file[6:] == "_layout.html" {
		tmpl, err := template.ParseFiles(filepath.Join(b.In, file))
		if err != nil {
			log.Printf("Error parsing _layout.html template: %v\n", err)
			return
		}

		b.layoutTemplate = tmpl
		b.invalidateSeenPageFiles()
	} else if filepath.Base(file) == "_details.html" {
		metadata, tmpl, err := b.loadPageFile(file)
		if err != nil {
			log.Printf("Error loading page file %s: %v\n", file, err)
			return
		}

		dataType := filepath.Base(filepath.Dir(file))
		switch dataType {
		case "posts":
			for _, post := range b.pageContext.Posts {
				id := dataType + "/" + post.ID

				renderedMetadata, err := renderPageMetadata(metadata, post)
				if err != nil {
					log.Printf("Error rendering data page %s metadata: %v\n", id, err)
					return
				}

				err = b.renderPageToFile(id, renderedMetadata, tmpl, post)
				if err != nil {
					log.Printf("Error rendering data page %s: %v\n", id, err)
					return
				}
			}
		case "projects":
			for _, project := range b.pageContext.Projects {
				id := dataType + "/" + project.ID

				renderedMetadata, err := renderPageMetadata(metadata, project)
				if err != nil {
					log.Printf("Error rendering data page %s metadata: %v\n", id, err)
					return
				}

				err = b.renderPageToFile(id, renderedMetadata, tmpl, project)
				if err != nil {
					log.Printf("Error rendering data page %s: %v\n", id, err)
					return
				}
			}
		default:
			log.Printf("Unknown data type %s for file %s\n", dataType, file)
			return
		}

		b.trackSeenPageFiles(file)
	} else {
		metadata, tmpl, err := b.loadPageFile(file)
		if err != nil {
			log.Printf("Error loading page file %s: %v\n", file, err)
			return
		}

		id := strings.TrimSuffix(file[6:], filepath.Ext(file))
		id = strings.TrimSuffix(id, string(filepath.Separator)+"index")

		err = b.renderPageToFile(id, metadata, tmpl, &b.pageContext)
		if err != nil {
			log.Printf("Error rendering page %s: %v\n", file, err)
			return
		}

		b.trackSeenPageFiles(file)
	}
}

func (b *Build) updateStaticFile(file string) {
	if filepath.Ext(file) == ".css" {
		b.updateStylesheet(file)
		return
	}

	src := filepath.Join(b.In, file)
	dst := filepath.Join(b.Out, file[7:])

	err := fs.CopyFile(src, dst)
	if err != nil {
		log.Printf("Error copying static file %s: %v\n", file, err)
		return
	}
}

func (b *Build) updateStylesheet(file string) {
	// Add stylesheet to list if not yet present
	found := false
	for _, stylesheet := range b.stylesheets {
		if stylesheet == file {
			found = true
		}
	}

	if !found {
		b.stylesheets = append(b.stylesheets, file)
	}

	// Output concatenated stylesheet.
	err := os.MkdirAll(b.Out, 0644)
	if err != nil {
		log.Printf("Error creating destination folder %s for stylesheet: %v\n", b.Out, err)
		return
	}

	var content bytes.Buffer
	hash := md5.New()
	writer := newNopWriteCloser(io.MultiWriter(hash, &content))

	if b.Minify {
		writer = b.minifier.Writer("text/css", writer)
	}

	for _, stylesheet := range b.stylesheets {
		source, err := os.Open(filepath.Join(b.In, stylesheet))
		if err != nil {
			log.Printf("Error reading stylesheet %s: %v\n", stylesheet, err)
			continue
		}

		defer source.Close()

		_, err = io.Copy(writer, source)
		if err != nil {
			log.Printf("Error copying stylesheet %s: %v\n", stylesheet, err)
			continue
		}
	}

	err = writer.Close()
	if err != nil {
		log.Printf("Error closing (minify) writer for stylesheets: %v\n", err)
		return
	}

	filenameWithHash := fmt.Sprintf("styles-%x.css", hash.Sum(nil)[:8])
	err = ioutil.WriteFile(filepath.Join(b.Out, filenameWithHash), content.Bytes(), 0644)
	if err != nil {
		log.Printf("Error creating new stylesheet: %v\n", err)
		return
	}

	if b.stylesheetName != "" {
		err = os.Remove(filepath.Join(b.Out, b.stylesheetName))
		if err != nil {
			log.Printf("Error removing old stylesheet: %v\n", err)
			return
		}
	}

	b.stylesheetName = filenameWithHash
	b.invalidateSeenPageFiles()
}

func (b *Build) loadPageFile(file string) (*pageMetadata, *template.Template, error) {
	metadata := pageMetadata{}
	content, err := fs.ReadFileWithMetadata(filepath.Join(b.In, file), &metadata)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading file with metadata %s: %v", file, err)
	}

	tmpl, err := template.New("").Parse(content)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing page template %s: %v", file, err)
	}

	return &metadata, tmpl, nil
}

func (b *Build) renderPageToFile(
	id string,
	metadata *pageMetadata,
	tmpl *template.Template,
	pageContext interface{},
) error {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, pageContext)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", id, err)
	}

	if b.layoutTemplate != nil {
		// TODO: fill ParentID and ParentTitle

		layoutContext := layoutContext{
			ID:          id,
			Title:       metadata.Title,
			Content:     template.HTML(buf.String()),
			ParentID:    "index",
			ParentTitle: "",
			Stylesheet:  b.stylesheetName,
		}

		destinationFileName := filepath.Join(b.Out, id+".html")
		err = os.MkdirAll(filepath.Dir(destinationFileName), 0644)
		if err != nil {
			return fmt.Errorf("error creating destination folder %s for page %s: %v", destinationFileName, id, err)
		}

		var destination io.WriteCloser
		destination, err = os.Create(destinationFileName)
		if err != nil {
			return fmt.Errorf("error creating destination %s for page %s: %v", destinationFileName, id, err)
		}

		if b.Minify {
			destination = b.minifier.Writer("text/html", destination)
		}

		err = b.layoutTemplate.Execute(destination, &layoutContext)
		if err != nil {
			destination.Close()
			return fmt.Errorf("error executing layout template for %s: %v", id, err)
		}

		err = destination.Close()
		if err != nil {
			return fmt.Errorf("error closing (minify) writer for %s: %v", id, err)
		}
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
