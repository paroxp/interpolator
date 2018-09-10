package interpolator

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
)

// Template would hold the replacement information for the variable to be found in text content.
type Template struct {
	Match string
	Name  string
	Value string
}

// NewTemplate is to return a ready to use struct that is capable of replacing a variable in a given string.
func NewTemplate(name, v string) *Template {
	return &Template{
		Match: v,
		Name:  name,
		Value: os.Getenv(name),
	}
}

// FindMatches will attempt to run regex expression over the content finding all the matches.
func FindMatches(content []byte, regex string) ([]*Template, error) {
	matches := []*Template{}
	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	for name, variable := range uniqueVariables(r.FindAllStringSubmatch(string(content), -1)) {
		matches = append(matches, NewTemplate(name, variable))
	}

	return matches, nil
}

// ParseContent will itterate through the matched variables and replace on the go.
// If failOnMissing is set to true, it will return an error
func ParseContent(content []byte, matches []*Template, failOnMissing bool) ([]byte, error) {
	for _, t := range matches {
		if failOnMissing && t.Value == "" {
			return nil, fmt.Errorf("expected %s environment variable to be set", t.Name)
		}

		content = bytes.Replace(content, []byte(t.Match), []byte(t.Value), -1)
	}

	return content, nil
}

func uniqueVariables(s [][]string) map[string]string {
	used := map[string]string{}
	for _, match := range s {
		used[match[1]] = match[0]
	}

	return used
}
