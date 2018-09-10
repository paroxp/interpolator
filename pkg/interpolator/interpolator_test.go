package interpolator_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/paroxp/interpolator/pkg/interpolator"
)

func TestNewTemplate(t *testing.T) {
	os.Setenv("TEST", "12345")
	os.Setenv("TESTER", "gopher")

	for _, c := range []struct {
		in       string
		inName   string
		outValue string
	}{
		{"${TEST}", "TEST", "12345"},
		{"${TESTER}", "TESTER", "gopher"},
		{"${UNBOUND_VARIABLE}", "UNBOUND_VARIABLE", ""},
	} {
		tmpl := interpolator.NewTemplate(c.inName, c.in)

		if tmpl.Value != c.outValue {
			t.Errorf("expected '%s', got '%s'", c.outValue, tmpl.Value)
		}
	}

	os.Unsetenv("TEST")
	os.Unsetenv("TESTER")
}

func TestFindMatches(t *testing.T) {
	for _, c := range []struct {
		inContext []byte
		inRegex   string
		outLen    int
		outErr    error
		outArr    []string
	}{
		{[]byte("Hello ${GREET}. How you doing ${GREET}?"), "\\${([a-zA-Z0-9_]+?)}", 1, nil, []string{"GREET"}},
		{[]byte("Hello ${GREET}. How you doing ${NAME}? Cool ${GREETz}"), "\\${([A-Z0-9_]+?)}", 2, nil, []string{"GREET", "NAME"}},
		{[]byte("Hello ${GREET}. How you doing ${GREET}?"), "~((?<=^[^\\s])|(?<=\\s[^\\s]))\\s(?=[^\\s](\\s|$))~", 0, fmt.Errorf("invalid or unsupported Perl syntax"), nil},
		{[]byte("Hello ((GREET)). How you doing ((GREET))?"), "\\(\\(([A-Z0-9]+)\\)\\)", 1, nil, []string{"GREET"}},
	} {
		matches, err := interpolator.FindMatches(c.inContext, c.inRegex)

		if c.outErr != nil {
			if !strings.Contains(err.Error(), c.outErr.Error()) {
				t.Errorf("expected error to be like '%s', got '%s'", c.outErr, err)
			}
		}

		if c.outErr == err && err != nil {
			t.Errorf("expected error to be %s, got %s", c.outErr, err)
		}

		if len(matches) != c.outLen {
			t.Errorf("expected the lenght of %v, to be %d, got %d", matches, c.outLen, len(matches))
		}

		for i := range matches {
			if c.outArr[i] != matches[i].Name {
				t.Errorf("expected %s to equal to %s", c.outArr[i], matches[i].Name)
			}
		}
	}
}

func TestParseContent(t *testing.T) {
	os.Setenv("GREET", "Tester")

	for _, c := range []struct {
		inContent  []byte
		inFailover bool
		out        []byte
		outErr     error
	}{
		{[]byte("Hello ${GREET}. How you doing ${GREET}?"), false, []byte("Hello Tester. How you doing Tester?"), nil},
		{[]byte("Hello ${GREET}. How you doing ${GRE3T}?"), false, []byte("Hello Tester. How you doing ?"), nil},
		{[]byte("Hello ${GREET}. How you doing ${GRE3T}?"), true, nil, fmt.Errorf("expected GRE3T environment variable to be set")},
	} {
		matches, _ := interpolator.FindMatches(c.inContent, "\\${([a-zA-Z0-9_]+?)}")
		output, err := interpolator.ParseContent(c.inContent, matches, c.inFailover)

		if c.outErr != nil {
			if !strings.Contains(err.Error(), c.outErr.Error()) {
				t.Errorf("expected error to be like '%s', got '%s'", c.outErr, err)
			}
		}

		if c.outErr == err && err != nil {
			t.Errorf("expected error to be %s, got %s", c.outErr, err)
		}

		if string(output) != string(c.out) {
			t.Errorf("expected '%s', got '%s'", c.out, output)
		}
	}

	os.Unsetenv("GREET")
}
