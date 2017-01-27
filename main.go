// Copyright (c) 2016, Michael Sonntag (michael.p.sonntag@gmail.com)
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
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/docopt/docopt-go"
)

const ver = "fetchem 0.1.1"

func main() {
	usage := `fetchem

Usage:
  fetchem (<url>) [-t <filetype>... | -r <fileRegExp>]
  fetchem --decode <url>
  fetchem --encode <url>
  fetchem -h | --help | --version

Options:
  <url>               url to fetch or from where to download, required option.
  -t <filetype>...    Specify which files to fetch, more than one option can be selected.
                          e.g. -t "png" -t "jpg"
                      If this option is not used, the code of the specified url will be printed
                      onto the screen.
  -r <fileRegExp>     For more fine grained specification of which files to fetch.
  --decode <url>      encoded URL, the decoded URL will be printed to the command line.
                      Alternative command line usage.
  --encode <url>      plain URL, the encoded URL will be printed to the command line.
                      Alternative command line usage.
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
		fmt.Fprintf(os.Stderr, "[Error] parsing command line options: %s\n", err.Error())
		os.Exit(-1)
	}

	if decode, ok := args["--decode"]; ok && decode != nil {
		shinyUrl, err := url.QueryUnescape(decode.(string))
		if err != nil {
			fmt.Fprintf(os.Stderr, "[Error] decoding url: %s\n", err.Error())
			os.Exit(-1)
		}
		fmt.Printf("Dumbly unescaped string: \n\n%s\n\n", shinyUrl)
		os.Exit(0)
	}

	if enc, ok := args["--encode"]; ok && enc != nil {
		encUrl := url.QueryEscape(enc.(string))
		fmt.Printf("Dumbly escaped string: \n\n%s\n\n", encUrl)
		os.Exit(0)
	}

	fetchThis := args["<url>"].(string)

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

	resp, err := http.Get(fetchThis)
	if err != nil {
		fmt.Printf("Error fetching url: %s\n", err.Error())
		os.Exit(-1)
	}
	defer resp.Body.Close()

	if !strings.Contains(resp.Status, "200 OK") {
		fmt.Printf("Page for address %s was not found, return code: %s\n", fetchThis, resp.Status)
		os.Exit(-1)
	}

	fmt.Printf("Status: %s, Proto: %s, ContentLength: %d\n", resp.Status, resp.Proto, resp.ContentLength)

	var checkExists string
	s := bufio.NewScanner(bufio.NewReader(resp.Body))
	for s.Scan() {
		if len(fity) > 0 {
			for _, v := range fity {
				exp := `[a-zA-Z\d./\\_+-]*` + v
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
