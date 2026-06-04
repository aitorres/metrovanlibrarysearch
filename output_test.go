package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

type stubAdapter struct {
	results []Result
	err     error
}

func (s *stubAdapter) Search(_ context.Context, _ string, _ int) ([]Result, error) {
	return s.results, s.err
}

func TestRenderTextAndJSON(t *testing.T) {
	reports := []LibraryReport{
		{
			Library: "Test Library A",
			Results: []Result{{
				Title: "Foo", Author: "Bar", Format: "Book",
				CopiesTotal: 5, CopiesAvailable: 2,
				Description: "About a thing.",
				URL:         "https://example.com/r/1",
			}},
		},
		{Library: "Test Library B", Results: []Result{}},
		{Library: "Test Library C", Error: "boom"},
	}

	var buf bytes.Buffer
	renderText(&buf, "foo", reports)
	out := buf.String()
	for _, want := range []string{
		"Test Library A", "1. Foo", "by Bar", "2 available of 5", "https://example.com/r/1",
		"Test Library B", "No results.",
		"Test Library C", "Error: boom",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("text output missing %q\n%s", want, out)
		}
	}

	buf.Reset()
	if err := renderJSON(&buf, reports); err != nil {
		t.Fatal(err)
	}
	var round []LibraryReport
	if err := json.Unmarshal(buf.Bytes(), &round); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, buf.String())
	}
	if len(round) != 3 || round[0].Results[0].CopiesAvailable != 2 {
		t.Errorf("round-trip mismatch: %+v", round)
	}
}

func TestSearchAllOrderingAndErrors(t *testing.T) {
	// Swap out the registry for the duration of the test.
	saved := Libraries
	defer func() { Libraries = saved }()

	Libraries = []Library{
		{Name: "Lib1", Adapter: &stubAdapter{results: []Result{{Title: "A", URL: "u"}}}},
		{Name: "Lib2", Adapter: &stubAdapter{err: errors.New("fail")}},
		{Name: "Lib3", Adapter: &stubAdapter{results: nil}},
	}

	reports := searchAll(context.Background(), "q", 3)
	if len(reports) != 3 {
		t.Fatalf("want 3 reports, got %d", len(reports))
	}
	if reports[0].Library != "Lib1" || len(reports[0].Results) != 1 {
		t.Errorf("Lib1 wrong: %+v", reports[0])
	}
	if reports[1].Error != "fail" {
		t.Errorf("Lib2 error: %q", reports[1].Error)
	}
	if reports[2].Library != "Lib3" || len(reports[2].Results) != 0 || reports[2].Error != "" {
		t.Errorf("Lib3 wrong: %+v", reports[2])
	}
}
