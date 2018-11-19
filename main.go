package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "\033[31;1mfatal error: %s\033[0m\n", fmt.Sprintf(format, a...))
	os.Exit(1)
}

func main() {
	start := time.Now()

	configReq, err := http.Get(os.Getenv("CONFIG_SRC"))
	if err != nil {
		fatal("could not fetch config: %v", err)
	}

	configData := &bytes.Buffer{}
	_, err = io.Copy(configData, configReq.Body)
	if err != nil {
		fatal("could not read config response data: %v", err)
	}

	config := (&Config{}).Unmarshall(configData.Bytes())

	requestData, err := config.Request.Marshall(config)
	if err != nil {
		fatal("could not marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", config.Endpoint, bytes.NewReader(requestData))
	if err != nil {
		fatal("could not create http request: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("bearer %v", config.Token))

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		fatal("query request failed: %v", err)
	}

	responseData := &bytes.Buffer{}
	_, err = io.Copy(responseData, res.Body)
	if err != nil {
		fatal("could not read response data: %v", err)
	}

	response, err := (&Response{}).Unmarshall(responseData.Bytes())
	if err != nil {
		fatal("could not parse response data: %v", err)
	}

	err = Render(config, response)
	if err != nil {
		fatal("could not render templates: %v", err)
	}

	fmt.Printf("Built from fresh data in %s\n", time.Since(start))
}
