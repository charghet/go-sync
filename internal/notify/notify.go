package notify

import (
	"os"
	"path/filepath"
	"time"

	"github.com/charghet/go-sync/pkg/logger"
	"github.com/fsnotify/fsnotify"
)

type Notify struct {
	watcher *fsnotify.Watcher
	Events  chan fsnotify.Event
	Errors  chan error
}

func NewNotify() (*Notify, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal("Failed to create fsnotify watcher:", err)
		return nil, err
	}
	return &Notify{watcher: watcher}, nil
}

func (n *Notify) Add(p string) error {
	info, err := os.Stat(p)
	if err != nil {
		logger.Warn("Failed to stat path:", p, "Error:", err)
		return err
	}

	if info.IsDir() {
		err = addRecursiveWatch(n.watcher, p)
		if err != nil {
			logger.Danger("Failed to add recursive watch for directory:", p, "Error:", err)
			return err
		}
	} else {
		err = n.watcher.Add(p)
		if err != nil {
			logger.Danger("Failed to add watcher:", p, "Error:", err)
			return err
		}
		logger.Info("Added watcher for file:", p)
	}

	n.Events = make(chan fsnotify.Event, 10)
	n.Errors = make(chan error, 10)
	go func() {
		for {
			select {
			case event, ok := <-n.watcher.Events:
				if !ok {
					return
				}
				if filepath.Base(event.Name) == ".git" {
					continue
				}
				logger.Info("Notify Received event:", event, "for path:", event.Name)
				if event.Op&fsnotify.Create == fsnotify.Create {
					time.Sleep(100 * time.Millisecond)
					info, err := os.Stat(event.Name)
					if err != nil {
						logger.Warn("Failed to stat created path:", event.Name, "Error:", err)
						continue
					}
					if info.IsDir() {
						err = addRecursiveWatch(n.watcher, event.Name)
						if err != nil {
							logger.Danger("Failed to add recursive watch for created directory:", event.Name, "Error:", err)
							continue
						}
						logger.Info("Added recursive watch for created directory:", event.Name)
					}
				}
				n.Events <- event
			case err, ok := <-n.watcher.Errors:
				if !ok {
					return
				}
				n.Errors <- err
			}
		}
	}()

	return nil
}

func (n *Notify) Close() error {
	if n.watcher != nil {
		err := n.watcher.Close()
		if err != nil {
			logger.Warn("Failed to close watcher:", err)
		}
	}
	return nil
}

func addRecursiveWatch(watcher *fsnotify.Watcher, root string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Warn("Failed to walk path:", path, "Error:", err)
			return nil // Continue walking
		}
		logger.Info("Walking path:", path)
		if info.IsDir() {
			if filepath.Base(path) == ".git" {
				return filepath.SkipDir
			}

			err = watcher.Add(path)
			if err != nil {
				logger.Danger("Failed to add watcher for directory:", path, "Error:", err)
				return err
			} else {
				logger.Info("Added watcher for directory:", path)
			}
		}
		return nil
	})
	if err != nil {
		logger.Danger("Failed to walk directory:", root, "Error:", err)
	}
	return err
}
