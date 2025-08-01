package main

import (
	"html/template"
	"log"
	"net/http"
)

func createProxyHandler(routes []Route, notFound *template.Template, fileLogger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		LogRequest(r, fileLogger)
		host := r.Host

		for _, route := range routes {
			if route.Host == host {
				route.ProxyHandler.ServeHTTP(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/html")
		if notFound != nil {
			notFound.Execute(w, nil)
		}
	}
}
