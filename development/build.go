package main

import (
	"bufio"
	"context"
	"fmt"
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
	if env.GraphQLToken == "" {
		fatal("missing GRAPHQL_TOKEN")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatal("could not create file watcher: %s", err)
	}
	defer watcher.Close()

	// Watch for config changes.
	watcher.Add(env.ConfigPath)

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
	fmt.Printf("watching \"%s\" and \"%s\"\n", env.TemplateDir, env.ConfigPath)

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

	configContent, err := os.ReadFile(env.ConfigPath)
	if err != nil {
		return fmt.Errorf("could not load local config")
	}

	config, err := (&website.Config{}).Parse(string(configContent))
	if err != nil {
		return fmt.Errorf("could not parse local config: %v", err)
	}

	data, err := config.FetchData(context.Background(), env.GraphQLEndpoint, env.GraphQLToken)
	if err != nil {
		return err
	}

	output, err := website.Render(env.TemplateDir, env.TemplateEntry, data)
	if err != nil {
		return fmt.Errorf("could not render templates: %v", err)
	}

	err = os.WriteFile(env.OutputFile, output.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("could not write rendered output to file: %v", err)
	}

	fmt.Printf("built in %v\n", time.Since(start))

	return nil
}
