package scraper

import "time"

type Scraper struct {
	// Original domain
	OldDomain string

	// New domain to rewrite the download HTML sites
	NewDomain string

	// Root domain
	Root string

	// Path where to save the downloads
	Path string

	// Use args on URLs
	UseQueries bool

	// Number of concurrent queries
	Simultaneus int

	// Scanning now
	Scanning chan int

	// New links found
	NewLinks chan []Links

	// Pages to save
	Pages chan Page

	// Attachments found
	Attachments chan []string

	// Started
	Started chan int

	// Finished
	Finished chan int

	// Indexed pages
	Indexed []string

	// Pages for sitemap
	ForSitemap []string

	// Files to download
	Files []string

	// Seen links
	Seen map[string]bool

	// Start time
	StartTime time.Time
}

type Conf struct {
	// Original domain
	OldDomain string

	// New domain to rewrite the download HTML sites
	NewDomain string

	// Root domain
	Root string

	// Path where to save the downloads
	Path string

	// Use args on URLs
	UseQueries bool

	// Number of concurrent queries
	Simultaneus int
}

// Links model
type Links struct {
	Href string
}

// Page model
type Page struct {
	URL       string
	Canonical string
	Links     []Links
	NoIndex   bool
	HTML      string
}
