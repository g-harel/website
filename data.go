package website

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Data represents the combined response data from Config's query.
type Data struct {
	User          *UserData
	Projects      []*ProjectData
	Contributions []*ContributionData
	Keyboards     []*CreationData
}

// UserData represents the user data from Config's query.
type UserData struct {
	Icon        string `json:"avatarUrl"`
	Email       string `json:"email"`
	Description string `json:"bio"`
	Name        string `json:"name"`
	Login       string `json:"login"`
	Location    string `json:"location"`
	URL         string `json:"url"`
}

// ProjectData represents a project data item from Config's query.
type ProjectData struct {
	FullName string `json:"nameWithOwner"`
	Name     string `json:"name"`
	Owner    struct {
		Login string `json:"login"`
	} `json:"owner"`
	Description string `json:"description"`
	Homepage    string `json:"homepageUrl"`
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

// ContributionData represents a contribution data item from Config's query.
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

// CreationData represents the creation data from Config.
type CreationData struct {
	Title           string `json:"title"`
	ImageURL        string `json:"imageUrl"`
	BackgroundColor string `json:"backgroundColor"`
}

// GQLResponse represents a a generic GraphQL json response.
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

// Parse populates the fields of its receiver with unmarshalled contents from the raw json data.
func (d *Data) Parse(text string) (*Data, error) {
	res := &GQLResponse{}
	err := json.Unmarshal([]byte(text), res)
	if err != nil {
		return nil, fmt.Errorf("could not parse raw data: %v", err)
	}

	if res.Data == nil {
		return nil, fmt.Errorf("malformed raw data: %v", text)
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
		project.Homepage = strings.TrimPrefix(project.Homepage, "https://")
		project.Homepage = strings.TrimSuffix(project.Homepage, "/")
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
