package main

import "context"

// Result is one catalog hit from a single library.
type Result struct {
	Title           string `json:"title"`
	Author          string `json:"author,omitempty"`
	Format          string `json:"format,omitempty"`
	Description     string `json:"description,omitempty"`
	URL             string `json:"url"`
	CopiesTotal     int    `json:"copies_total"`
	CopiesAvailable int    `json:"copies_available"`
}

// LibraryReport holds the outcome of searching one library.
type LibraryReport struct {
	Library string   `json:"library"`
	Results []Result `json:"results"`
	Error   string   `json:"error,omitempty"`
}

// Adapter implements catalog search for one library system.
type Adapter interface {
	Search(ctx context.Context, query string, limit int, format string) ([]Result, error)
}
