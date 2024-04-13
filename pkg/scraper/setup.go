package scraper

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Config holds the scraper configuration
type Config struct {
	// Original domain
	OldDomain string `long:"u" short:"u"`

	// New domain to rewrite the download HTML sites
	NewDomain string `long:"new" short:"new"`

	// URL prefixes/roots that should be included in the scraper
	IncludedURLs string `long:"r" short:"r"`

	// Roots contains a range of URLs that can be considered the root
	// This is useful for scraping sites where content is hosted on a CDN
	// Not a flag. This will be filled by the scraper uppon setup
	Roots []string

	// Path where to save the downloads
	DownloadPath string `long:"path" short:"path"`

	// Use args on URLs
	UseQueries bool `long:"q" short:"q"`

	// Number of concurrent queries
	Simultaneous int `long:"s" short:"s"`
}

// validateFlags ensures all required flags are set and values are valid
func validateFlags(conf *Config) error {
	if conf.OldDomain == "" {
		return errors.New("missing required flag: -u (URL)")
	}

	if conf.Simultaneous <= 0 {
		return errors.New("invalid number of connections: -s (must be at least 1)")
	}

	return nil
}

// parseFlags parses command line arguments and validates them
func ParseFlags() (*Config, error) {
	conf := &Config{
		Simultaneous: 3, // Set default value
		DownloadPath: "./website",
		UseQueries:   false,
	}
	flag.StringVar(&conf.OldDomain, "u", "", "URL to download content from. (required)")
	flag.StringVar(&conf.NewDomain, "new", "", "New URL to use for downloaded content (optional)")
	flag.StringVar(&conf.IncludedURLs, "r", "", "URL prefixes/root paths that should be included (optional)")
	flag.IntVar(&conf.Simultaneous, "s", conf.Simultaneous, "Number of concurrent connections (default: 3, minimum: 1)")
	flag.BoolVar(&conf.UseQueries, "q", conf.UseQueries, "Ignore query strings in URLs (optional)")
	flag.StringVar(&conf.DownloadPath, "path", conf.DownloadPath, "Local path to save downloaded files (default: ./website)")

	help := flag.Bool("h", false, "Show this help message")

	flag.Parse()

	if *help {
		PrintUsage()
		return nil, errors.New("help requested")
	}

	if err := validateFlags(conf); err != nil {
		PrintUsage()
		return nil, err
	}

	return conf, nil
}

// PrintUsage prints the usage message
func PrintUsage() {
	fmt.Println("Usage: ./go-download-web [options] -u <URL>")
	fmt.Println("Downloads website content and saves it locally.")
	flag.PrintDefaults()
}

// New creates a new Scraper
func New(conf *Config, getter HttpGet, con Console) (*Scraper, error) {

	con.AddStatus("Checking domain")

	// Get the root domain
	final, status, _, err := getter.Get(conf.OldDomain)
	if err != nil {
		return nil, fmt.Errorf("error getting domain: %s", err)
	}

	if status == http.StatusMovedPermanently || status == http.StatusFound {
		con.AddStatus(fmt.Sprintf("Redirected to %s", final))
	} else if status != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d on %s", status, conf.OldDomain)
	}

	correct := RemoveLastSlash(final)

	// Prepare the roots
	conf.Roots = append(conf.Roots, correct)
	if len(conf.IncludedURLs) > 0 {
		var urls = strings.Split(conf.IncludedURLs, ",")
		for _, url := range urls {
			if len(url) == 0 {
				continue
			}
			url = strings.TrimSpace(url)
			url = RemoveLastSlash(url)
			conf.Roots = append(conf.Roots, url)
		}
	}

	con.AddDomain(correct)

	con.AddStatus("Initiating scraper")

	return &Scraper{
		OldDomain:    conf.OldDomain,
		NewDomain:    conf.NewDomain,
		Roots:        conf.Roots,
		DownloadPath: conf.DownloadPath,
		UseQueries:   conf.UseQueries,

		Scanning:    make(chan int, conf.Simultaneous), // Semaphore
		NewLinks:    make(chan []Links, 100000),        // New links to scan
		Pages:       make(chan Page, 100000),           // Pages scanned
		Attachments: make(chan []string, 100000),       // Attachments
		Started:     make(chan int, 100000),            // Crawls started
		Finished:    make(chan int, 100000),            // Crawls finished

		Indexed:    []string{},
		ForSitemap: []string{},
		Files:      []string{},
		StartTime:  time.Now(),

		Seen: make(map[string]bool),

		Get: getter,
		Con: con,
	}, nil
}

// Close closes the channels
func (s *Scraper) Close() {
	close(s.Scanning)
	close(s.NewLinks)
	close(s.Pages)
	close(s.Attachments)
	close(s.Started)
	close(s.Finished)
}
