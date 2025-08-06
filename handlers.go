package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var rateLimiter = NewRateLimiter() // Global rate limiter instance

func createProxyHandler(routes []Route, notFound *template.Template, fileLogger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request:", r.Method, r.URL.Path)
		r.Body = http.MaxBytesReader(w, r.Body, 1048576) // Limit request body size to 1MB
		clientIP := getClientIP(r)
		if !rateLimiter.Allow(clientIP) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		LogRequest(r, fileLogger)
		host := r.Host

		for _, route := range routes {
			if route.Host == host {
				route.ProxyHandler.ServeHTTP(w, r)
				return
			}
		}

		// Fix: Set headers BEFORE WriteHeader
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		if notFound != nil {
			notFound.Execute(w, nil)
		}
	}
}
