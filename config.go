package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type Route struct {
	Host         string
	TargetURL    string
	ProxyHandler http.Handler
}

func parseConfigFile(filename string) ([]Route, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var routes []Route
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "->")
		if len(parts) != 2 {
			continue
		}

		hostName := strings.TrimSpace(parts[0])
		targetURLStr := strings.TrimSpace(parts[1])

		targetURL, err := url.Parse(targetURLStr)
		if err != nil {
			log.Printf("Warning: Invalid target URL %s: %v", targetURLStr, err)
			continue
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		routes = append(routes, Route{
			Host:         hostName,
			TargetURL:    targetURLStr,
			ProxyHandler: proxy,
		})

		fmt.Printf("Added route: %s -> %s\n", hostName, targetURLStr)
	}

	return routes, scanner.Err()
}
