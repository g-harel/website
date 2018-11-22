package website

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
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
			return fmt.Errorf("could not read template file: %v", err)
		}

		_, err = tmpl.New(strings.TrimPrefix(path, dir+"/")).Parse(string(b))
		if err != nil {
			return fmt.Errorf("could not parse template: %v", err)
		}

		return nil
	}
}

// Render executes all templates in the provided directory using the given data.
// Output is "minified" by trimming the space around each line and current time is prepended.
func Render(dir, entry string, d *Data) (*bytes.Buffer, error) {
	tmpl := template.New("")
	dir = filepath.Clean(dir)

	seen := map[string]bool{}
	tmpl = tmpl.Funcs(template.FuncMap{
		// Used to check if blocks have already been included.
		"first": func(path string) (f bool) {
			f, seen[path] = !seen[path], true
			return
		},
		"rootData": func() *Data {
			return d
		},
	})

	err := filepath.Walk(dir, walker(dir, tmpl))
	if err != nil {
		return nil, fmt.Errorf("could not collect all templates: %v", err)
	}

	original := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(original, entry, d)
	if err != nil {
		return nil, fmt.Errorf("could not execute template: %v", err)
	}

	transformed := &bytes.Buffer{}
	fmt.Fprintf(transformed, "<!-- %v -->\n", time.Now().Format(time.RFC1123))
	for {
		line, readErr := original.ReadBytes('\n')
		if readErr != nil && readErr != io.EOF {
			return nil, fmt.Errorf("could not read line from template output: %v", err)
		}

		_, err = transformed.Write(bytes.TrimSpace(line))
		if err != nil {
			return nil, fmt.Errorf("could not write to transformed template output: %v", err)
		}

		if readErr == io.EOF {
			break
		}
	}

	return transformed, nil
}
