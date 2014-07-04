package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"
)

var Log *log.Logger

func init() {
	Log = log.New(os.Stdout, "[kati] ", log.Lmicroseconds)
}

type PathMapping struct {
	Regex    *regexp.Regexp
	HostPort string
}

func NewRegexRP(pathMappings []PathMapping) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		lastIndex := len(pathMappings) - 1
		for index, mapping := range pathMappings {
			if mapping.Regex.MatchString(req.URL.Path) || index == lastIndex {
				req.URL.Scheme = "http" //just http for now
				req.Host = mapping.HostPort
				req.URL.Host = mapping.HostPort
				Log.Printf("Proxying: %s -> %s", req.URL.Path, mapping.HostPort)
				return
			}
		}
	}

	return &httputil.ReverseProxy{Director: director}
}

func parseFlagMappings(flagMappings []string) ([]PathMapping, error) {
	var err error
	pathMappings := make([]PathMapping, len(flagMappings))
	for index, mapping := range flagMappings {
		parts := strings.Split(mapping, "->")
		if len(parts) > 2 { //case where an arrow is part of the regex
			lastIndex := len(parts) - 1
			parts = []string{strings.Join(parts[0:lastIndex], "->"), parts[lastIndex]}
		}

		if len(parts) != 2 {
			return nil, errors.New(fmt.Sprintf("Invalid mapping syntax for: %v", mapping))
		}
		regex := strings.TrimSpace(parts[0])
		hostPort := strings.TrimSpace(parts[1])
		compiledRegex, err := regexp.Compile(regex)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid regex: %v", regex))
		}
		pathMappings[index] = PathMapping{Regex: compiledRegex, HostPort: hostPort}
	}

	return pathMappings, err
}

func runProxyServer(c *cli.Context) {
	mappings, parseErr := parseFlagMappings(c.StringSlice("proxy"))
	if parseErr != nil {
		fmt.Printf("Error parsing arguments: %v", parseErr)
	} else {
		listenErr := http.ListenAndServe(fmt.Sprintf(":%d", c.Int("http-port")), NewRegexRP(mappings))
		if listenErr != nil {
			fmt.Printf("Error starting server: %v", listenErr)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "kati"
	app.Version = "1.0.0"
	app.Usage = "Simple proxy server to send requests to different hosts based on path matched by a regex."
	app.Action = runProxyServer
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{"proxy", &cli.StringSlice{}, "Proxy mappings look like this:  \"/api/.* -> api.example.com\""},
		cli.IntFlag{"http-port", 80, "Port to listen for HTTP (not TLS)"},
	}
	app.Run(os.Args)
}
