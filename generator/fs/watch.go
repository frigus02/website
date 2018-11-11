package fs

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// RecursiveWatcher watches the given folder and all subfolders for file and
// folder changes and notifies about them through channels.
type RecursiveWatcher struct {
	*fsnotify.Watcher
	basepath string
	Files    chan string
}

// NewRecursiveWatcher creates a new watcher for the given path and returns
// both the watcher and the initial files list.
func NewRecursiveWatcher(path string) (rw *RecursiveWatcher, files []string, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return nil, nil, err
	}

	path = path + string(filepath.Separator)

	files, folders, err := ListRecursive(path)
	if err != nil {
		return nil, nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, err
	}

	rw = &RecursiveWatcher{Watcher: watcher, basepath: path}
	rw.Files = make(chan string, 10)

	if err = rw.Add(path); err != nil {
		rw.Close()
		return nil, nil, err
	}

	for _, folder := range folders {
		if err = rw.Add(filepath.Join(path, folder)); err != nil {
			rw.Close()
			return nil, nil, err
		}
	}

	return rw, files, nil
}

// Run starts the watching.
func (watcher *RecursiveWatcher) Run() {
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fi, err := os.Stat(event.Name)
				if err != nil {
					log.Printf("Got creation event for %s, but could not find file: %v\n", event.Name, err)
					break
				}

				strippedPath := strings.TrimPrefix(event.Name, watcher.basepath)

				// File or directory was created
				if event.Op&fsnotify.Create == fsnotify.Create {
					if fi.IsDir() {
						err := watcher.Add(event.Name)
						if err != nil {
							log.Printf("Error watching newly created folder %s: %v\n", event.Name, err)
						}
					} else {
						watcher.Files <- strippedPath
					}
				}

				// File or directory was modified
				if event.Op&fsnotify.Write == fsnotify.Write {
					if !fi.IsDir() {
						watcher.Files <- strippedPath
					}
				}

			case err := <-watcher.Errors:
				log.Printf("Error watching files: %v\n", err)
			}
		}
	}()
}