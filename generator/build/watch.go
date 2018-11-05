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

	for _, file := range files {
		b.handleFile(file)
	}

	log.Printf("Built site. Now watching for changes (press Ctrl+C to stop)...\n")

	watcher.Run()
	defer watcher.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

out:
	for {
		select {
		case file := <-watcher.Files:
			log.Printf("File changed: %s\n", file)
			b.handleFile(file)
		case sig := <-sigs:
			log.Printf("Received signal: %v\n", sig)
			break out
		}
	}
}
