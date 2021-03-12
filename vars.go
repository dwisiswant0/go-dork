package main

import (
	"net/http"
	"sync"
)

var (
	query, engine, proxy string
	headers              customHeaders
	silent               bool
	page                 int

	queries []string
	client  http.Client
	wg      sync.WaitGroup
)
