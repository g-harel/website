package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
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

// Render generates the static files from the config and API data.
func Render(c *Config, r *Response) error {
	tmpl := template.New("")
	dir := filepath.Clean(c.Templates)

	err := filepath.Walk(dir, walker(dir, tmpl))
	if err != nil {
		return fmt.Errorf("could not collect all templates: %s", err)
	}

	b := bytes.Buffer{}
	err = tmpl.ExecuteTemplate(&b, c.RootName, r)
	if err != nil {
		return fmt.Errorf("could not execute template: %s", err)
	}

	f, err := os.Create(c.OutPath)
	if err != nil {
		return fmt.Errorf("could not create output file: %s", err)
	}

	for {
		line, readErr := b.ReadBytes('\n')
		if readErr != nil && readErr != io.EOF {
			return fmt.Errorf("could not read line from executed template: %s", err)
		}

		_, err = f.Write(bytes.TrimSpace(line))
		if err != nil {
			return fmt.Errorf("could not write to output file: %s", err)
		}

		if readErr == io.EOF {
			break
		}
	}

	return nil
}
