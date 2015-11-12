// Heavily inspired by github.com/nathany/looper.
package main

import (
	"github.com/chinthakagodawita/docker-unisync/Godeps/_workspace/src/gopkg.in/fsnotify.v1"
	"os"
	"path/filepath"
	"strings"
)

type Watcher struct {
	*fsnotify.Watcher
	Files   chan string
	Folders chan string
}

func NewWatcher(path string) (*Watcher, error) {
	folders := SubFolders(path)

	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	defer fsWatcher.Close()

	watcher := &Watcher{Watcher: fsWatcher}

	watcher.Files = make(chan string, 10)
	watcher.Folders = make(chan string, len(folders))

	for _, dir := range folders {
		LogInfo(dir)
		if err = watcher.AddFolder(dir); err != nil {
			return nil, err
		}
	}

	return watcher, nil
}

func (watcher *Watcher) AddFolder(folder string) error {
	if err := watcher.Add(folder); err != nil {
		return err
	}
	watcher.Folders <- folder
	return nil
}

func (watcher *Watcher) Run() {
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					info, err := os.Stat(event.Name)
					if err != nil {
						// eg. stat .subl513.tmp : no such file or directory
						LogError(err.Error())
					} else if info.IsDir() {
						if !isIgnored(filepath.Base(event.Name)) {
							watcher.AddFolder(event.Name)
						}
					} else {
						watcher.Files <- event.Name
					}
				}

				LogInfo(event.Name)

				if event.Op&fsnotify.Write == fsnotify.Write {
					watcher.Files <- event.Name
				}
			}
		}
	}()
}

func SubFolders(path string) (paths []string) {
	filepath.Walk(path, func(curPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			dirName := info.Name()

			if isIgnored(dirName) {
				return filepath.SkipDir
			}

			// fmt.Println(curPath)
			paths = append(paths, curPath)
		}

		return nil
	})

	return paths
}

func isIgnored(name string) bool {
	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")
}
