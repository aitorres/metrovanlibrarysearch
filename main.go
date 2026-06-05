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

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s <command> [flags]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  query   search every library for a term")
	fmt.Fprintln(os.Stderr, "  serve   start a local web UI")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, "run '%s <command> --help' for command-specific flags\n", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "query":
		runQuery(os.Args[2:])
	case "serve":
		runServe(os.Args[2:])
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "error: unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func runQuery(args []string) {
	fs := flag.NewFlagSet("query", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "output JSON instead of human-readable text")
	limit := fs.Int("limit", defaultLimit, "max results per library")
	format := fs.String("format", "", "filter by BiblioCommons format code (e.g. BK, EBOOK, AB, DVD)")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s query [--json] [--limit N] [--format CODE] \"<query>\"\n", os.Args[0])
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	query := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if query == "" {
		fs.Usage()
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

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	addr := fs.String("addr", ":8080", "listen address")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s serve [--addr ADDR]\n", os.Args[0])
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	if fs.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "error: serve does not accept positional arguments")
		fs.Usage()
		os.Exit(2)
	}
	if err := runServer(*addr); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
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
