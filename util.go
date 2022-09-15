package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/logrusorgru/aurora/v3"
	log "github.com/projectdiscovery/gologger"
)

func isStdin() bool {
	f, e := os.Stdin.Stat()
	if e != nil {
		return false
	}
	if f.Mode()&os.ModeNamedPipe == 0 {
		return false
	}

	return true
}

func isURL(s string) bool {
	_, e := url.ParseRequestURI(s)
	if e != nil {
		return false
	}

	u, e := url.Parse(s)
	if e != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func isError(e error) {
	if e != nil {
		log.Fatal().Msgf("%s\n", e)
	}
}

func showBanner() {
	fmt.Fprintf(os.Stderr, "%s\n", aurora.Cyan(banner))
}

func (opt *options) search() (bool, error) {
	queryEsc := url.QueryEscape(opt.Query)
	var regexes, baseURL, params string

	switch opt.Engine {
	case "google":
		regexes = `"><a href="\/url\?q=(.*?)&amp;sa=U&amp;`
		baseURL = "https://www.google.com/search"
		params = ("q=" + queryEsc + "&gws_rd=cr,ssl&client=ubuntu&ie=UTF-8&start=")
	case "shodan":
		regexes = `\"><a href=\"/host/(.*?)\">`
		baseURL = "https://www.shodan.io/search"
		params = ("query=" + queryEsc + "&page=")
	case "bing":
		regexes = `</li><li class=\"b_algo\"><h2><a href=\"(.*?)\" h=\"ID=SERP,`
		baseURL = "https://www.bing.com/search"
		params = ("q=" + queryEsc + "&first=")
	case "duck":
		regexes = `<a rel=\"nofollow\" href=\"//duckduckgo.com/l/\?kh=-1&amp;uddg=(.*?)\">`
		baseURL = "https://html.duckduckgo.com/html/"
		params = ("q=" + queryEsc + "&_=")
	case "yahoo":
		regexes = `\" ac-algo fz-l ac-21th lh-24\" href=\"(.*?)\" referrerpolicy=\"origin`
		baseURL = "https://search.yahoo.com/search"
		params = ("q=" + queryEsc + "&b=")
	case "ask":
		regexes = `target=\"_blank\" href='(.*?)' data-unified=`
		baseURL = "https://www.ask.com/web"
		params = ("q=" + queryEsc + "&page=")
	default:
		return true, errors.New("engine not found! Please choose one available")
	}

iterPage:
	for p := 1; p <= opt.Page; p++ {
		page := strconv.Itoa(p)
		switch opt.Engine {
		case "google":
			page += "0"
		case "yahoo", "bing":
			page += "1"
		}

		scrape := opt.get(baseURL + "?" + params + page)
		result := parser(scrape, regexes)
		for i := range result {
			url, err := url.QueryUnescape(result[i][1])
			if err != nil {
				return false, fmt.Errorf("when querying '%s' on page %d", queryEsc, p)
			}

			if !isURL(url) {
				break iterPage
			}

			fmt.Printf("%s\n", url)
		}

		if opt.Engine == "duck" && p == 1 {
			break
		}
	}

	return false, nil
}
