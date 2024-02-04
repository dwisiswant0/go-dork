# go-dork

[![License](https://img.shields.io/badge/license-MIT-_red.svg)](https://opensource.org/licenses/MIT)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/dwisiswant0/go-dork/issues)

The fastest dork scanner written in Go.

<img src="https://user-images.githubusercontent.com/25837540/111008561-f22f9c80-83c3-11eb-8500-fb63456a4614.png" height="350">

There are also various search engines supported by go-dork, including Google, Shodan, Bing, Duck, Yahoo and Ask.

- [Install](#install)
- [Usage](#usage)
  - [Basic Usage](#basic-usage)
  - [Flags](#flags)
  - [Querying](#querying)
  - [Defining engine](#defining-engine)
  - [Pagination](#pagination)
  - [Adding custom headers](#adding-headers)
  - [Using proxy](#using-proxy)
  - [Chained with other tools](#chained-with-other-tools)
- [Supporting Materials](#supporting-materials)
- [Help & Bugs](#help--bugs)
- [TODOs](#todos)
- [License](#license)
- [Version](#version)

## Install

- [Download](https://github.com/dwisiswant0/go-dork/releases) a prebuilt binary from releases page, unpack and run! or
- If you have [Go 1.15+](https://golang.org/dl/) compiler installed and configured:

```bash
> GO111MODULE=on go install github.com/dwisiswant0/go-dork@latest
```

## Usage

### Basic Usage

It's fairly simple, go-dork can be run with:

```bash
> go-dork -q "inurl:'...'"
```

### Flags

```bash
> go-dork -h
```

This will display help for the tool. Here are all the switches it supports.

| Flag           | Description                                          |
|----------------|------------------------------------------------------|
| -q/--query     | Search query _(required)_                            |
| -e/--engine    | Provide search engine (default: Google)              |
|                | _(options: Google, Shodan, Bing, Duck, Yahoo, Ask)_  |
| -p/--page      | Specify number of pages (default: 1)                 |
| -H/--header    | Pass custom header to search engine                  |
| -x/--proxy     | Use proxy to surfing                                 |
| -s/--silent    | Silent mode, prints only results in output           |

### Querying

```bash
> go-dork -q "inurl:..."
```

Queries can also be input with stdin

```bash
> cat dorks.txt | go-dork -p 5
```

### Defining engine

Search engine can be changed from the available engines: Google, Shodan, Bing, Duck, Yahoo, Ask.
However, if the `-e` flag is not defined, it will use the Google search engine by default.

```bash
> go-dork -e bing -q ".php?id="
```

This will do a search by the Bing engine.

### Pagination

By default, go-dork scrapes the first page, you can customize using the `-p` flag.

```bash
> go-dork -q "intext:'jira'" -p 5
```

It will search sequentially from pages 1 to 5.

### Adding custom headers

Maybe you want to use a search filter on the Shodan engine, you can use custom headers to add cookies or other header parts.

```bash
> go-dork -q "org:'Target' http.favicon.hash:116323821" \
  --engine shodan -H "Cookie: ..." -H "User-Agent: ..."
```

### Using proxy

Using a proxy, this can also be useful if Google or other engines meet Captcha.

```bash
> go-dork -q "intitle:'BigIP'" -p 2 -x http://127.0.0.1:8989
```

### Chained with other tools

If you want to chain the `go-dork` results with another tool, use the `-s` flag.

```bash
> cat dorks.txt | go-dork | pwntools
> go-dork -q "inurl:'/secure' intext:'jira' site:org" -s | nuclei -t workflows/jira-exploitaiton-workflow.yaml
```

## Supporting Materials

- Hazana. _[Dorking on Steroids](https://hazanasec.github.io/2021-03-11-Dorking-on-Steriods/)_, 11 Mar. 2021, https://hazanasec.github.io/2021-03-11-Dorking-on-Steriods/.

## Help & Bugs

If you are still confused or found a bug, please [open the issue](https://github.com/dwisiswant0/go-dork/issues). All bug reports are appreciated, some features have not been tested yet due to lack of free time.

## TODOs

- [ ] Fixes Yahoo regexes
- [ ] Fixes Google regexes if using custom User-Agent
- [x] Stopping if there's no results & page flag was set
- [ ] DuckDuckGo next page

## License

MIT. See `LICENSE` for more details.
