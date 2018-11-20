package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/g-harel/website"
)

var env = struct {
	ConfigPath      string
	TemplateDir     string
	TemplateEntry   string
	GraphQLEndpoint string
	GraphQLToken    string
	OutputFile      string
}{
	ConfigPath:      ".config",
	TemplateDir:     "templates",
	TemplateEntry:   "index.tmpl",
	GraphQLEndpoint: "https://api.github.com/graphql",
	GraphQLToken:    os.Getenv("GRAPHQL_TOKEN"),
	OutputFile:      "index.html",
}

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "\033[31;1m\nerror: %s\033[0m\n\n", fmt.Sprintf(format, a...))
	os.Exit(1)
}

func main() {
	start := time.Now()

	configContent, err := ioutil.ReadFile(env.ConfigPath)
	if err != nil {
		fatal("could not load local config")
	}

	config, err := (&website.Config{}).Parse(string(configContent))
	if err != nil {
		fatal("could not parse local config: %v", err)
	}

	query, err := config.Query()
	if err != nil {
		fatal("could not generate query from config: %s", err)
	}

	dataReq, err := http.NewRequest("POST", env.GraphQLEndpoint, bytes.NewReader(query))
	if err != nil {
		fatal("could not create data request: %v", err)
	}
	dataReq.Header.Add("Authorization", fmt.Sprintf("bearer %v", env.GraphQLToken))

	dataRes, err := (&http.Client{}).Do(dataReq)
	if err != nil {
		fatal("request for data failed")
	}

	dataBody := &bytes.Buffer{}
	_, err = io.Copy(dataBody, dataRes.Body)
	if err != nil {
		fatal("could not read data response: %v", err)
	}

	data, err := (&website.Data{}).Parse(dataBody.String())
	if err != nil {
		fatal("could not parse received data: %v", err)
	}

	output, err := website.Render(env.TemplateDir, env.TemplateEntry, data)
	if err != nil {
		fatal("could not render templates: %v", err)
	}

	err = ioutil.WriteFile(env.OutputFile, output.Bytes(), 0644)
	if err != nil {
		fatal("could not write rendered output to file: %s", err)
	}

	fmt.Printf("\nbuilt in %s\n\n", time.Since(start))
}
