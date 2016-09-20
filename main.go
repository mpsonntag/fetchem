// Copyright (c) 2016, Michael Sonntag (sonntag@bio.lmu.de)
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted under the terms of the BSD License. See
// LICENSE file in the root of the Project.

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"regexp"
)

const ver = "fetchem 0.1.0"

func main() {
	usage := `fetchem

Usage:
  fetchem (<url>) [-t <filetype>... | -r <fileRegExp>]
  fetchem -h | --help | --version

Options:
  <url>               url to fetch or from where to download, required option.
  -t <filetype>...    Specify which files to fetch, more than one option can be selected.
                          e.g. -t "png" -t "jpg"
                      If this option is not used, the code of the specified url will be printed
                      onto the screen.
  -r <fileRegExp>     For more fine grained specification of which files to fetch.
  -h  --help          Show this screen.
  --version           Show version.
`

	// examples:
	// go run main.go https://web.archive.org/web/20150320093805/http://lovecraftismissing.com/?p=8539 -t png -t jpg
	// go run main.go https://web.archive.org/web/20150320093805/http://lovecraftismissing.com/?p=8549 -r "(/web){1}(.)*/[0-9_+-]*.jpg"
	// go run main.go http://static.nichtlustig.de/toondb/150421.html -t png -t jpg
	// go run main.go http://static.nichtlustig.de/toondb/150421.html -r
	// go run main.go http://static.nichtlustig.de/toondb/150421.html -r "(//static){1}[0-9a-zA-Z._+-/:]*(/st/){1}[0-9a-zA-Z._+-/:]*.png"

	args, err := docopt.Parse(usage, nil, true, ver, false)
	if err != nil {
		fmt.Printf("An error has occurred trying to parse the cli options: %s\n", err.Error())
	}

	url := args["<url>"].(string)

	var fity []string
	if args["-t"] != nil {
		fity = args["-t"].([]string)
	}

	var fileReg string
	var expression *regexp.Regexp
	if args["-r"] != nil {
		fileReg = args["-r"].(string)
		expression = regexp.MustCompile(fileReg)
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching url: %s\n", err.Error())
		os.Exit(-1)
	}
	defer resp.Body.Close()

	if !strings.Contains(resp.Status, "200 OK") {
		fmt.Printf("Page for address %s was not found, return code: %s\n", url, resp.Status)
		os.Exit(-1)
	}

	fmt.Printf("Status: %s, Proto: %s, ContentLength: %d\n", resp.Status, resp.Proto, resp.ContentLength)

	var checkExists string
	s := bufio.NewScanner(bufio.NewReader(resp.Body))
	for s.Scan() {
		if len(fity) > 0 {
			for _, v := range fity {
				exp := `[a-zA-Z\d./\\_+-]*`+ v
				re := regexp.MustCompile(exp)
				checkExists = findRegexp(re, s.Text(), checkExists)
			}
		} else if fileReg != "" {
			checkExists = findRegexp(expression, s.Text(), checkExists)
		} else {
			fmt.Printf("%s\n", s.Text())
		}
	}

	os.Exit(0)
}

// findRegexp finds all occurrences that match a regular expression in a string
// and prints each match, if it does not already exist in a second check string.
// The check string is updated with new occurrences and returned.
func findRegexp(exp *regexp.Regexp, text string, checkExists string) string {
	match := exp.FindAllString(text, -1)
	for _, val := range match {
		if !strings.Contains(checkExists, val) {
			fmt.Println(val)
			checkExists = checkExists + val
		}
	}
	return checkExists
}
