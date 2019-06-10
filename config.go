package website

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"
)

// Config represents website configuration settings.
type Config struct {
	Login         string
	Projects      []*ProjectConfig
	Contributions []*ContributionConfig
}

// ProjectConfig represents a project config item.
type ProjectConfig struct {
	Owner string
	Name  string
}

// ContributionConfig represents a contribution config item.
type ContributionConfig struct {
	Owner string
	Name  string
	Pull  int
	Issue int
}

// Parse populates the fields of its receiver with unmarshalled contents from the raw config.
func (c *Config) Parse(text string) (*Config, error) {
	sections := strings.Split(strings.TrimSpace(text), "\n\n")
	if len(sections) < 3 {
		return nil, fmt.Errorf("missing config sections (1/%v of login, projects, config)", len(sections))
	}
	if len(sections) > 3 {
		return nil, fmt.Errorf("malformed config, too many sections (%v)", len(sections))
	}

	c.Login = string(sections[0])
	c.Projects = []*ProjectConfig{}
	c.Contributions = []*ContributionConfig{}

	projects := []string{}
	if sections[1] != "" {
		projects = strings.Split(sections[1], "\n")
	}

	contributions := []string{}
	if sections[2] != "" {
		contributions = strings.Split(sections[2], "\n")
	}

	for _, project := range projects {
		name := strings.Split(project, "/")
		if len(name) != 2 {
			return nil, fmt.Errorf("malformed project config: \"%v\"", project)
		}

		c.Projects = append(c.Projects, &ProjectConfig{
			Owner: string(name[0]),
			Name:  string(name[1]),
		})
	}

	for _, contribution := range contributions {
		parts := strings.Split(contribution, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("malformed contribution config: \"%v\"", contribution)
		}

		name := strings.Split(parts[0], "/")
		if len(name) != 2 {
			return nil, fmt.Errorf("malformed contribution config: \"%v\"", contribution)
		}

		pull, err := strconv.Atoi(string(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("could not parse contribution pull number in \"%v\": %v", contribution, err)
		}

		issue, err := strconv.Atoi(string(parts[2]))
		if err != nil {
			return nil, fmt.Errorf("could not parse contribution issue number in \"%v\": %v", contribution, err)
		}

		c.Contributions = append(c.Contributions, &ContributionConfig{
			Owner: string(name[0]),
			Name:  string(name[1]),
			Pull:  pull,
			Issue: issue,
		})
	}

	return c, nil
}

// Query generates a GraphQL query string from its receiver's fields.
func (c *Config) Query() ([]byte, error) {
	tmpl, err := template.New("test").Parse(query)
	if err != nil {
		return nil, fmt.Errorf("could not parse query template: %v", err)
	}

	templateContent := bytes.Buffer{}
	err = tmpl.Execute(&templateContent, c)
	if err != nil {
		return nil, fmt.Errorf("could not execute query template: %v", err)
	}

	r := &struct {
		Query string `json:"query"`
	}{
		Query: templateContent.String(),
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("could not marshal cuest: %v", err)
	}

	return b, nil
}

// GraphQL query template sent to GitHub api.
var query = `
	query {
		user(login: "{{.Login}}") {
			avatarUrl
			email
			bio
			name
			login
			location
			url
		}
		{{range $i, $p := .Projects}}
			p{{$i}}: repository(owner: "{{$p.Owner}}", name: "{{$p.Name}}") {
				...RepoInfo
			}
		{{end}}
		{{range $i, $c :=.Contributions}}
			c{{$i}}: repository(owner: "{{$c.Owner}}", name: "{{$c.Name}}") {
				name
				owner {
					login
				}
				url
				{{if not (eq $c.Pull 0)}}
					pullRequest(number: {{$c.Pull}}) {
						number
						title
						url
					}
				{{end}}
				{{if not (eq $c.Issue 0)}}
					issue(number: {{$c.Issue}}) {
						number
						url
					}
				{{end}}
			}
		{{end}}
	}

	fragment RepoInfo on Repository {
		nameWithOwner
		name
		owner {
			login
		}
		description
		url
		homepageUrl
		stargazers {
			totalCount
		}
		languages(first: 3, orderBy: {field: SIZE, direction: DESC}) {
			nodes {
				name
				color
			}
		}
	}
`
