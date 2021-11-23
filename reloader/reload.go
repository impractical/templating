package reloader

// TODO: license compliance
// this file is heavily based on https://gist.github.com/ericlagergren/1e6ea402afc1727ff787

import (
	"fmt"
	"html/template"
	"strings"
	"sync"

	"github.com/pkg/errors"
	fsnotify "gopkg.in/fsnotify.v1"
)

type Templates struct {
	templateDirs []string
	templateExt  string
	templates    map[string]*template.Template
	mu           *sync.RWMutex
	watcher      *fsnotify.Watcher
}

func (t *Templates) Get(name string) *template.Template {
	t.mu.RLock()
	defer t.mu.Unlock()
	if tmpl, ok := t.templates[name]; ok {
		return tmpl
	}
	return nil
}

func (t *Templates) Close() {
	t.watcher.Close()
}

func New(ext string, dirs ...string) (*Templates, error) {
	tmpls := map[string]*template.Template{}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, path := range dirs {
		tmpls[path], err = template.ParseFiles(path)
		if err != nil {
			return nil, errors.Wrapf(err, "error parsing template %s", path)
		}
		watcher.Add(path)
	}

	return &Templates{
		templateDirs: dirs,
		templateExt:  ext,
		templates:    tmpls,
		watcher:      watcher,
		mu:           new(sync.RWMutex),
	}, nil
}

func (t *Templates) Watch() {
	for {
		select {
		case evt, ok := <-t.watcher.Events:
			if !ok {
				return
			}
			if evt.Op&fsnotify.Write == fsnotify.Write {
				fmt.Printf("File: %s Event: %s. Reloading template.\n",
					evt.Name, evt.String())

				if err := t.reload(evt.Name); err != nil {
					fmt.Printf("Error reloading %s: %s\n", evt.Name, err)
				}
			}
		case err := <-t.watcher.Errors:
			fmt.Println("Error in watcher:", err)
		}
	}
}

func (t *Templates) reload(name string) error {
	if len(name) >= len(t.templateExt) &&
		strings.HasSuffix(name, t.templateExt) {

		tmpl := template.Must(template.ParseFiles(name))

		t.mu.Lock()
		t.templates[name] = tmpl
		t.mu.Unlock()

		return nil
	}

	return fmt.Errorf("Unable to reload file %s\n", name)

}
