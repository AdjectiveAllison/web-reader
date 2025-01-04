package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"syscall/js"
	"time"

	"github.com/AdjectiveAllison/web-reader/model"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
)

type ConvertHandler struct {
	db        *sql.DB
	converter *converter.Converter
}

func NewConvertHandler(db *sql.DB) *ConvertHandler {
	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
		),
	)

	return &ConvertHandler{
		db:        db,
		converter: conv,
	}
}

func (h *ConvertHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetRequest(w, r)
	case http.MethodPost:
		h.handlePostRequest(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ConvertHandler) handleGetRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handler received GET request with query: %s", r.URL.RawQuery)
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		log.Printf("No URL found in query parameters")
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}
	log.Printf("Processing URL: %s", targetURL)

	if err := h.processURL(w, targetURL, r.Header.Get("X-No-Cache") == "true"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ConvertHandler) handlePostRequest(w http.ResponseWriter, r *http.Request) {
	var req model.ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.processURL(w, req.URL, r.Header.Get("X-No-Cache") == "true"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ConvertHandler) processURL(w http.ResponseWriter, url string, bypassCache bool) error {
	log.Printf("Processing URL %s (bypass cache: %v)", url, bypassCache)
	
	// Check cache first if not bypassing
	if !bypassCache {
		log.Printf("Checking cache for URL: %s", url)
		var cache model.PageCache
		err := h.db.QueryRow("SELECT url, title, markdown, fetched_at FROM page_cache WHERE url = ?", url).Scan(
			&cache.URL, &cache.Title, &cache.Markdown, &cache.FetchedAt,
		)
		if err == nil {
			log.Printf("Cache hit for URL: %s (fetched at: %d)", url, cache.FetchedAt)
			return h.writeResponse(w, cache)
		}
		log.Printf("Cache miss for URL: %s (err: %v)", url, err)
	}

	// Use the Fetch API
	fetch := js.Global().Get("fetch")
	fetchPromise := fetch.Invoke(url)
	respValue := await(fetchPromise)
	if respValue.Get("ok").Bool() {
		bodyPromise := respValue.Call("text")
		html := await(bodyPromise).String()

		// Convert to markdown
		markdown, err := h.converter.ConvertString(html)
		if err != nil {
			return err
		}

		// Create cache entry
		cache := model.PageCache{
			URL:       url,
			Title:     respValue.Get("url").String(), // Use response URL
			Markdown:  markdown,
			FetchedAt: time.Now().Unix(),
		}

		// Store in cache
		log.Printf("Storing in cache: %s", cache.URL)
		_, err = h.db.Exec(
			"INSERT OR REPLACE INTO page_cache (url, title, markdown, fetched_at) VALUES (?, ?, ?, ?)",
			cache.URL, cache.Title, cache.Markdown, cache.FetchedAt,
		)
		if err != nil {
			return err
		}

		return h.writeResponse(w, cache)
	}
	return fmt.Errorf("failed to fetch URL: %s", respValue.Get("statusText").String())
}

// Helper to await JavaScript promises
func await(promise js.Value) js.Value {
	done := make(chan js.Value)
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		done <- args[0]
		return nil
	}))
	return <-done
}

func (h *ConvertHandler) writeResponse(w http.ResponseWriter, cache model.PageCache) error {
	// Format the response like Jina does
	output := fmt.Sprintf("Title: %s\n\nURL Source: %s\n\nMarkdown Content:\n%s\n",
		cache.Title,
		cache.URL,
		cache.Markdown,
	)

	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte(output))
	return err
}
