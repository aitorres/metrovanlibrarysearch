package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func renderJSON(w io.Writer, reports []LibraryReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(reports)
}

func renderText(w io.Writer, query string, reports []LibraryReport) {
	fmt.Fprintf(w, "Search results for %q across Metro Vancouver libraries\n", query)
	fmt.Fprintln(w, strings.Repeat("=", 70))

	for _, r := range reports {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "## %s\n", r.Library)
		fmt.Fprintln(w, strings.Repeat("-", len(r.Library)+3))

		switch {
		case r.Error != "":
			fmt.Fprintf(w, "  Error: %s\n", r.Error)
		case len(r.Results) == 0:
			fmt.Fprintln(w, "  No results.")
		default:
			for i, res := range r.Results {
				fmt.Fprintf(w, "  %d. %s\n", i+1, res.Title)
				if res.Author != "" {
					fmt.Fprintf(w, "     by %s\n", res.Author)
				}
				if res.Format != "" {
					fmt.Fprintf(w, "     Format: %s\n", res.Format)
				}
				fmt.Fprintf(w, "     Copies: %d available of %d total\n",
					res.CopiesAvailable, res.CopiesTotal)
				if res.Description != "" {
					fmt.Fprintf(w, "     %s\n", res.Description)
				}
				fmt.Fprintf(w, "     %s\n", res.URL)
				if i < len(r.Results)-1 {
					fmt.Fprintln(w)
				}
			}
		}
	}
}
