package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

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
