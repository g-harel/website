package functions

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/g-harel/website"
)

var env = struct {
	ConfigSrc         string
	TemplateDir       string
	TemplateEntry     string
	GraphQLEndpoint   string
	GraphQLToken      string
	UploadBucket      string
	UploadName        string
	UploadCredentials string
}{
	ConfigSrc:         os.Getenv("CONFIG_SRC"),
	TemplateDir:       os.Getenv("TEMPLATE_DIR"),
	TemplateEntry:     os.Getenv("TEMPLATE_ENTRY"),
	GraphQLEndpoint:   os.Getenv("GRAPHQL_ENDPOINT"),
	GraphQLToken:      os.Getenv("GRAPHQL_TOKEN"),
	UploadBucket:      os.Getenv("UPLOAD_BUCKET"),
	UploadName:        os.Getenv("UPLOAD_NAME"),
	UploadCredentials: os.Getenv("UPLOAD_CREDENTIALS"),
}

// Returned errors are not yet logged (go cloud function alpha).
func fatal(format string, a ...interface{}) error {
	err := fmt.Errorf("fatal error: %v", fmt.Sprintf(format, a...))
	fmt.Fprint(os.Stderr, err.Error())
	return err
}

// Build is a background function triggered by messages to a pub/sub topic.
// It will build a new version of the website from a remote config and updated api data.
// The result is uploaded to a public cloud storage bucket.
func Build(ctx context.Context, _ interface{}) error {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	configReq, err := http.NewRequest("GET", env.ConfigSrc, nil)
	if err != nil {
		return fatal("could not create remote config request: %v", err)
	}
	configReq = configReq.WithContext(ctx)

	configRes, err := httpClient.Do(configReq)
	if err != nil {
		return fatal("request for remote config failed: %v", err)
	}

	configBody := &bytes.Buffer{}
	_, err = io.Copy(configBody, configRes.Body)
	if err != nil {
		return fatal("could not read remote config response: %v", err)
	}

	config, err := (&website.Config{}).Parse(configBody.String())
	if err != nil {
		return fatal("could not parse remote config: %v", err)
	}

	query, err := config.Query()
	if err != nil {
		return fatal("could not generate query from config: %s", err)
	}

	dataReq, err := http.NewRequest("POST", env.GraphQLEndpoint, bytes.NewReader(query))
	if err != nil {
		return fatal("could not create data request: %v", err)
	}
	dataReq = dataReq.WithContext(ctx)
	dataReq.Header.Add("Authorization", fmt.Sprintf("bearer %v", env.GraphQLToken))

	dataRes, err := httpClient.Do(dataReq)
	if err != nil {
		return fatal("request for data failed")
	}

	dataBody := &bytes.Buffer{}
	_, err = io.Copy(dataBody, dataRes.Body)
	if err != nil {
		return fatal("could not read data response: %v", err)
	}

	data, err := (&website.Data{}).Parse(dataBody.String())
	if err != nil {
		return fatal("could not parse received data: %v", err)
	}

	output, err := website.Render(env.TemplateDir, env.TemplateEntry, data)
	if err != nil {
		return fatal("could not render templates: %v", err)
	}

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return fatal("could not create storage client: %v", err)
	}

	storageObject := storageClient.Bucket(env.UploadBucket).Object(env.UploadName).NewWriter(ctx)

	_, err = io.Copy(storageObject, output)
	if err != nil {
		return fatal("could not write output to storage: %s", err)
	}

	err = storageObject.Close()
	if err != nil {
		return fatal("storage write failed: %s", err)
	}

	return nil
}
