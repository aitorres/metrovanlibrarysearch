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
	format := flag.String("format", "", "filter by BiblioCommons format code (e.g. BK, EBOOK, AB, DVD)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [--json] [--limit N] [--format CODE] \"<query>\"\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	var positional []string
	for flag.NArg() > 0 {
		positional = append(positional, flag.Arg(0))
		if err := flag.CommandLine.Parse(flag.Args()[1:]); err != nil {
			os.Exit(2)
		}
	}

	if len(positional) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	query := strings.TrimSpace(strings.Join(positional, " "))
	if query == "" {
		flag.Usage()
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), overallTimeout)
	defer cancel()

	reports := searchAll(ctx, query, *limit, *format)

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
func searchAll(ctx context.Context, query string, limit int, format string) []LibraryReport {
	reports := make([]LibraryReport, len(Libraries))
	var wg sync.WaitGroup
	for i, lib := range Libraries {
		wg.Add(1)
		go func(idx int, lib Library) {
			defer wg.Done()
			report := LibraryReport{Library: lib.Name, Results: []Result{}}
			results, err := lib.Adapter.Search(ctx, query, limit, format)
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
