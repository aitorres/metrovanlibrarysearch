package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const defaultLimit = 3
const overallTimeout = 30 * time.Second

func main() {
	jsonOut := flag.Bool("json", false, "output JSON instead of human-readable text")
	limit := flag.Int("limit", defaultLimit, "max results per library")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [--json] [--limit N] \"<query>\"\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}
	query := strings.TrimSpace(strings.Join(flag.Args(), " "))
	if query == "" {
		flag.Usage()
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), overallTimeout)
	defer cancel()

	reports := searchAll(ctx, query, *limit)

	if *jsonOut {
		if err := renderJSON(os.Stdout, reports); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	} else {
		renderText(os.Stdout, query, reports)
	}
}

// searchAll fans out to every library and returns the reports in registry order.
func searchAll(ctx context.Context, query string, limit int) []LibraryReport {
	reports := make([]LibraryReport, len(Libraries))
	var wg sync.WaitGroup
	for i, lib := range Libraries {
		wg.Add(1)
		go func(idx int, lib Library) {
			defer wg.Done()
			report := LibraryReport{Library: lib.Name, Results: []Result{}}
			results, err := lib.Adapter.Search(ctx, query, limit)
			if err != nil {
				report.Error = err.Error()
			} else if results != nil {
				report.Results = results
			}
			reports[idx] = report
		}(i, lib)
	}
	wg.Wait()
	return reports
}
