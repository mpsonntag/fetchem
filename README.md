![run-tests](https://github.com/mpsonntag/fetchem/workflows/run-tests/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/mpsonntag/fetchem/badge.svg?branch=master)](https://coveralls.io/github/mpsonntag/fetchem?branch=master)

# Fetchem

Small command line tool to manipulate URLs on the command line.

```Usage:
  fetchem (<url>) [-t <filetype>... | -r <fileRegExp>]
  fetchem -d <url>
  fetchem -e <url>
  fetchem -h | --help | --version

Options:
  <url>               url to fetch or from where to download, required option.
  -t <filetype>...    Specify which files to fetch, more than one option can be selected.
                          e.g. -t "png" -t "md"
                      If this option is not used, the code of the specified url will be printed
                      onto the screen.
  -r <fileRegExp>     For more fine grained specification of which files to fetch.
  -d <url>            Decode encoded URL, the decoded URL will be printed to the command line.
                      Alternative command line usage.
  -e <url>            Encode plain URL, the encoded URL will be printed to the command line.
                      Alternative command line usage.
  -h  --help          Show this screen.
  --version           Show version.```
  