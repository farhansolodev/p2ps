package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var peers = make(map[string]string)

func handlePostedPeer(w http.ResponseWriter, r *http.Request) {
	// Read the body as a raw byte slice
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error reading request body: %s", err.Error())
		return
	}

	// Get the client's source IP address
	clientIP := r.RemoteAddr
	if colonIndex := strings.LastIndex(clientIP, ":"); colonIndex != -1 {
		clientIP = clientIP[:colonIndex] // Remove port from the address
	}

	// Store the client's IP and POST body in the global map
	peers[clientIP] = string(body)

	fmt.Printf("Current peers map: %+v\n", peers)
}

func main() {
    http.HandleFunc("/", handlePostedPeer)
    fmt.Println("HTTP server starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}