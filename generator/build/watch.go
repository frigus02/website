package build

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/frigus02/website/generator/fs"
)

// Watch builds the site once and then watches files for changes and updates
// the site. Stops when receiving SIGINT or SIGTERM.
func (b *Build) Watch() {
	watcher, files, err := fs.NewRecursiveWatcher(b.In)
	if err != nil {
		log.Fatal(err)
	}

	if b.ConcatCSS {
		b.fileReader = newCSSConcatFileReader()
	} else {
		b.fileReader = newSimpleFileReader()
	}

	b.renderCtx = newRenderContext(b.Out)
	b.items = make(map[string]item)
	if b.Minify {
		b.renderCtx.settings.minifier = newMinifier()
	}

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

	log.Printf("Watching for changes (press Ctrl+C to stop)...\n")

	watcher.Run()
	defer watcher.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

out:
	for {
		select {
		case file := <-watcher.Files:
			log.Printf("File changed: %s\n", file)
			err = b.handleFile(file)
			if err != nil {
				log.Printf("Error: %v\n", err)
			} else {
				err = b.render()
				if err != nil {
					log.Printf("Error: %v\n", err)
				}
			}
		case sig := <-sigs:
			log.Printf("Received signal: %v\n", sig)
			break out
		}
	}
}
