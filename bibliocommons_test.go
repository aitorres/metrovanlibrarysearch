package main
package main

import (
	"encoding/xml"
	"os"
	"strings"
	"testing"
)

func TestParseRSSFixture(t *testing.T) {
	data, err := os.ReadFile("testdata/vpl-rss.xml")
	if err != nil {
		t.Fatal(err)
	}
	var feed rssFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(feed.Items) == 0 {
		t.Fatal("expected items in feed")
	}
	first := feed.Items[0]
	if first.Title != "House of Leaves" {
		t.Errorf("title = %q, want %q", first.Title, "House of Leaves")
	}
	if first.Creator != "Danielewski, Mark Z." {
		t.Errorf("creator = %q", first.Creator)
	}
	if !strings.Contains(first.Link, "/v2/record/S38C737861") {
		t.Errorf("link = %q", first.Link)
	}
	if first.Format != "BK" {
		t.Errorf("format = %q", first.Format)
	}
}

func TestParseAvailabilityFixture(t *testing.T) {
	data, err := os.ReadFile("testdata/vpl-availability.html")
	if err != nil {
		t.Fatal(err)
	}
	total, avail, err := parseAvailability(data)
	if err != nil {
		t.Fatal(err)
	}
	if total != 22 {
		t.Errorf("total = %d, want 22", total)
	}
	if avail != 1 {
		t.Errorf("available = %d, want 1", avail)
	}
}

func TestRecordIDFromURL(t *testing.T) {
	cases := []struct{ in, want string }{
		{"https://vpl.bibliocommons.com/v2/record/S38C737861", "S38C737861"},
		{"https://burnaby.bibliocommons.com/v2/record/S75C2242373?ref=x", "S75C2242373"},
		{"https://example.com/foo", ""},
	}
	for _, c := range cases {
		got := recordIDFromURL(c.in)
		if got != c.want {
			t.Errorf("recordIDFromURL(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestCleanTitle(t *testing.T) {
	if got := cleanTitle("House of Five Leaves", "1"); got != "House of Five Leaves" {
		t.Errorf("numeric subtitle should be dropped, got %q", got)
	}
	if got := cleanTitle("The House of Blue Leaves", "A Play in Two Acts"); got != "The House of Blue Leaves: A Play in Two Acts" {
		t.Errorf("got %q", got)
	}
	if got := cleanTitle("Plain", ""); got != "Plain" {
		t.Errorf("got %q", got)
	}
}

func TestFormatLabel(t *testing.T) {
	if formatLabel("BK") != "Book" {
		t.Error("BK")
	}
	if formatLabel("COMIC_BK") != "Comic Book" {
		t.Error("COMIC_BK")
	}
	if formatLabel("UNKNOWN_FMT") != "UNKNOWN_FMT" {
		t.Error("unknown should pass through")
	}
}

func TestTruncate(t *testing.T) {
	short := "hello world"
	if got := truncate(short, 100); got != short {
		t.Errorf("got %q", got)
	}
	long := strings.Repeat("a ", 300)
	got := truncate(long, 50)
	if len(got) > 52 {
		t.Errorf("not truncated: %d", len(got))
	}
	if !strings.HasSuffix(got, "…") {
		t.Errorf("missing ellipsis: %q", got)
	}
}
