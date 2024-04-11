package main

import (
	"errors"
	"flag"
	"log"

	"github.com/antsanchez/go-download-web/scraper"
)

// ParseFlags parses the flags
func parseFlags() (conf scraper.Conf, err error) {
	conf.OldDomain = *flag.String("u", "", "URL to copy")
	conf.NewDomain = *flag.String("new", "", "New URL")
	conf.Simultaneus = *flag.Int("s", 3, "Number of concurrent connections")
	conf.UseQueries = *flag.Bool("q", false, "Ignore queries on URLs")
	conf.Path = *flag.String("path", "./website", "Local path for downloaded files")
	flag.Parse()

	if conf.OldDomain == "" {
		err = errors.New("URL cannot be empty! Please, use '-u <URL>'")
		return
	}

	if conf.Simultaneus <= 0 {
		err = errors.New("the number of concurrent connections be at least 1'")
		return
	}

	log.Println("Domain:", conf.OldDomain)
	if conf.NewDomain != "" {
		log.Println("New Domain: ", conf.NewDomain)
	}
	log.Println("Simultaneus:", conf.Simultaneus)
	log.Println("Use Queries:", conf.UseQueries)

	return
}
