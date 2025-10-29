package fileio

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	fsWatcher		*fsnotify.Watcher
	mu				sync.Mutex
	done			chan struct{}
	Events			chan string
}

// NewWatcher initializes a new file watcher
func NewWatcher(rootPath string) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	wr := &Watcher{
		fsWatcher: w,
		done: make(chan struct{}),
		Events: make(chan string, 32),
	}

	// Watch recursively
	err = filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return w.Add(path)
		}
		return nil
	})
	if err != nil {
		wr.Close()
		return nil, err
	}

	go wr.run()
	return wr, nil
}

func (w *Watcher) run() {
	for {
		select {
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}
			w.Events <- event.Name
		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			_ = err // ignored
		case <-w.done:
			return
		}
	}
}

// Close stops the watcher
func (w *Watcher) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	select {
	case <-w.done:
	default:
		close(w.done)
		w.fsWatcher.Close()
		close(w.Events)
	}
}
