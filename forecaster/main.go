package main

import (
	"log"
	"net"
	"net/http"
	"os"
)

var (
	apiBaseURL = os.Getenv("API_URL")
	apiKey     = os.Getenv("API_KEY")
	host       = os.Getenv("HOST")
	port       = os.Getenv("PORT")
)

func main() {
	server := NewServer(NewAPIClient(apiBaseURL, apiKey))
	addr := net.JoinHostPort(host, port)
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, server))
}
