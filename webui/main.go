package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"tact-webui/api"
	"tact-webui/handlers"
)

func main() {
	// Configuration flags
	port := flag.Int("port", 2200, "Port to run the web server on")
	apiURL := flag.String("api", "", "Backend API URL (default: http://localhost:2100)")
	flag.Parse()

	// Determine API URL (flag > env > default)
	backendURL := *apiURL
	if backendURL == "" {
		backendURL = os.Getenv("TACT_API_URL")
	}
	if backendURL == "" {
		backendURL = "http://localhost:2100"
	}

	// Initialize API client
	client := api.NewClient(backendURL)

	// Create router and register handlers
	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, client)

	// Start server
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting TACT Web UI on http://localhost%s", addr)
	log.Printf("Backend API: %s", backendURL)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
