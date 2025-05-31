//go:build dev

package app

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// StaticFSHandler serves static files from the filesystem during development
func StaticFSHandler() http.Handler {
	// Get current wokring directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// Build path to web/static (absolute)
	// Make sure to run main from root
	staticPath := filepath.Join(wd, "web", "static")

	log.Printf("DEV MODE: Serving static files from: %s", staticPath)
	return http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath)))
}
