package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

const port = ":80"

func main() {
	logFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	fileLogger := log.New(logFile, "", log.LstdFlags)

	routes, err := parseConfigFile("config.cfg")
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	notFound, err := template.ParseFiles("404.html")
	if err != nil {
		log.Printf("Warning: Could not load 404 template: %v", err)
	}

	fmt.Println("Reverse proxy running on http://localhost" + port)

	http.HandleFunc("/", createProxyHandler(routes, notFound, fileLogger))

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
