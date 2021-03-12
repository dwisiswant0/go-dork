package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/logrusorgru/aurora/v3"
	log "github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

func init() {
	flag.StringVar(&query, "q", "", "")
	flag.StringVar(&query, "query", "", "")

	flag.StringVar(&engine, "e", "", "")
	flag.StringVar(&engine, "engine", "google", "")

	flag.IntVar(&page, "p", 1, "")
	flag.IntVar(&page, "page", 1, "")

	flag.Var(&headers, "header", "")
	flag.Var(&headers, "H", "")

	flag.StringVar(&proxy, "x", "", "")
	flag.StringVar(&proxy, "proxy", "", "")

	flag.BoolVar(&silent, "s", false, "")
	flag.BoolVar(&silent, "silent", false, "")

	flag.Usage = func() {
		h := []string{
			"Options:",
			"  -q, --query <query>          Search query",
			"  -e, --engine <engine>        Provide search engine (default: Google)",
			"                               (options: Google, Shodan, Bing, Duck, Yahoo, Ask)",
			"  -p, --page <i>               Specify number of pages (default: 1)",
			"  -H, --header <header>        Pass custom header to search engine",
			"  -x, --proxy <proxy_url>      Use proxy to surfing (HTTP/SOCKSv5 proxy)",
			"  -s, --silent                 Silent mode",
			"\n",
		}
		showBanner()
		fmt.Fprintf(os.Stderr, "%s", aurora.Green(strings.Join(h, "\n")))
	}
	flag.Parse()

	engine = strings.ToLower(engine)

	maxLog := levels.LevelDebug
	if silent {
		maxLog = levels.LevelSilent
	}
	log.DefaultLogger.SetMaxLevel(maxLog)

	showBanner()
}

func main() {
	if isStdin() {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			q := sc.Text()
			queries = append(queries, q)
		}
	} else {
		if query == "" {
			log.Fatal().Msgf("Missing required -q flag!")
			os.Exit(2)
		}

		queries = []string{query}
	}

	log.Info().Msgf("Query : %+v", queries)
	log.Info().Msgf("Page  : %s", strconv.Itoa(page))
	if proxy != "" {
		log.Info().Msgf("Proxy : %s", proxy)
	}
	if len(headers) > 0 {
		log.Info().Msgf("Header: [%+v]", headers)
	}
	log.Info().Msgf("Engine: %s", strings.Title(engine))
	log.Warning().Msg("Use at your own risk! Developers assume no responsibility...")
	log.Warning().Msg("If your IP address has been blocked by search engine providers or other reason.\n\n")

	for _, q := range queries {
		wg.Add(1)
		go func(dork string) {
			defer wg.Done()
			opts := options{
				Query:   dork,
				Engine:  engine,
				Page:    page,
				Proxy:   proxy,
				Headers: headers,
			}

			fatal, err := opts.search()
			if err != nil {
				if fatal {
					isError(err)
				}
				log.Warning().Msgf("Something error %s.", err.Error())
			}
		}(q)
	}
	wg.Wait()
}
