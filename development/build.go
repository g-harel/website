package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/g-harel/website"
	_ "github.com/joho/godotenv/autoload"
)

var env = struct {
	RebuildBatch    int
	ConfigPath      string
	TemplateDir     string
	TemplateEntry   string
	GraphQLEndpoint string
	GraphQLToken    string
	OutputFile      string
}{
	RebuildBatch:    123,
	ConfigPath:      ".config",
	TemplateDir:     "templates",
	TemplateEntry:   "entry.html",
	GraphQLEndpoint: "https://api.github.com/graphql",
	GraphQLToken:    os.Getenv("GRAPHQL_TOKEN"),
	OutputFile:      "index.html",
}

// Fatal errors are formatted/printed to stderr and the process exits.
func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "\033[31;1m\nerror: %v\033[0m\n\n", fmt.Sprintf(format, a...))
	os.Exit(1)
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatal("could not create file watcher: %s", err)
	}
	defer watcher.Close()

	// Recursively walk the templates directory to watch the entire tree.
	err = filepath.Walk(env.TemplateDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				return fmt.Errorf("could not add path to watcher: %s", err)
			}
		}

		return nil
	})
	if err != nil {
		fatal("could not recursively watch files in \"%v\": %v", env.TemplateDir, err)
	}
	fmt.Printf("watching \"%s\"\n", env.TemplateDir)

	// Poll rebuild variable to batch multiple watcher events.
	rebuild := true
	go func() {
		for {
			if rebuild {
				rebuild = false
				fmt.Printf("> ")
				err := Build()
				if err != nil {
					fmt.Printf("build failed: %v\n", err)
				}
			}
			time.Sleep(time.Duration(env.RebuildBatch) * time.Millisecond)
		}
	}()

	// Listen for manual rebuilds.
	go func() {
		r := bufio.NewReader(os.Stdin)
		for {
			s, err := r.ReadString('\n')
			if err != nil {
				fatal("could not read user input: %v", err)
			}
			if s == ".\n" {
				rebuild = true
			}
		}
	}()

	// React to events/errors from watcher channels.
	for {
		select {
		case _ = <-watcher.Events:
			rebuild = true
		case err := <-watcher.Errors:
			fatal("watcher error: %v", err)
		}
	}

}

// Build renders output from local config to local file, API data is queried on every call.
func Build() error {
	start := time.Now()

	configContent, err := ioutil.ReadFile(env.ConfigPath)
	if err != nil {
		return fmt.Errorf("could not load local config")
	}

	config, err := (&website.Config{}).Parse(string(configContent))
	if err != nil {
		return fmt.Errorf("could not parse local config: %v", err)
	}

	query, err := config.Query()
	if err != nil {
		return fmt.Errorf("could not generate query from config: %v", err)
	}

	dataReq, err := http.NewRequest("POST", env.GraphQLEndpoint, bytes.NewReader(query))
	if err != nil {
		return fmt.Errorf("could not create data request: %v", err)
	}
	dataReq.Header.Add("Authorization", fmt.Sprintf("bearer %v", env.GraphQLToken))

	dataRes, err := (&http.Client{}).Do(dataReq)
	if err != nil {
		return fmt.Errorf("request for data failed")
	}

	dataBody := &bytes.Buffer{}
	_, err = io.Copy(dataBody, dataRes.Body)
	if err != nil {
		return fmt.Errorf("could not read data response: %v", err)
	}

	data, err := (&website.Data{}).Parse(dataBody.String())
	if err != nil {
		return fmt.Errorf("could not parse received data: %v", err)
	}

	// TODO add to build function too.
	for i := 0; i < len(config.Creations); i++ {
		data.Creations = append(data.Creations, &website.CreationData{
			ImageURL: config.Creations[i].ImageURL,
		})
	}

	output, err := website.Render(env.TemplateDir, env.TemplateEntry, data)
	if err != nil {
		return fmt.Errorf("could not render templates: %v", err)
	}

	err = ioutil.WriteFile(env.OutputFile, output.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("could not write rendered output to file: %v", err)
	}

	fmt.Printf("built in %v\n", time.Since(start))

	return nil
}
