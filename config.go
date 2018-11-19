package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func warn(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "\033[33;1mwarning: %s\033[0m\n", fmt.Sprintf(format, a...))
}

func splitStrict(b []byte, sep []byte) [][]byte {
	if len(b) == 0 {
		return [][]byte{}
	}
	return bytes.Split(b, sep)
}

func splitExact(b []byte, sep []byte, n int) [][]byte {
	split := bytes.SplitN(b, sep, n)
	for len(split) < n {
		split = append(split, []byte{})
	}
	return split
}

// Config configures how the output is generated.
type Config struct {
	Templates string
	RootName  string
	OutPath   string

	Query    string
	Endpoint string
	Token    string

	Request *Request
}

func (c *Config) Unmarshall(spec []byte) *Config {
	c.Templates = os.Getenv("TEMPLATE_DIR")
	c.RootName = os.Getenv("TEMPLATE_ENTRY")
	c.OutPath = os.Getenv("TEMPLATE_OUT")

	c.Query = os.Getenv("QUERY_TEMPLATE")
	c.Endpoint = os.Getenv("QUERY_ENDPOINT")
	c.Token = os.Getenv("QUERY_TOKEN")

	sections := splitExact(bytes.TrimSpace(spec), []byte{'\n', '\n'}, 3)

	c.Request = &Request{
		Login:         string(sections[0]),
		Projects:      []*ProjectRequest{},
		Contributions: []*ContributionRequest{},
	}

	projects := splitStrict(sections[1], []byte{'\n'})
	for _, project := range projects {
		name := splitExact(project, []byte{'/'}, 2)
		if len(name[0]) == 0 {
			warn("empty project owner: \"%s\"", project)
			break
		}
		if len(name[1]) == 0 {
			warn("empty project name: \"%s\"", project)
			break
		}

		c.Request.Projects = append(c.Request.Projects, &ProjectRequest{
			Owner: string(name[0]),
			Name:  string(name[1]),
		})
	}

	contributions := splitStrict(sections[2], []byte{'\n'})
	for _, contribution := range contributions {
		parts := splitExact(contribution, []byte{' '}, 3)

		name := splitExact(parts[0], []byte{'/'}, 2)
		if len(name[0]) == 0 {
			warn("empty contribution project owner: \"%s\"", contribution)
			break
		}
		if len(name[1]) == 0 {
			warn("empty contribution project name: \"%s\"", contribution)
			break
		}

		pull, err := strconv.Atoi(string(parts[1]))
		if err != nil {
			warn("could not parse contribution pull number: \"%s\"", contribution)
			pull = 0
		}

		issue, err := strconv.Atoi(string(parts[2]))
		if err != nil {
			warn("could not parse contribution issue number: \"%s\"", contribution)
			issue = 0
		}

		c.Request.Contributions = append(c.Request.Contributions, &ContributionRequest{
			Owner: string(name[0]),
			Name:  string(name[1]),
			Pull:  pull,
			Issue: issue,
		})
	}

	return c
}
