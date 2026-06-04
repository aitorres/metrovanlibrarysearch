package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const userAgent = "metrovanlibsearch/1.0"

// httpClient is shared across adapters.
var httpClient = &http.Client{Timeout: 15 * time.Second}

// BiblioCommonsAdapter searches a BiblioCommons-powered catalog by hitting
// the public RSS gateway for the result list and the per-record availability
// page for copy counts.
type BiblioCommonsAdapter struct {
	Subdomain string
}

// rssFeed mirrors only the fields we read from the BiblioCommons RSS output.
type rssFeed struct {
	XMLName xml.Name  `xml:"rss"`
	Items   []rssItem `xml:"channel>item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Creator     string `xml:"http://purl.org/dc/elements/1.1/ creator"`
	Format      string `xml:"format"`
	Subtitle    string `xml:"subtitle"`
}

func (a *BiblioCommonsAdapter) Search(ctx context.Context, query string, limit int, format string) ([]Result, error) {
	feed, err := a.fetchRSS(ctx, query, format)
	if err != nil {
		return nil, err
	}
	if limit > len(feed.Items) {
		limit = len(feed.Items)
	}
	results := make([]Result, limit)

	var wg sync.WaitGroup
	for i := 0; i < limit; i++ {
		item := feed.Items[i]
		results[i] = Result{
			Title:       cleanTitle(item.Title, item.Subtitle),
			Author:      strings.TrimSpace(item.Creator),
			Format:      formatLabel(item.Format),
			Description: truncate(strings.TrimSpace(item.Description), 400),
			URL:         strings.TrimSpace(item.Link),
		}
		recordID := recordIDFromURL(results[i].URL)
		if recordID == "" {
			continue
		}
		wg.Add(1)
		go func(idx int, recID string) {
			defer wg.Done()
			total, avail, err := a.fetchAvailability(ctx, recID)
			if err != nil {
				return // leave zero counts on availability lookup failure
			}
			results[idx].CopiesTotal = total
			results[idx].CopiesAvailable = avail
		}(i, recordID)
	}
	wg.Wait()
	return results, nil
}

func (a *BiblioCommonsAdapter) fetchRSS(ctx context.Context, query, format string) (*rssFeed, error) {
	body, err := httpGet(ctx, buildRSSURL(a.Subdomain, query, format))
	if err != nil {
		return nil, err
	}
	var feed rssFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("parse RSS: %w", err)
	}
	return &feed, nil
}

// buildRSSURL builds the BiblioCommons RSS search URL, optionally filtering
// by a format code (e.g. "BK", "EBOOK"). Empty format means no filter.
func buildRSSURL(subdomain, query, format string) string {
	u := fmt.Sprintf(
		"https://gateway.bibliocommons.com/v2/libraries/%s/rss/search?query=%s&searchType=keyword",
		url.PathEscape(subdomain), url.QueryEscape(query),
	)
	if f := strings.ToUpper(strings.TrimSpace(format)); f != "" {
		u += "&f_FORMAT=" + url.QueryEscape(f)
	}
	return u
}

var (
	totalCopiesRe = regexp.MustCompile(`"totalCopies":\s*(\d+)`)
	availCopiesRe = regexp.MustCompile(`"availableCopies":\s*(\d+)`)
)

func (a *BiblioCommonsAdapter) fetchAvailability(ctx context.Context, recordID string) (total, available int, err error) {
	u := fmt.Sprintf("https://%s.bibliocommons.com/v2/availability/%s",
		a.Subdomain, url.PathEscape(recordID))
	body, err := httpGet(ctx, u)
	if err != nil {
		return 0, 0, err
	}
	return parseAvailability(body)
}

func parseAvailability(body []byte) (total, available int, err error) {
	if m := totalCopiesRe.FindSubmatch(body); m != nil {
		total, _ = strconv.Atoi(string(m[1]))
	}
	if m := availCopiesRe.FindSubmatch(body); m != nil {
		available, _ = strconv.Atoi(string(m[1]))
	}
	return total, available, nil
}

// recordIDFromURL extracts the record ID from a BiblioCommons URL
func recordIDFromURL(u string) string {
	const marker = "/v2/record/"
	i := strings.Index(u, marker)
	if i < 0 {
		return ""
	}
	id := u[i+len(marker):]
	if j := strings.IndexAny(id, "?#/"); j >= 0 {
		id = id[:j]
	}
	return id
}

// cleanTitle joins a title with its subtitle when the subtitle is meaningful
func cleanTitle(title, subtitle string) string {
	t := strings.TrimSpace(title)
	s := strings.TrimSpace(subtitle)
	if s == "" {
		return t
	}
	// Skip purely numeric subtitles (volume numbers etc.)
	allDigits := true
	for _, r := range s {
		if r < '0' || r > '9' {
			allDigits = false
			break
		}
	}
	if allDigits {
		return t
	}
	return t + ": " + s
}

// formatLabel converts BiblioCommons format codes to friendlier labels.
func formatLabel(code string) string {
	switch strings.ToUpper(strings.TrimSpace(code)) {
	case "":
		return ""
	case "BK":
		return "Book"
	case "EBOOK":
		return "eBook"
	case "AB":
		return "Audiobook"
	case "EAUDIO", "AUDIO":
		return "Audiobook"
	case "DVD":
		return "DVD"
	case "BLU_RAY":
		return "Blu-ray"
	case "MUSIC_CD":
		return "Music CD"
	case "MUSIC_ONLINE":
		return "Streaming Music"
	case "VIDEO_ONLINE":
		return "Streaming Video"
	case "COMIC_BK":
		return "Comic Book"
	case "GRAPHIC_NOVEL":
		return "Graphic Novel"
	case "MAG":
		return "Magazine"
	case "EMAG":
		return "eMagazine"
	default:
		return code
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	cut := s[:n]
	if i := strings.LastIndex(cut, " "); i > n/2 {
		cut = cut[:i]
	}
	return cut + "…"
}

func httpGet(ctx context.Context, u string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "*/*")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d for %s", resp.StatusCode, u)
	}
	return io.ReadAll(resp.Body)
}
