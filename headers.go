package main

import "strings"

type customHeaders []string

func (h customHeaders) String() string {
	return strings.Join(h, ", ")
}

func (h *customHeaders) Set(val string) error {
	*h = append(*h, val)
	return nil
}
