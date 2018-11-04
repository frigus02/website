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
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/frigus02/website/generator/data"
	"github.com/frigus02/website/generator/fs"
)

type pageMetadata struct {
	Title string `yaml:"title"`
}

// Watch builds the website and watches files for changes.
type Watch struct {
	In             string
	Out            string
	stylesheets    []string
	stylesheetName string
	layoutTemplate *template.Template
	seenPageFiles  []string
	pageContext    pageContext
}

// Watch starts the watching.
func (w *Watch) Watch() {
	watcher, files, err := fs.NewRecursiveWatcher(w.In)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		w.handleFile(file)
	}

	watcher.Run()
	defer watcher.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

out:
	for {
		select {
		case file := <-watcher.Files:
			log.Printf("File changed: %s\n", file)
			w.handleFile(file)
		case sig := <-sigs:
			log.Printf("Received signal: %v\n", sig)
			break out
		}
	}
}

func (w *Watch) handleFile(file string) {
	topDir := file[0:strings.Index(file, string(filepath.Separator))]
	if topDir == "data" {
		w.updateDataItem(file)
	} else if topDir == "pages" {
		log.Printf("Pages file %s\n", file)
		w.updatePageFile(file)
	} else if topDir == "static" {
		w.updateStaticFile(file)
	} else {
		log.Printf("Unknown file: %s\n", file)
	}
}

func (w *Watch) invalidateSeenPageFiles() {
	for _, file := range w.seenPageFiles {
		w.handleFile(file)
	}
}

func (w *Watch) trackSeenPageFiles(file string) {
	found := false
	for _, seenPageFile := range w.seenPageFiles {
		if seenPageFile == file {
			found = true
		}
	}

	if !found {
		w.seenPageFiles = append(w.seenPageFiles, file)
	}
}

func (w *Watch) updateDataItem(file string) {
	item, err := data.GetItem(filepath.Join(w.In, file))
	if err != nil {
		log.Printf("Error getting data item for file %s: %v\n", file, err)
		return
	}

	switch item := item.(type) {
	case *data.Post:
		for i, post := range w.pageContext.Posts {
			if post.ID == item.ID {
				w.pageContext.Posts[i] = item
				return
			}
		}

		w.pageContext.Posts = append(w.pageContext.Posts, item)
	case *data.Project:
		for i, project := range w.pageContext.Projects {
			if project.ID == item.ID {
				w.pageContext.Projects[i] = item
				return
			}
		}

		w.pageContext.Projects = append(w.pageContext.Projects, item)
	default:
		log.Fatal("Unexpected data item type")
	}

	w.invalidateSeenPageFiles()
}

func (w *Watch) updatePageFile(file string) {
	if file[6:] == "_layout.html" {
		tmpl, err := template.ParseFiles(filepath.Join(w.In, file))
		if err != nil {
			log.Printf("Error parsing _layout.html template: %v\n", err)
			return
		}

		w.layoutTemplate = tmpl
		w.invalidateSeenPageFiles()
	} else if filepath.Base(file) == "_details.html" {
		log.Printf("Details template %s\n", file)
	} else {
		metadata := pageMetadata{}
		content, err := fs.ReadFileWithMetadata(filepath.Join(w.In, file), &metadata)
		if err != nil {
			log.Printf("Error reading file with metadata %s: %v\n", file, err)
			return
		}

		tmpl, err := template.New("").Parse(content)
		if err != nil {
			log.Printf("Error parsing page template %s: %v\n", file, err)
			return
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, w.pageContext)
		if err != nil {
			log.Printf("Error executing template %s: %v\n", file, err)
			return
		}

		if w.layoutTemplate != nil {
			// TODO: fill ParentID and ParentTitle

			id := strings.TrimSuffix(file[6:], filepath.Ext(file))
			id = strings.TrimSuffix(id, string(filepath.Separator)+"index")

			layoutContext := layoutContext{
				ID:          id,
				Title:       metadata.Title,
				Content:     template.HTML(buf.String()),
				ParentID:    "index",
				ParentTitle: "",
				Stylesheet:  w.stylesheetName,
			}

			destinationFileName := filepath.Join(w.Out, file[6:])
			err = os.MkdirAll(filepath.Dir(destinationFileName), 0644)
			if err != nil {
				log.Printf("Error creating destination folder %s for page %s: %v\n", destinationFileName, file, err)
				return
			}

			destination, err := os.Create(destinationFileName)
			if err != nil {
				log.Printf("Error creating destination %s for page %s: %v\n", destinationFileName, file, err)
				return
			}

			defer destination.Close()

			err = w.layoutTemplate.Execute(destination, &layoutContext)
			if err != nil {
				log.Printf("Error executing layout template for %s: %v\n", file, err)
				return
			}
		}

		w.trackSeenPageFiles(file)
	}
}

func (w *Watch) updateStaticFile(file string) {
	if filepath.Ext(file) == ".css" {
		w.updateStylesheet(file)
		return
	}

	dest := filepath.Join(w.Out, file[7:])

	source, err := os.Open(filepath.Join(w.In, file))
	if err != nil {
		log.Printf("Error reading static file %s: %v\n", file, err)
		return
	}

	defer source.Close()

	err = os.MkdirAll(filepath.Dir(dest), 0644)
	if err != nil {
		log.Printf("Error creating destination folder %s for static file %s: %v\n", dest, file, err)
		return
	}

	destination, err := os.Create(dest)
	if err != nil {
		log.Printf("Error creating destination %s for static file %s: %v\n", dest, file, err)
		return
	}

	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		log.Printf("Error copying static file %s: %v\n", file, err)
		return
	}
}

func (w *Watch) updateStylesheet(file string) {
	// Add stylesheet to list if not yet present
	found := false
	for _, stylesheet := range w.stylesheets {
		if stylesheet == file {
			found = true
		}
	}

	if !found {
		w.stylesheets = append(w.stylesheets, file)
	}

	// Output concatenated stylesheet.
	err := os.MkdirAll(w.Out, 0644)
	if err != nil {
		log.Printf("Error creating destination folder %s for stylesheet: %v\n", w.Out, err)
		return
	}

	var content bytes.Buffer
	hash := md5.New()
	hashAndFile := io.MultiWriter(hash, &content)

	for _, stylesheet := range w.stylesheets {
		source, err := os.Open(filepath.Join(w.In, stylesheet))
		if err != nil {
			log.Printf("Error reading stylesheet %s: %v\n", stylesheet, err)
			continue
		}

		defer source.Close()

		_, err = io.Copy(hashAndFile, source)
		if err != nil {
			log.Printf("Error copying stylesheet %s: %v\n", stylesheet, err)
			continue
		}
	}

	filenameWithHash := fmt.Sprintf("styles-%x.css", hash.Sum(nil)[:8])
	err = ioutil.WriteFile(filepath.Join(w.Out, filenameWithHash), content.Bytes(), 0644)
	if err != nil {
		log.Printf("Error creating new stylesheet: %v\n", err)
		return
	}

	if w.stylesheetName != "" {
		err = os.Remove(filepath.Join(w.Out, w.stylesheetName))
		if err != nil {
			log.Printf("Error removing old stylesheet: %v\n", err)
			return
		}
	}

	w.stylesheetName = filenameWithHash
	w.invalidateSeenPageFiles()
}
