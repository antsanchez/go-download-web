package scraper

import (
	"bytes"
	"time"
)

// Console interface
type Console interface {
	AddDomain(string)
	AddStatus(string)
	AddStarted()
	AddFinished()
	AddAttachments()
	AddDownloaded()
	AddDownloading()
	AddErrors(string)
}

// HttpGet interface
type HttpGet interface {
	ParseURL(baseURLString, relativeURLString string) (final string, err error)
	Get(link string) (final string, status int, buff *bytes.Buffer, err error)
}

type Scraper struct {
	// Original domain
	OldDomain string

	// New domain to rewrite the download HTML sites
	NewDomain string

	// Roots contains a range of URLs that can be considered the root
	// This is useful for scraping sites where content is hosted on a CDN
	Roots []string

	// Path where to save the downloads
	DownloadPath string

	// Use args on URLs
	UseQueries bool

	// Number of concurrent queries
	Simultaneous int

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

	// GetInterface
	Get HttpGet

	// Console
	Con Console
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
	HTML      string
}
