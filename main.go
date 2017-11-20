// Copyright (c) 2016, Michael Sonntag (michael.p.sonntag@gmail.com)
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of the copyright holder nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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
	"github.com/atotto/clipboard"
)

const ver = "fetchem 0.1.1"

func main() {
	usage := `fetchem

Usage:
  fetchem (<url>) [-t <filetype>... | -r <fileRegExp>]
  fetchem -d <url>
  fetchem -e <url>
  fetchem -h | --help | --version

Options:
  <url>               url to fetch or from where to download, required option.
  -t <filetype>...    Specify which files to fetch, more than one option can be selected.
                          e.g. -t "png" -t "jpg"
                      If this option is not used, the code of the specified url will be printed
                      onto the screen.
  -r <fileRegExp>     For more fine grained specification of which files to fetch.
  -d <url>            Decode encoded URL, the decoded URL will be printed to the command line.
                      Alternative command line usage.
  -e <url>            Encode plain URL, the encoded URL will be printed to the command line.
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

	if args["-d"] != nil {
		err = decodeLink(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n%s\n\n", err.Error())
			os.Exit(-1)
		}
		os.Exit(0)
	}

	if args["-e"] != nil {
		encUrl := url.QueryEscape(args["-e"].(string))
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

// decodeLink splits multiple urlencoded urls from a string, decodes them
// and prints them to the commandline. If a supported clipboard handler is
// available, the second url is copied to the clipboard.
func decodeLink(args map[string]interface{}) error {
	const sep = "http"
	shinyUrl, err := url.QueryUnescape(args["-d"].(string))
	if err != nil {
		err = fmt.Errorf("[Error] decoding url: %s\n", err.Error())
		return err
	}

	shinySplit := strings.Split(shinyUrl, sep)
	for i, v := range shinySplit {
		if len(v) > 0 {
			fmt.Printf("\nDumbly unescaped string #%d: \n%s%s\n", i, sep, v)
		}

		// Copy the second occurrence to clipboard
		if i == 2 {
			err = clipboard.WriteAll(fmt.Sprintf("%s%s", sep, v))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
