package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	httpPort  = ":80"
	httpsPort = ":443"
	certFile  = "/etc/letsencrypt/live/skyefactory.com/fullchain.pem"
	keyFile   = "/etc/letsencrypt/live/skyefactory.com/privkey.pem"
)

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

	handler := createProxyHandler(routes, notFound, fileLogger)

	// Check if certificates exist
	if _, err := os.Stat(certFile); err == nil {
		if _, err := os.Stat(keyFile); err == nil {
			// Start HTTPS server on port 443
			fmt.Println("Reverse proxy running on port 443")
			http.HandleFunc("/", handler)
			if err := http.ListenAndServeTLS(httpsPort, certFile, keyFile, nil); err != nil {
				log.Fatalf("HTTPS server failed: %v", err)
			}
		} else {
			log.Printf("Private key not found: %v", err)
		}
	} else {
		log.Printf("Certificate not found: %v", err)
	}
}
