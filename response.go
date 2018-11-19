package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// TODO comments

type UserResponse struct {
	Icon        string `json:"avatarUrl"`
	Email       string `json:"email"`
	Description string `json:"bio"`
	Name        string `json:"name"`
	Login       string `json:"login"`
	Location    string `json:"location"`
	URL         string `json:"url"`
}

type ProjectResponse struct {
	FullName string `json:"nameWithOwner"`
	Name     string `json:"name"`
	Owner    struct {
		Login string `json:"login"`
	} `json:"owner"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Stargazers  struct {
		Count int `json:"totalCount"`
	} `json:"stargazers"`
	Languages struct {
		Nodes []struct {
			Name  string `json:"name"`
			Color string `json:"color"`
		} `json:"nodes"`
	} `json:"languages"`
}

type ContributionResponse struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	Pull struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		URL    string `json:"url"`
	} `json:"pullRequest"`
	Issue struct {
		Number int    `json:"number"`
		URL    string `json:"url"`
	} `json:"issue"`
}

type Response struct {
	User          *UserResponse
	Projects      []*ProjectResponse
	Contributions []*ContributionResponse
}

type GQLResponse struct {
	Data   map[string]*json.RawMessage `json:"data"`
	Errors []struct {
		Message   string `json:"message"`
		Locations []struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"location"`
	} `json:"errors"`
}

func (res *Response) Unmarshall(b []byte) (*Response, error) {
	queryResponse := &GQLResponse{}
	err := json.Unmarshal(b, queryResponse)
	if err != nil {
		return nil, fmt.Errorf("could not parse response data: %v", err)
	}

	if queryResponse.Data == nil {
		return nil, fmt.Errorf("missing response data: %s", b)
	}
	if len(queryResponse.Errors) > 0 {
		return nil, fmt.Errorf("query error (1/%v): %v", len(queryResponse.Errors), queryResponse.Errors[0].Message)
	}

	res.User = &UserResponse{}
	res.Projects = []*ProjectResponse{}
	res.Contributions = []*ContributionResponse{}

	err = json.Unmarshal(*queryResponse.Data["user"], res.User)
	if err != nil {
		return nil, fmt.Errorf("could not parse user data: %v", err)
	}

	index := 0
	for {
		data, ok := queryResponse.Data["p"+strconv.Itoa(index)]
		if !ok {
			break
		}

		project := &ProjectResponse{}
		err = json.Unmarshal(*data, project)
		if err != nil {
			return nil, fmt.Errorf("could not parse project data: %v", err)
		}
		res.Projects = append(res.Projects, project)

		index++
	}

	index = 0
	for {
		data, ok := queryResponse.Data["c"+strconv.Itoa(index)]
		if !ok {
			break
		}

		contribution := &ContributionResponse{}
		err = json.Unmarshal(*data, contribution)
		if err != nil {
			return nil, fmt.Errorf("could not parse contribution data: %v", err)
		}
		res.Contributions = append(res.Contributions, contribution)

		index++
	}

	return res, nil
}
