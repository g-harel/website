package website

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Generator for a file tree walker which recursively collects files into a template.
// Templates are named relative to the root directory.
func walker(dir string, tmpl *template.Template) func(string, os.FileInfo, error) error {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read template file: %s", err)
		}

		_, err = tmpl.New(strings.TrimPrefix(path, dir+"/")).Parse(string(b))
		if err != nil {
			return fmt.Errorf("could not parse template: %s", err)
		}

		return nil
	}
}

func Render(dir, entry string, d *Data) (*bytes.Buffer, error) {
	tmpl := template.New("")
	dir = filepath.Clean(dir)

	err := filepath.Walk(dir, walker(dir, tmpl))
	if err != nil {
		return nil, fmt.Errorf("could not collect all templates: %s", err)
	}

	b := bytes.Buffer{}
	err = tmpl.ExecuteTemplate(&b, entry, d)
	if err != nil {
		return nil, fmt.Errorf("could not execute template: %s", err)
	}

	return &b, nil
}
