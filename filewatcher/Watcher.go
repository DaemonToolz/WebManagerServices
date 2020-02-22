package main

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher

func initFileWatcher(root string) {
	watcher, _ = fsnotify.NewWatcher()

	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk(root, watchDir); err != nil {
		failOnError(err, "Couldn't get to the folder")
	}

}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}
