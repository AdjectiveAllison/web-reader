package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/AdjectiveAllison/web-reader/handler"
	"github.com/syumai/workers"
	_ "github.com/syumai/workers/cloudflare/d1"
)

func main() {
	db, err := sql.Open("d1", "DB")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	convertHandler := handler.NewConvertHandler(db)
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.RequestURI)
		
		if r.Method == http.MethodGet {
			if r.URL.Path == "" || r.URL.Path == "/" {
				http.Error(w, "URL is required", http.StatusBadRequest)
				return
			}
			
			// Extract URL from path
			targetURL := strings.TrimPrefix(r.URL.Path, "/")
			
			log.Printf("Extracted target URL: %s", targetURL)
			
			// Create a new request to pass to the handler
			newReq := *r
			newReq.URL.RawQuery = fmt.Sprintf("url=%s", targetURL)
			
			log.Printf("Modified request with query: %s", newReq.URL.RawQuery)
			convertHandler.ServeHTTP(w, &newReq)
			return
		}
		
		convertHandler.ServeHTTP(w, r)
	})

	workers.Serve(handler)
}
