package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
)

// Loader handles configuration file loading and hot-reload
type Loader struct {
	filepath string
	watcher  *fsnotify.Watcher
	stopCh   chan struct{}
}

// NewLoader creates a new configuration loader
func NewLoader(filepath string) (*Loader, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Loader{
		filepath: filepath,
		watcher:  watcher,
		stopCh:   make(chan struct{}),
	}, nil
}

// Load loads the configuration from file
func (l *Loader) Load() error {
	data, err := os.ReadFile(l.filepath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return GetInstance().Load(data)
}

// StartWatching starts watching for configuration file changes
func (l *Loader) StartWatching() error {
	if err := l.watcher.Add(l.filepath); err != nil {
		return fmt.Errorf("failed to watch config file: %w", err)
	}

	go l.watchLoop()
	return nil
}

func (l *Loader) watchLoop() {
	for {
		select {
		case <-l.stopCh:
			return
		case event, ok := <-l.watcher.Events:
			if !ok {
				return
			}

			// Reload config on write or create events
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				if err := l.Load(); err != nil {
					fmt.Printf("Error reloading config: %v\n", err)
				}
			}

		case err, ok := <-l.watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("Watcher error: %v\n", err)
		}
	}
}

// Stop stops watching for changes
func (l *Loader) Stop() error {
	close(l.stopCh)
	return l.watcher.Close()
}
