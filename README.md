# go-dork

[![License](https://img.shields.io/badge/license-MIT-_red.svg)](https://opensource.org/licenses/MIT)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/dwisiswant0/go-dork/issues)

The fastest dork scanner written in Go.

<img src="https://user-images.githubusercontent.com/25837540/87547986-153f6a80-c6d6-11ea-92ef-bdc23d60f79e.png" height="350" width="566">

There are also various search engines supported by go-dork, including Google, Shodan, Bing, Duck, Yahoo and Ask.

- [Install](#install)
- [Usage](#usage)
  - [Basic Usage](#basic-usage)
  - [Flags](#flags)
  - [Querying](#querying)
  - [Defining engine](#defining-engine)
  - [Pagination](#pagination)
  - [Adding headers](#adding-headers)
  - [Using Proxy](#using-proxy)
  - [Chained with other tools](#chained-with-other-tools)
- [Help & Bugs](#help--bugs)
- [TODOs](#todos)
- [License](#license)
- [Version](#version)

## Install

- [Download](https://github.com/dwisiswant0/go-dork/releases) a prebuilt binary from releases page, unpack and run! or
- If you have go compiler installed and configured:

```bash
> GO111MODULE=on go get -v github.com/dwisiswant0/go-dork/...
```

## Usage

### Basic Usage

Simply, go-dork can be run with:

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
| -nc/--no-color | Disable colored output results                       |

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

By default, go-dork selects the first page, you can customize using the `-p` flag.

```bash
> go-dork -q "intext:'jira'" -p 5
```

It will search sequentially from pages 1 to 5.

### Adding Headers

Maybe you want to use a search filter on the Shodan engine, you can use custom headers to add cookies or other header parts.

```bash
> go-dork -q "org:'Target' http.favicon.hash:116323821" \
  --engine shodan -H "Cookie: ..." -H "User-Agent: ..."
```

### Using Proxy

Using a proxy, this can also be useful if Google or other engines meet Captcha.

```bash
> go-dork -q "intitle:'BigIP'" -p 2 -x http://127.0.0.1:8989
```

### Chained with other tools

If you want to chain the `go-dork` results with another tool, use the `-s` flag.

```bash
> go-dork -q "inurl:'/secure' intext:'jira' site:org" -s | nuclei -t workflows/jira-exploitaiton-workflow.yaml
```

## Help & Bugs

If you are still confused or found a bug, please [open the issue](https://github.com/dwisiswant0/go-dork/issues). All bug reports are appreciated, some features have not been tested yet due to lack of free time.

## TODOs

- [ ] Fixes Yahoo regexes
- [ ] Fixes Google regexes if using custom User-Agent
- [ ] Stopping if there's no results & page flag was set
- [ ] DuckDuckGo next page

## License

MIT. See `LICENSE` for more details.

## Version

**Current version is 0.0.1** and still development.
