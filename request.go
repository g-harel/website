package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
)

type ProjectRequest struct {
	Owner string
	Name  string
}

type ContributionRequest struct {
	Owner string
	Name  string
	Pull  int
	Issue int
}

type Request struct {
	Login         string
	Projects      []*ProjectRequest
	Contributions []*ContributionRequest
}

type GQLRequest struct {
	Query string `json:"query"`
}

func (req *Request) Marshall(c *Config) ([]byte, error) {
	tmpl, err := template.ParseFiles(c.Query)
	if err != nil {
		return nil, fmt.Errorf("could not parse query template: %v", err)
	}

	templateContent := bytes.Buffer{}
	err = tmpl.Execute(&templateContent, req)
	if err != nil {
		return nil, fmt.Errorf("could not execute query template: %v", err)
	}

	r := &GQLRequest{
		Query: templateContent.String(),
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %v", err)
	}

	return b, nil
}
