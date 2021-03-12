package main

import "regexp"

func parser(html string, pattern string) [][]string {
	regex := regexp.MustCompile(pattern)
	match := regex.FindAllStringSubmatch(html, -1)[0:]
	return match
}
