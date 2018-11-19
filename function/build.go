package function

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
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
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

// Returned errors are not yet logged (alpha).
func fatal(format string, a ...interface{}) error {
	err := fmt.Errorf("fatal error: %v", fmt.Sprintf(format, a...))
	fmt.Fprint(os.Stderr, err.Error())
	return err
}

func Build(ctx context.Context, _ interface{}) error {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	configReq, err := http.NewRequest("GET", env.ConfigSrc, nil)
	if err != nil {
		return fatal("could not create remote config request: %v", err)
	}

	configRes, err := httpClient.Do(configReq)
	if err != nil {
		return fatal("request for remote config failed")
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

	credentials, err := google.CredentialsFromJSON(ctx, []byte(env.UploadCredentials))
	if err != nil {
		return fmt.Errorf("could not create credentials from env: %v", err)
	}

	storageClient, err := storage.NewClient(ctx, option.WithCredentials(credentials))
	if err != nil {
		return fmt.Errorf("could not create storage client: %v", err)
	}

	storageObject := storageClient.Bucket(env.UploadBucket).Object(env.UploadName).NewWriter(ctx)
	for {
		line, readErr := output.ReadBytes('\n')
		if readErr != nil && readErr != io.EOF {
			return fmt.Errorf("could not read line from render output: %s", err)
		}

		_, err = storageObject.Write(bytes.TrimSpace(line))
		if err != nil {
			return fmt.Errorf("could not write to response: %s", err)
		}

		if readErr == io.EOF {
			break
		}
	}

	return nil
}
