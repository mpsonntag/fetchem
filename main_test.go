package main

import (
	"regexp"
	"strings"
	"testing"
)

func TestFindRegexp(t *testing.T) {
	exp := `[a-zA-Z\d./\\_+-]*jpg`
	re := regexp.MustCompile(exp)
	res := ""
	restring := "abc.jpg"
	res = findRegexp(re, restring, res)

	if !strings.Contains(res, restring) {
		t.Fatalf("Failed to match string: '%s'", restring)
	}

	restring = "abc.png"
	res = findRegexp(re, restring, res)
	if strings.Contains(res, restring) {
		t.Fatalf("Invalidly matched string: '%s'", restring)
	}
}
