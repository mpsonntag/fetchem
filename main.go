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
)

const ver = "fetchem 0.1.0"

func main() {
	usage := `fetchem

Usage:
  fetchem <url>
  fetchem <url> [--type="<filetype>"]
  fetchem -h | --help | --version

Options:
  -h  --help       Show this screen.
  --version        Show version.
`

	args, err := docopt.Parse(usage, nil, true, ver, false)
	if err != nil {
		fmt.Printf("An error has occurred trying to parse the cli options: %s\n", err.Error())
	}

	url := args["<url>"].(string)

	var fity string
	tmp := args["--type"]
	if tmp != nil {
		fity = tmp.(string)
		fmt.Printf("Content of --type: %s\n", tmp)
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching url: %s\n", err.Error())
		os.Exit(-1)
	}
	defer resp.Body.Close()

	if !strings.Contains(resp.Status, "200 OK") {
		fmt.Printf("Page was not found, return code: %s\n", resp.Status)
		os.Exit(-1)
	}

	fmt.Printf("Status: %s, Proto: %s, ContentLength: %d\n", resp.Status, resp.Proto, resp.ContentLength)

	s := bufio.NewScanner(bufio.NewReader(resp.Body))
	for s.Scan() {
		if fity == "" {
			fmt.Printf("%s\n", s.Text())
		} else if strings.Contains(s.Text(), fity) {
			fmt.Printf("%s\n", s.Text())
		}
	}

	os.Exit(0)
}
