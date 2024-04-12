package scraper_test

import (
	"testing"

	"github.com/antsanchez/go-download-web/pkg/console"
	"github.com/antsanchez/go-download-web/pkg/get"
	"github.com/antsanchez/go-download-web/pkg/scraper"
)

type Path struct {
	Original string
	Folder   string
	Filename string
}

// Test for PreparePathsFile
func TestPreparePathsFile(t *testing.T) {
	var paths = []Path{
		{Original: "https://example.com/path/to/file.txt", Folder: "/path/to/", Filename: "file.txt"},
		{Original: "https://example.com/path/to/file.txt/", Folder: "/path/to/", Filename: "file.txt"},
		{Original: "https://example.com/path/to/file", Folder: "/path/to/", Filename: "file"},
		{Original: "https://example.com/path/to/file/", Folder: "/path/to/", Filename: "file"},
		{Original: "https://example.com/path/to/", Folder: "/path/", Filename: "to"},
		{Original: "https://example.com/path/to", Folder: "/path/", Filename: "to"},
		{Original: "https://example.com/path/", Folder: "/", Filename: "path"},
		{Original: "https://example.com/site.webmanifest", Folder: "/", Filename: "site.webmanifest"},
		{Original: "https://example.com/file.mp3", Folder: "/", Filename: "file.mp3"},
		{Original: "https://example.com/path/paht3/file.mp3", Folder: "/path/paht3/", Filename: "file.mp3"},
		{Original: "https://example.com/path/paht3/file.mp3", Folder: "/path/paht3/", Filename: "file.mp3"},
	}

	// Create a new scraper
	s := scraper.New(&scraper.Conf{OldDomain: "https://example.com"}, get.New(), console.New())

	for _, path := range paths {
		folder, filename := s.PreparePathsFile(path.Original)
		if folder != path.Folder || filename != path.Filename {
			t.Errorf("Expected %s, %s, got %s, %s on %s", path.Folder, path.Filename, folder, filename, path.Original)
		}
	}
}

// Test for PreparePathsPage
func TestPreparePathsPage(t *testing.T) {
	var paths = []Path{
		{Original: "https://example.com/path/to/index.html", Folder: "/path/to/", Filename: "index.html"},
		{Original: "https://example.com/path/to/index.html/", Folder: "/path/to/", Filename: "index.html"},
		{Original: "https://example.com/index.html", Folder: "/", Filename: "index.html"},
		{Original: "https://example.com/", Folder: "/", Filename: "index.html"},
		{Original: "https://example.com", Folder: "/", Filename: "index.html"},
		{Original: "https://example.com/path/to/file", Folder: "/path/to/file/", Filename: "index.html"},
		{Original: "https://example.com/path/to/file/", Folder: "/path/to/file/", Filename: "index.html"},
		{Original: "https://example.com/path/to/", Folder: "/path/to/", Filename: "index.html"},
		{Original: "https://example.com/path/to/test.html", Folder: "/path/to/", Filename: "test.html"},
		{Original: "https://example.com/path/to/feed.php", Folder: "/path/to/", Filename: "feed.php"},
		{Original: "https://example.com/feed.php", Folder: "/", Filename: "feed.php"},
		{Original: "https://example.com/moviecam.shtml", Folder: "/", Filename: "moviecam.shtml"},
	}

	// Create a new scraper
	s := scraper.New(&scraper.Conf{OldDomain: "https://example.com"}, get.New(), console.New())

	for _, path := range paths {
		folder, filename := s.PreparePathsPage(path.Original)
		if folder != path.Folder || filename != path.Filename {
			t.Errorf("Expected %s, %s, got %s, %s on %s", path.Folder, path.Filename, folder, filename, path.Original)
		}
	}
}
