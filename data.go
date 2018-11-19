package website

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Data struct {
	User          *UserData
	Projects      []*ProjectData
	Contributions []*ContributionData
}

type UserData struct {
	Icon        string `json:"avatarUrl"`
	Email       string `json:"email"`
	Description string `json:"bio"`
	Name        string `json:"name"`
	Login       string `json:"login"`
	Location    string `json:"location"`
	URL         string `json:"url"`
}

type ProjectData struct {
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

type ContributionData struct {
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

func (d *Data) Parse(text string) (*Data, error) {
	res := &GQLResponse{}
	err := json.Unmarshal([]byte(text), res)
	if err != nil {
		return nil, fmt.Errorf("could not parse input: %v", err)
	}

	if res.Data == nil {
		return nil, fmt.Errorf("malformed input format, no data: %v", text)
	}
	if len(res.Errors) > 0 {
		return nil, fmt.Errorf("error (1/%v): %v", len(res.Errors), res.Errors[0].Message)
	}

	d.User = &UserData{}
	d.Projects = []*ProjectData{}
	d.Contributions = []*ContributionData{}

	err = json.Unmarshal(*res.Data["user"], d.User)
	if err != nil {
		return nil, fmt.Errorf("could not parse user data: %v", err)
	}

	index := 0
	for {
		data, ok := res.Data["p"+strconv.Itoa(index)]
		if !ok {
			break
		}

		project := &ProjectData{}
		err = json.Unmarshal(*data, project)
		if err != nil {
			return nil, fmt.Errorf("could not parse project data: %v", err)
		}
		d.Projects = append(d.Projects, project)

		index++
	}

	index = 0
	for {
		data, ok := res.Data["c"+strconv.Itoa(index)]
		if !ok {
			break
		}

		contribution := &ContributionData{}
		err = json.Unmarshal(*data, contribution)
		if err != nil {
			return nil, fmt.Errorf("could not parse contribution data: %v", err)
		}
		d.Contributions = append(d.Contributions, contribution)

		index++
	}

	return d, nil
}
