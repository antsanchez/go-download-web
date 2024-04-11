package main

import (
	"errors"
	"flag"
	"log"

	"github.com/antsanchez/go-download-web/scraper"
)

// ParseFlags parses the flags
func parseFlags() (conf scraper.Conf, err error) {
	flag.StringVar(&conf.OldDomain, "u", "", "URL to download")
	flag.StringVar(&conf.NewDomain, "new", "", "New URL")
	flag.StringVar(&conf.IncludedURLs, "r", "", "URL prefixes/root paths that should be included in the scraper, in addition to the domain")
	flag.IntVar(&conf.Simultaneous, "s", 3, "Number of concurrent connections")
	flag.BoolVar(&conf.UseQueries, "q", false, "Ignore queries on URLs")
	flag.StringVar(&conf.Path, "path", "./website", "Local path for downloaded files")
	flag.Parse()

	if conf.OldDomain == "" {
		err = errors.New("URL cannot be empty! Please, use '-u <URL>'")
		return
	}

	if conf.Simultaneous <= 0 {
		err = errors.New("the number of concurrent connections be at least 1'")
		return
	}

	log.Println("Domain:", conf.OldDomain)
	if conf.NewDomain != "" {
		log.Println("New Domain: ", conf.NewDomain)
	}
	log.Println("Simultaneous:", conf.Simultaneous)
	log.Println("Use Queries:", conf.UseQueries)

	return
}
