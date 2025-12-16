package main

import (
	"io"
	"net/http"
	"strings"

	"ktbs.dev/mubeng/pkg/mubeng"
)

// maxResponseSize limits response body to 10MB to prevent resource exhaustion
const maxResponseSize = 10 * 1024 * 1024

func (opt *options) get(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	// Check error immediately to prevent nil pointer dereference
	if err != nil {
		isError(err)
		return ""
	}

	for _, h := range opt.Headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			continue
		}
		req.Header.Set(parts[0], parts[1])
	}

	// Use local client to avoid race condition on shared client.Transport
	httpClient := &client
	if opt.Proxy != "" {
		transport, err := mubeng.Transport(opt.Proxy)
		if err != nil {
			isError(err)
			return ""
		}
		httpClient = &http.Client{Transport: transport}
		req.Header.Add("Connection", "close")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		isError(err)
		return ""
	}
	defer resp.Body.Close()

	// Limit response size to prevent resource exhaustion
	limitedReader := io.LimitReader(resp.Body, maxResponseSize)
	data, _ := io.ReadAll(limitedReader)
	return string(data)
}
