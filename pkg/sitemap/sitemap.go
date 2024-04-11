package sitemap

import (
	"encoding/xml"
	"os"
	"path/filepath"
)

const SitemapFile = "sitemap.xml"

// URLSitemap is the model for every url entry on the sitemap
type URLSitemap struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

// sitemapPath returns the path to the sitemap file
func sitemapPath(filename string) string {
	if filename != "" {
		return filepath.Join(filename, SitemapFile)
	}
	return SitemapFile
}

// CreateSitemap creates the sitemap
func CreateSitemap(links []string, filename string) error {
	filename = sitemapPath(filename)

	// Create the urlset
	urlset := &struct {
		XMLName xml.Name `xml:"urlset"`
		XMLNS   string   `xml:"xmlns,attr"`
		URLs    []URLSitemap
	}{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]URLSitemap, len(links)),
	}
	for i, link := range links {
		url := &URLSitemap{Loc: link}
		urlset.URLs[i] = *url
	}

	// Marshal the urlset to XML
	output, err := xml.MarshalIndent(urlset, "", "  ")
	if err != nil {
		return err
	}

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the XML declaration to the file
	_, err = file.WriteString(xml.Header)
	if err != nil {
		return err
	}

	// Write the urlset to the file
	_, err = file.Write(output)
	return err
}
