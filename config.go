package website

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"regexp"
	"strconv"
	"strings"
)

const (
	expectedSectionCount      = 6 // TODO(2024-07-29): use additional sections.
	loginSectionIndex         = 0
	projectsSectionIndex      = 1
	contributionsSectionIndex = 2
	keyboardsSectionIndex     = 3
)

// Config represents website configuration settings.
type Config struct {
	Login         string
	Projects      []*ProjectConfig
	Contributions []*ContributionConfig
	Keyboards     []*CreationConfig
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

// CreationConfig represents a creation config item.
type CreationConfig struct {
	Title           string
	ImageURL        string
	BackgroundColor string
	Link            string
}

// Parse populates the fields of its receiver with unmarshalled contents from the raw config.
func (c *Config) Parse(text string) (*Config, error) {
	// Remove comments.
	commentPattern := regexp.MustCompile("(^|\n)#[^\n]*")
	simplifiedText := commentPattern.ReplaceAllString(text, "")

	sections := strings.Split(strings.TrimSpace(simplifiedText), "\n\n")
	if len(sections) < expectedSectionCount {
		return nil, fmt.Errorf("missing config sections (%v/%v of login, projects, contributions, creations)", len(sections), expectedSectionCount)
	}
	if len(sections) > expectedSectionCount {
		return nil, fmt.Errorf("malformed config, too many sections (%v > %v)", len(sections), expectedSectionCount)
	}

	c.Login = sections[loginSectionIndex]
	c.Projects = []*ProjectConfig{}
	c.Contributions = []*ContributionConfig{}
	c.Keyboards = []*CreationConfig{}

	projects := []string{}
	if sections[projectsSectionIndex] != "" {
		projects = strings.Split(sections[projectsSectionIndex], "\n")
	}

	contributions := []string{}
	if sections[contributionsSectionIndex] != "" {
		contributions = strings.Split(sections[contributionsSectionIndex], "\n")
	}

	keyboards := []string{}
	if sections[keyboardsSectionIndex] != "" {
		keyboards = strings.Split(sections[keyboardsSectionIndex], "\n")
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
		parts := splitRepeated(contribution, " ")
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

	for _, keyboard := range keyboards {
		if !strings.Contains(keyboard, "|") {
			return nil, fmt.Errorf("creation missing separator: \"%v\"", keyboard)
		}
		sections := splitRepeated(keyboard, "|")
		if len(sections) != 2 {
			return nil, fmt.Errorf("malformed creation config: \"%v\"", keyboard)
		}
		metadata := splitRepeated(sections[1], " ")
		if len(metadata) != 3 {
			return nil, fmt.Errorf("missing creation metadata: \"%v\"", keyboard)
		}

		c.Keyboards = append(c.Keyboards, &CreationConfig{
			Title:           strings.TrimSpace(sections[0]),
			ImageURL:        metadata[1],
			BackgroundColor: metadata[0],
			Link:            metadata[2],
		})
	}

	return c, nil
}

func splitRepeated(s, sep string) []string {
	raw := strings.Split(s, sep)
	clean := []string{}
	for i := 0; i < len(raw); i++ {
		if raw[i] != "" {
			clean = append(clean, raw[i])
		}
	}
	return clean
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
