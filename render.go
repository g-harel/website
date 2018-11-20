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

	original := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(original, entry, d)
	if err != nil {
		return nil, fmt.Errorf("could not execute template: %s", err)
	}

	transformed := &bytes.Buffer{}
	fmt.Fprintf(transformed, "<!-- %v -->\n", time.Now().Format(time.RFC1123))
	for {
		line, readErr := original.ReadBytes('\n')
		if readErr != nil && readErr != io.EOF {
			return nil, fmt.Errorf("could not read line from template output: %s", err)
		}

		_, err = transformed.Write(bytes.TrimSpace(line))
		if err != nil {
			return nil, fmt.Errorf("could not write to transformed template output: %s", err)
		}

		if readErr == io.EOF {
			break
		}
	}

	return transformed, nil
}
