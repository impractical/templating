package compiled

// this should be used in conjunction with

import (
	"html/template"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Templates struct {
	templates map[string]*template.Template
}

func (t *Templates) Get(name string) *template.Template {
	if tmpl, ok := t.templates[name]; ok {
		return tmpl
	}
	return nil
}

func New(ext string, asset func(path string) ([]byte, error), assetDir func(path string) ([]string, error), dirs ...string) (*Templates, error) {
	tmpls := map[string]*template.Template{}
	for _, dir := range dirs {
		// TODO: assetDir isn't recursive, it'll return only one level of results
		// need to fix that
		files, err := assetDir(dir)
		if err != nil {
			return nil, err
		}

		for _, name := range files {
			if strings.HasSuffix(name, ext) {
				path := filepath.Join(dir, name)
				file, err := asset(path)
				if err != nil {
					return nil, err
				}
				tmpls[path], err = template.New(filepath.Base(path)).Parse(string(file))
				if err != nil {
					return nil, errors.Wrapf(err, "error parsing template %s", path)
				}

			}
		}
	}

	return &Templates{
		templates: tmpls,
	}, nil
}
