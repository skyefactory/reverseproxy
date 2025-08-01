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
	"time"
	"html/template"
)

const port = ":80"

type Route struct {
	Host         string
	TargetURL    string
	ProxyHandler http.Handler
}

func main() {
	// Setup logging to file
	logFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	// Create a multi-writer to log to both console and file
	fileLogger := log.New(logFile, "", log.LstdFlags)

	// Parse routes from config file
	routes, err := parseConfigFile("config.cfg")
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	notFound, err := template.ParseFiles("404.html")
    if err != nil {
        log.Printf("Warning: Could not load 404 template: %v", err)
    }

	fmt.Println("Reverse proxy running on http://localhost" + port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		LogRequest(r, fileLogger)
		host := r.Host

		// Check if we have a route for this host
		for _, route := range routes {
			if route.Host == host {
				// We found a matching route, proxy the request
				route.ProxyHandler.ServeHTTP(w, r)
				return
			}
		}

		// No matching route found, return 404
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/html")
		if notFound != nil {
			notFound.Execute(w, nil)
		}
	})

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
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
		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "->")
		if len(parts) != 2 {
			continue // Invalid line format
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

func LogRequest(r *http.Request, fileLogger *log.Logger) {
	host := r.Host
	url := r.URL.String()
	clientIP := r.RemoteAddr
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		clientIP = forwardedFor
	}

	method := r.Method
	userAgent := r.UserAgent()
	protocol := r.Proto
	log.Printf("Request received: Host=%s, IP=%s, URL=%s", host, clientIP, url)

	// Log detailed information to file
	logEntry := fmt.Sprintf(
		"Time: %s, Host: %s, URL Path: %s, Client IP: %s, Method: %s, Protocol: %s, User-Agent: %s",
		time.Now().Format(time.RFC3339),
		host,
		url,
		clientIP,
		method,
		protocol,
		userAgent,
	)

	fileLogger.Println(logEntry)
}
