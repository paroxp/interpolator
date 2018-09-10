package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/paroxp/interpolator/pkg/interpolator"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	dry          = kingpin.Flag("dry", "Dry run, will tell you which variables would be replaced.").Default("false").Bool()
	failover     = kingpin.Flag("failover", "Fail, if a variable is missing. Otherwise, an empty string would be placed.").Short('f').Default("false").Bool()
	matchPattern = kingpin.Flag("pattern", "Pattern to be used while matching the file for variables.").Short('m').Default("\\${([a-zA-Z0-9_]+?)}").OverrideDefaultFromEnvar("MATCH_PATTERN").String()
)

func main() {
	kingpin.Version("0.1.0")
	kingpin.Parse()

	buf := &bytes.Buffer{}
	n, err := io.Copy(buf, os.Stdin)
	if err != nil {
		log.Fatalln(err)
	} else if n <= 1 { // buffer always contains '\n'
		log.Fatalln("no input provided")
	}

	content := buf.Bytes()
	matches, err := interpolator.FindMatches(content, *matchPattern)
	if err != nil {
		log.Fatalln(err)
	}

	if *dry {
		for _, m := range matches {
			fmt.Printf("%s:\t\t'%s'\n", m.Name, m.Value)
		}
		os.Exit(0)
	}

	output, err := interpolator.ParseContent(content, matches, *failover)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(output))
}
