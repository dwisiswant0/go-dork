package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	log "github.com/projectdiscovery/gologger"
)

func Err(e error) {
	if e != nil {
		log.Errorf("%s\n", e)
		os.Exit(1)
	}
}

type Options struct {
	Query, Engine, Proxy string
	Page                 int
	Headers              []string
}

type Headers []string

const banner = `
                   __         __  
  ___ ____  ______/ /__  ____/ /__
 / _ '/ _ \/__/ _  / _ \/ __/  '_/
 \_, /\___/   \_,_/\___/_/ /_/\_\ 
/___/
       v0.0.1 - @dwisiswant0
`

var query, engine, proxy string
var headers Headers
var page int
var noColor, silent bool

func (h Headers) String() string {
	return strings.Join(h, ", ")
}

func (h *Headers) Set(val string) error {
	*h = append(*h, val)
	return nil
}

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

	flag.BoolVar(&noColor, "nc", false, "")
	flag.BoolVar(&noColor, "no-color", false, "")

	flag.BoolVar(&silent, "s", false, "")
	flag.BoolVar(&silent, "silent", false, "")

	flag.Usage = func() {
		c := color.New(color.FgCyan, color.Bold).FprintfFunc()
		h := []string{
			banner,
			"Options:",
			"  -q, --query <query>          Search query",
			"  -e, --engine <engine>        Provide search engine (default: Google)",
			"                               (options: Google, Shodan, Bing, Duck, Yahoo, Ask)",
			"  -p, --page <i>               Specify number of pages (default: 1)",
			"  -H, --header <header>        Pass custom header to search engine",
			"  -x, --proxy <proxy_url>      Use proxy to surfing",
			"  -s, --silent                 Silent mode",
			"  -nc, --no-color              Disable colored output results",
			"",
		}

		c(os.Stderr, strings.Join(h, "\n"))
	}
}

func main() {
	flag.Parse()

	engine = strings.ToLower(engine)

	if !silent {
		c := color.New(color.FgCyan, color.Bold)
		c.Println(banner)
		log.Labelf("Use at your own risk! Developers assume no responsibility")
		log.Labelf("If your IP address has been blocked by search engine providers or other reason.")
		log.Infof("Query : %s", query)
		log.Infof("Page  : %s", strconv.Itoa(page))
		if proxy != "" {
			log.Infof("Proxy : %s", proxy)
		}
		if len(headers) > 0 {
			for _, h := range headers {
				log.Infof("Header: %s", h)
			}
		}
		log.Infof("Engine: %s\n\n", strings.Title(engine))
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			query = sc.Text()
			opts := Options{
				Query:   query,
				Engine:  engine,
				Page:    page,
				Proxy:   proxy,
				Headers: headers,
			}

			err := opts.Search(noColor)
			Err(err)
		}
	} else {
		if query == "" {
			log.Fatalf("Missing required -q flag!")
			os.Exit(2)
		}

		opts := Options{
			Query:   query,
			Engine:  engine,
			Page:    page,
			Proxy:   proxy,
			Headers: headers,
		}

		err := opts.Search(noColor)
		Err(err)
	}
}

func Parser(html string, pattern string) [][]string {
	regex := regexp.MustCompile(pattern)
	match := regex.FindAllStringSubmatch(html, -1)[0:]
	return match
}

func (opt *Options) Get(url string) string {
	client := Client(opt.Proxy)
	req, err := http.NewRequest("GET", url, nil)
	for _, h := range opt.Headers {
		parts := strings.SplitN(h, ":", 2)

		if len(parts) != 2 {
			continue
		}
		req.Header.Set(parts[0], parts[1])
	}
	Err(err)

	resp, err := client.Do(req)
	Err(err)
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	body := string(data)
	return body
}

func (opt *Options) Search(noColor bool) error {
	queryEsc := url.QueryEscape(opt.Query)
	var regexes, baseURL, params string
	var o = color.New(color.FgGreen)

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
		return errors.New("Engine not found! Please choose one available.")
	}

	for p := 1; p <= opt.Page; p++ {
		page := strconv.Itoa(p)

		switch opt.Engine {
		case "google":
			page += "0"
		case "yahoo", "bing":
			page += "1"
		}

		html := opt.Get(baseURL + "?" + params + page)
		result := Parser(html, regexes)
		for i := range result {
			url, err := url.QueryUnescape(result[i][1])
			Err(err)
			if noColor {
				fmt.Println(url)
			} else {
				o.Println(url)
			}
		}

		if opt.Engine == "duck" && p == 1 {
			break
		}
	}

	return nil
}

func Client(proxy string) *http.Client {
	tr := &http.Transport{
		MaxIdleConns:    30,
		IdleConnTimeout: time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   time.Second * 10,
			KeepAlive: time.Second,
		}).DialContext,
	}

	if proxy != "" {
		if p, err := url.Parse(proxy); err == nil {
			tr.Proxy = http.ProxyURL(p)
		}
	}

	re := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &http.Client{
		Transport:     tr,
		CheckRedirect: re,
		Timeout:       time.Second * 10,
	}
}
