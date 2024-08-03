package functions

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/g-harel/website"
)

var env = struct {
	ConfigBucket    string
	ConfigObject    string
	GraphQLEndpoint string
	GraphQLToken    string
	TemplateBucket  string
	TemplateObject  string
	TemplateEntry   string
	UploadBucket    string
	UploadObject    string
}{
	ConfigBucket:    os.Getenv("CONFIG_BUCKET"),
	ConfigObject:    os.Getenv("CONFIG_OBJECT"),
	GraphQLEndpoint: os.Getenv("GRAPHQL_ENDPOINT"),
	GraphQLToken:    os.Getenv("GRAPHQL_TOKEN"),
	TemplateBucket:  os.Getenv("TEMPLATE_BUCKET"),
	TemplateObject:  os.Getenv("TEMPLATE_OBJECT"),
	TemplateEntry:   os.Getenv("TEMPLATE_ENTRY"),
	UploadBucket:    os.Getenv("UPLOAD_BUCKET"),
	UploadObject:    os.Getenv("UPLOAD_OBJECT"),
}

// Build is a background function triggered by messages to a pub/sub topic.
// It will build a new version of the website from a remote config, remote templates and updated api data.
// The result is uploaded to a public cloud storage bucket.
func Build(ctx context.Context, _ interface{}) error {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("create storage client: %v", err)
	}

	remoteConfig, err := storageClient.Bucket(env.ConfigBucket).Object(env.ConfigObject).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("create remote config reader: %v", err)
	}

	configBody := &bytes.Buffer{}
	_, err = io.Copy(configBody, remoteConfig)
	if err != nil {
		return fmt.Errorf("read remote config response: %v", err)
	}

	err = remoteConfig.Close()
	if err != nil {
		return fmt.Errorf("close remote config: %v", err)
	}

	config, err := (&website.Config{}).Parse(configBody.String())
	if err != nil {
		return fmt.Errorf("parse remote config: %v", err)
	}

	query, err := config.Query()
	if err != nil {
		return fmt.Errorf("generate query from config: %v", err)
	}

	dataReq, err := http.NewRequest("POST", env.GraphQLEndpoint, bytes.NewReader(query))
	if err != nil {
		return fmt.Errorf("create data request: %v", err)
	}
	dataReq = dataReq.WithContext(ctx)
	dataReq.Header.Add("Authorization", fmt.Sprintf("bearer %v", env.GraphQLToken))

	dataRes, err := httpClient.Do(dataReq)
	if err != nil {
		return fmt.Errorf("request for data failed")
	}

	dataBody := &bytes.Buffer{}
	_, err = io.Copy(dataBody, dataRes.Body)
	if err != nil {
		return fmt.Errorf("read data response: %v", err)
	}

	data, err := (&website.Data{}).Parse(dataBody.String())
	if err != nil {
		return fmt.Errorf("parse received data: %v", err)
	}

	templates, err := storageClient.Bucket(env.TemplateBucket).Object(env.TemplateObject).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("create remote template reader: %v", err)
	}

	templateDir, err := unzip(templates)
	if err != nil {
		return fmt.Errorf("unzip templates: %v", err)
	}

	err = templates.Close()
	if err != nil {
		return fmt.Errorf("close remote templates: %v", err)
	}

	for i := 0; i < len(config.Keyboards); i++ {
		data.Keyboards = append(data.Keyboards, &website.CreationData{
			Title:           config.Keyboards[i].Title,
			ImageURL:        config.Keyboards[i].ImageURL,
			BackgroundColor: config.Keyboards[i].BackgroundColor,
			Link:            config.Keyboards[i].Link,
		})
	}
	for i := 0; i < len(config.Illustrations); i++ {
		data.Illustrations = append(data.Illustrations, &website.CreationData{
			Title:           config.Illustrations[i].Title,
			ImageURL:        config.Illustrations[i].ImageURL,
			BackgroundColor: config.Illustrations[i].BackgroundColor,
			Link:            config.Illustrations[i].Link,
		})
	}
	for i := 0; i < len(config.Woodworking); i++ {
		data.Woodworking = append(data.Woodworking, &website.CreationData{
			Title:           config.Woodworking[i].Title,
			ImageURL:        config.Woodworking[i].ImageURL,
			BackgroundColor: config.Woodworking[i].BackgroundColor,
			Link:            config.Woodworking[i].Link,
		})
	}

	output, err := website.Render(templateDir, env.TemplateEntry, data)
	if err != nil {
		return fmt.Errorf("render templates: %v", err)
	}

	storageObject := storageClient.Bucket(env.UploadBucket).Object(env.UploadObject).NewWriter(ctx)

	_, err = io.Copy(storageObject, output)
	if err != nil {
		return fmt.Errorf("write output: %v", err)
	}

	err = storageObject.Close()
	if err != nil {
		return fmt.Errorf("write output to storage: %v", err)
	}

	return nil
}

// Unzip writes the archive data to a temporary directory and returns its path.
// The data is first written to a temporary file and decompressed from there.
func unzip(data io.ReadCloser) (string, error) {
	archiveFile, err := ioutil.TempFile("", "*.zip")
	if err != nil {
		return "", fmt.Errorf("create temp archive: %v", err)
	}

	_, err = io.Copy(archiveFile, data)
	if err != nil {
		return "", fmt.Errorf("write archive data to file: %v", err)
	}
	data.Close()

	outputDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", fmt.Errorf("create output directory: %v", err)
	}

	archive, err := zip.OpenReader(archiveFile.Name())
	if err != nil {
		return "", fmt.Errorf("create archive reader: %v", err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		if f.FileInfo().IsDir() {
			continue
		}

		outPath := filepath.Join(outputDir, f.Name)

		err := os.MkdirAll(path.Dir(outPath), os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("create file output path: %v", err)
		}

		contents, err := f.Open()
		if err != nil {
			return "", fmt.Errorf("open archive file: %v", err)
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			return "", fmt.Errorf("create output file: %v", err)
		}

		_, err = io.Copy(outFile, contents)
		if err != nil {
			return "", fmt.Errorf("write archive file: %v", err)
		}

		outFile.Close()
		contents.Close()
	}

	return outputDir, nil
}
