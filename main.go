// Copyright (c) 2016, Michael Sonntag (sonntag@bio.lmu.de)
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted under the terms of the BSD License. See
// LICENSE file in the root of the Project.

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/docopt/docopt-go"
)

const ver = "fetchem 0.1.0"

func main() {
	usage := `fetchem

Usage:
  fetchem <url>
  fetchem <url> [--out=<output>]
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

	output := args["--out"]
	if output != nil {
		fmt.Printf("Content of --out: %s\n", output)
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching url: %s\n", err.Error())
		os.Exit(-1)
	}

	fmt.Printf("Status: %s, Proto: %s, ContentLength: %d\n%b\n", resp.Status, resp.Proto, resp.ContentLength, resp.Body)
	os.Exit(0)
}
