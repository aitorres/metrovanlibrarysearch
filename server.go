package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//go:embed index.html
var indexHTML []byte

const maxServeLimit = 25

func runServer(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/healthcheck", healthCheckAPIHandler)
	mux.HandleFunc("/api/search", searchAPIHandler)
	log.Printf("listening on http://localhost%s", normalizeAddr(addr))
	return http.ListenAndServe(addr, mux)
}

func normalizeAddr(addr string) string {
	if strings.HasPrefix(addr, ":") {
		return addr
	}
	return "/" + addr
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(indexHTML)
}

func searchAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		http.Error(w, "missing q parameter", http.StatusBadRequest)
		return
	}
	limit := defaultLimit
	if s := r.URL.Query().Get("limit"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			limit = n
		}
	}
	if limit > maxServeLimit {
		limit = maxServeLimit
	}
	format := strings.TrimSpace(r.URL.Query().Get("format"))

	ctx, cancel := context.WithTimeout(r.Context(), overallTimeout)
	defer cancel()

	reports := searchAll(ctx, q, limit, format)

	w.Header().Set("Content-Type", "application/json")
	if err := renderJSON(w, reports); err != nil {
		log.Printf("renderJSON error: %v", err)
		fmt.Fprintln(w)
	}
}

func healthCheckAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
