package scraper_test

import (
	"net/http"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/antsanchez/go-download-web/pkg/scraper"
	"github.com/stretchr/testify/assert"
)

func TestIsInternLink(t *testing.T) {
	s := scraper.New(&scraper.Conf{OldDomain: "http://example.com"})

	// Test intern links
	assert.True(t, s.IsInternLink("http://example.com/path"))
	assert.True(t, s.IsInternLink("http://example.com/path/"))
	assert.True(t, s.IsInternLink("http://example.com/path/file.html"))
	assert.True(t, s.IsInternLink("http://example.com/path/file.html?"))
	assert.True(t, s.IsInternLink("http://example.com/path/file.html?query=value"))
	assert.True(t, s.IsInternLink("http://example.com/path/file.html#fragment"))
	assert.True(t, s.IsInternLink("http://example.com/path/file.html?query=value#fragment"))
	assert.True(t, s.IsInternLink("http://example.com:8080/path"))

	// Test extern links
	assert.False(t, s.IsInternLink("http://other.com/path"))
	assert.False(t, s.IsInternLink("http://example.org/path"))
	assert.False(t, s.IsInternLink("https://example.com/path"))
	assert.False(t, s.IsInternLink("http://user@example.com/path"))
	assert.False(t, s.IsInternLink("http://hello.com/http://example.com/path"))
	assert.False(t, s.IsInternLink("http:// example.com/"))
}

func TestRemoveQuery(t *testing.T) {
	s := scraper.New(&scraper.Conf{OldDomain: "http://example.com"})

	// Test removing query
	assert.Equal(t, "http://example.com/path", s.RemoveQuery("http://example.com/path?"))
	assert.Equal(t, "http://example.com/path", s.RemoveQuery("http://example.com/path?query=value"))
	assert.Equal(t, "http://example.com/path", s.RemoveQuery("http://example.com/path?query=value#fragment"))

	// Test no query
	assert.Equal(t, "http://example.com/path", s.RemoveQuery("http://example.com/path"))
	assert.Equal(t, "http://example.com/path/", s.RemoveQuery("http://example.com/path/"))
	assert.Equal(t, "http://example.com/path#fragment", s.RemoveQuery("http://example.com/path#fragment"))
}

func TestIsStart(t *testing.T) {
	s := scraper.New(&scraper.Conf{OldDomain: "http://example.com"})

	// Test start links
	assert.True(t, s.IsStart("http://example.com"))

	// Test non-start links
	assert.False(t, s.IsStart("http://example.com/path"))
	assert.False(t, s.IsStart("http://example.org/path"))
	assert.False(t, s.IsStart("http://other.com"))
	assert.False(t, s.IsStart("http://example.org"))
	assert.False(t, s.IsStart("https://example.com"))
}

func TestSanitizeURL(t *testing.T) {
	s := scraper.New(&scraper.Conf{OldDomain: "http://example.com"})

	// Test valid links
	assert.Equal(t, "http://example.com/path/", s.SanitizeURL("http://example.com/path"))
	assert.Equal(t, "http://example.com/path/", s.SanitizeURL("http://example.com/path/"))
	assert.Equal(t, "http://example.com/path/", s.SanitizeURL("http://example.com/path#fragment"))

	// Test invalid links
	assert.Equal(t, "", s.SanitizeURL("mailto:user@example.com"))
	assert.Equal(t, "", s.SanitizeURL("javascript:alert('hello')"))
	assert.Equal(t, "", s.SanitizeURL("tel:+1234567890"))
	assert.Equal(t, "", s.SanitizeURL("whatsapp:+1234567890"))
	assert.Equal(t, "", s.SanitizeURL("callto:+1234567890"))
	assert.Equal(t, "", s.SanitizeURL("wtai://wp/mc;+1234567890"))
	assert.Equal(t, "", s.SanitizeURL("sms:+1234567890"))
	assert.Equal(t, "", s.SanitizeURL("market://details?id=com.example.app"))
	assert.Equal(t, "", s.SanitizeURL("geopoint:37.786971,-122.399677"))
	assert.Equal(t, "", s.SanitizeURL("ymsgr:sendim?example"))
	assert.Equal(t, "", s.SanitizeURL("msnim:chat?contact=example@hotmail.com"))
	assert.Equal(t, "", s.SanitizeURL("gtalk:chat?jid=example@gmail.com"))
	assert.Equal(t, "", s.SanitizeURL("skype:example?chat"))

	// Test query removal
	assert.Equal(t, "http://example.com/path/", s.SanitizeURL("http://example.com/path?"))
	assert.Equal(t, "http://example.com/path/", s.SanitizeURL("http://example.com/path?query=value"))
	assert.Equal(t, "http://example.com/path/", s.SanitizeURL("http://example.com/path?query=value#fragment"))

	// Test trailing slash
	assert.Equal(t, "http://example.com/path/", s.SanitizeURL("http://example.com/path"))
	assert.Equal(t, "http://example.com/path/", s.SanitizeURL("http://example.com/path/"))
	assert.Equal(t, "http://example.com/path/", s.SanitizeURL("http://example.com/path#fragment"))
}

func TestIsValidExtension(t *testing.T) {
	s := scraper.New(&scraper.Conf{OldDomain: "http://example.com"})

	// Test valid extensions
	assert.False(t, s.IsValidExtension("http://example.com/path/file.html"))
	assert.False(t, s.IsValidExtension("http://example.com/path/file.htm"))
	assert.False(t, s.IsValidExtension("http://example.com/path/file.shtml"))
	assert.False(t, s.IsValidExtension("http://example.com/path/file.php"))
	assert.False(t, s.IsValidExtension("http://example.com/path/file.asp"))
	assert.False(t, s.IsValidExtension("http://example.com/path/file.aspx"))
	assert.False(t, s.IsValidExtension("http://example.com/path/file.jsp"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.css"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.js"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.json"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.xml"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.svg"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.gif"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.jpg"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.jpeg"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.png"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.ico"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.txt"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.pdf"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.doc"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.docx"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.xls"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.xlsx"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.ppt"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.pptx"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.odt"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.ods"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.odp"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.zip"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.tar.gz"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.tgz"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.rar"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.7z"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.mp3"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.ogg"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.wav"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.aiff"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.mp4"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.mpeg"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.avi"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.wmv"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.flv"))
	assert.True(t, s.IsValidExtension("http://example.com/path/file.mkv"))

	// Test invalid extensions
	assert.False(t, s.IsValidExtension("http://example.com/path/file"))

}

func TestIsValidAttachment(t *testing.T) {
	s := scraper.New(&scraper.Conf{OldDomain: "http://example.com"})

	// Test valid attachments
	assert.True(t, s.IsValidAttachment("http://example.com/path/file.zip"))
	assert.True(t, s.IsValidAttachment("http://example.com/path/file.tar.gz"))
	assert.True(t, s.IsValidAttachment("http://example.com/path/file.tgz"))
	assert.True(t, s.IsValidAttachment("http://example.com/path/file.rar"))
	assert.True(t, s.IsValidAttachment("http://example.com/path/file.7z"))
	assert.True(t, s.IsValidAttachment("http://example.com/path/file.mp3"))
	assert.True(t, s.IsValidAttachment("http://example.com/path/file.ogg"))
	assert.True(t, s.IsValidAttachment("http://example.com/path/file.wav"))
	assert.True(t, s.IsValidAttachment("http://example.com/path/file.aiff"))
	assert.True(t, s.IsValidAttachment("http://example.com/path/file.mp4"))
	assert.True(t, s.IsValidAttachment("http://example.com/path/file.mpeg"))
	assert.True(t, s.IsValidAttachment("http://example.com/path/file.avi"))
	assert.True(t, s.IsValidAttachment("http://example.com/path/file.wmv"))
	assert.True(t, s.IsValidAttachment("http://example.com/path/file.flv"))
	assert.True(t, s.IsValidAttachment("http://example.com/path/file.mkv"))

	// Test invalid attachments
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.html"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.htm"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.shtml"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.php"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.asp"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.aspx"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.jsp"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.css"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.js"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.json"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.xml"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.svg"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.gif"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.jpg"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.jpeg"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.png"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.ico"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.txt"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.pdf"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.doc"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.docx"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.xls"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.xlsx"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.ppt"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.pptx"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.odt"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.ods"))
	assert.False(t, s.IsValidAttachment("http://example.com/path/file.odp"))
}

func TestDoesLinkExist(t *testing.T) {
	s := scraper.New(&scraper.Conf{OldDomain: "http://example.com"})

	// Test existing links
	links := []scraper.Links{
		{Href: "http://example.com/path/file.html"},
		{Href: "http://example.com/path/file.htm"},
		{Href: "http://example.com/path/file.shtml"},
		{Href: "http://example.com/path/file.php"},
		{Href: "http://example.com/path/file.asp"},
		{Href: "http://example.com/path/file.aspx"},
		{Href: "http://example.com/path/file.jsp"},
		{Href: "http://example.com/path/file.css"},
		{Href: "http://example.com/path/file.js"},
		{Href: "http://example.com/path/file.json"},
		{Href: "http://example.com/path/file.xml"},
		{Href: "http://example.com/path/file.svg"},
		{Href: "http://example.com/path/file.gif"},
		{Href: "http://example.com/path/file.jpg"},
		{Href: "http://example.com/path/file.jpeg"},
		{Href: "http://example.com/path/file.png"},
		{Href: "http://example.com/path/file.ico"},
		{Href: "http://example.com/path/file.txt"},
		{Href: "http://example.com/path/file.pdf"},
		{Href: "http://example.com/path/file.doc"},
		{Href: "http://example.com/path/file.docx"},
		{Href: "http://example.com/path/file.xls"},
		{Href: "http://example.com/path/file.xlsx"},
		{Href: "http://example.com/path/file.ppt"},
		{Href: "http://example.com/path/file.pptx"},
		{Href: "http://example.com/path/file.odt"},
		{Href: "http://example.com/path/file.ods"},
		{Href: "http://example.com/path/file.odp"},
		{Href: "http://example.com/path/file.zip"},
		{Href: "http://example.com/path/file.tar.gz"},
		{Href: "http://example.com/path/file.tgz"},
		{Href: "http://example.com/path/file.rar"},
		{Href: "http://example.com/path/file.7z"},
		{Href: "http://example.com/path/file.mp3"},
		{Href: "http://example.com/path/file.ogg"},
		{Href: "http://example.com/path/file.wav"},
		{Href: "http://example.com/path/file.aiff"},
		{Href: "http://example.com/path/file.mp4"},
		{Href: "http://example.com/path/file.mpeg"},
		{Href: "http://example.com/path/file.avi"},
		{Href: "http://example.com/path/file.wmv"},
		{Href: "http://example.com/path/file.flv"},
		{Href: "http://example.com/path/file.mkv"},
	}
	for _, link := range links {
		assert.True(t, s.DoesLinkExist(link, links))
	}

	// Test non-existing links
	nonExistingLinks := []scraper.Links{
		{Href: "http://example.com/path/file.html.html"},
		{Href: "http://example.com/path/file.html.htm"},
		{Href: "http://example.com/path/file.html.shtml"},
		{Href: "http://example.com/path/file.html.php"},
		{Href: "http://example.com/path/file.html.asp"},
		{Href: "http://example.com/path/file.html.aspx"},
		{Href: "http://example.com/path/file.html.jsp"},
		{Href: "http://example.com/path/file.html.css"},
		{Href: "http://example.com/path/file.html.js"},
		{Href: "http://example.com/path/file.html.json"},
		{Href: "http://example.com/path/file.html.xml"},
		{Href: "http://example.com/path/file.html.svg"},
		{Href: "http://example.com/path/file.html.gif"},
		{Href: "http://example.com/path/file.html.jpg"},
		{Href: "http://example.com/path/file.html.jpeg"},
		{Href: "http://example.com/path/file.html.png"},
		{Href: "http://example.com/path/file.html.ico"},
		{Href: "http://example.com/path/file.html.txt"},
		{Href: "http://example.com/path/file.html.pdf"},
		{Href: "http://example.com/path/file.html.doc"},
		{Href: "http://example.com/path/file.html.docx"},
		{Href: "http://example.com/path/file.html.xls"},
		{Href: "http://example.com/path/file.html.xlsx"},
		{Href: "http://example.com/path/file.html.ppt"},
		{Href: "http://example.com/path/file.html.pptx"},
		{Href: "http://example.com/path/file.html.odt"},
		{Href: "http://example.com/path/file.html.ods"},
		{Href: "http://example.com/path/file.html.odp"},
		{Href: "http://example.com/path/file.html.zip"},
		{Href: "http://example.com/path/file.html.tar.gz"},
		{Href: "http://example.com/path/file.html.tgz"},
		{Href: "http://example.com/path/file.html.rar"},
		{Href: "http://example.com/path/file.html.7z"},
		{Href: "http://example.com/path/file.html.mp3"},
		{Href: "http://example.com/path/file.html.ogg"},
		{Href: "http://example.com/path/file.html.wav"},
		{Href: "http://example.com/path/file.html.aiff"},
		{Href: "http://example.com/path/file.html.mp4"},
		{Href: "http://example.com/path/file.html.mpeg"},
		{Href: "http://example.com/path/file.html.avi"},
		{Href: "http://example.com/path/file.html.wmv"},
		{Href: "http://example.com/path/file.html.flv"},
		{Href: "http://example.com/path/file.html.mkv"},
	}
	for _, link := range nonExistingLinks {
		assert.False(t, s.DoesLinkExist(link, links))
	}
}

func TestIsURLInSlice(t *testing.T) {
	s := scraper.New(&scraper.Conf{OldDomain: "http://example.com"})

	// Test existing URLs
	urls := []string{
		"http://example.com/path/file.html",
		"http://example.com/path/file.htm",
		"http://example.com/path/file.shtml",
		"http://example.com/path/file.php",
		"http://example.com/path/file.asp",
		"http://example.com/path/file.aspx",
		"http://example.com/path/file.jsp",
		"http://example.com/path/file.css",
		"http://example.com/path/file.js",
		"http://example.com/path/file.json",
		"http://example.com/path/file.xml",
		"http://example.com/path/file.svg",
		"http://example.com/path/file.gif",
		"http://example.com/path/file.jpg",
		"http://example.com/path/file.jpeg",
		"http://example.com/path/file.png",
		"http://example.com/path/file.ico",
		"http://example.com/path/file.txt",
		"http://example.com/path/file.pdf",
		"http://example.com/path/file.doc",
		"http://example.com/path/file.docx",
		"http://example.com/path/file.xls",
		"http://example.com/path/file.xlsx",
		"http://example.com/path/file.ppt",
		"http://example.com/path/file.pptx",
		"http://example.com/path/file.odt",
		"http://example.com/path/file.ods",
		"http://example.com/path/file.odp",
		"http://example.com/path/file.zip",
		"http://example.com/path/file.tar.gz",
		"http://example.com/path/file.tgz",
		"http://example.com/path/file.rar",
		"http://example.com/path/file.7z",
		"http://example.com/path/file.mp3",
		"http://example.com/path/file.ogg",
		"http://example.com/path/file.wav",
		"http://example.com/path/file.aiff",
		"http://example.com/path/file.mp4",
		"http://example.com/path/file.mpeg",
		"http://example.com/path/file.avi",
		"http://example.com/path/file.wmv",
		"http://example.com/path/file.flv",
		"http://example.com/path/file.mkv",
	}
	for _, url := range urls {
		assert.True(t, s.IsURLInSlice(url, urls))
		assert.True(t, s.IsURLInSlice(url+"/", urls))
	}

	// Test non-existing URLs
	nonExistingURLs := []string{
		"http://example.com/path/file.html.html",
		"http://example.com/path/file.html.htm",
		"http://example.com/path/file.html.shtml",
		"http://example.com/path/file.html.php",
		"http://example.com/path/file.html.asp",
		"http://example.com/path/file.html.aspx",
		"http://example.com/path/file.html.jsp",
		"http://example.com/path/file.html.css",
		"http://example.com/path/file.html.js",
		"http://example.com/path/file.html.json",
		"http://example.com/path/file.html.xml",
		"http://example.com/path/file.html.svg",
		"http://example.com/path/file.html.gif",
		"http://example.com/path/file.html.jpg",
		"http://example.com/path/file.html.jpeg",
		"http://example.com/path/file.html.png",
		"http://example.com/path/file.html.ico",
		"http://example.com/path/file.html.txt",
		"http://example.com/path/file.html.pdf",
		"http://example.com/path/file.html.doc",
		"http://example.com/path/file.html.docx",
		"http://example.com/path/file.html.xls",
		"http://example.com/path/file.html.xlsx",
		"http://example.com/path/file.html.ppt",
		"http://example.com/path/file.html.pptx",
		"http://example.com/path/file.html.odt",
		"http://example.com/path/file.html.ods",
		"http://example.com/path/file.html.odp",
		"http://example.com/path/file.html.zip",
		"http://example.com/path/file.html.tar.gz",
		"http://example.com/path/file.html.tgz",
		"http://example.com/path/file.html.rar",
		"http://example.com/path/file.html.7z",
		"http://example.com/path/file.html.mp3",
		"http://example.com/path/file.html.ogg",
		"http://example.com/path/file.html.wav",
		"http://example.com/path/file.html.aiff",
		"http://example.com/path/file.html.mp4",
		"http://example.com/path/file.html.mpeg",
		"http://example.com/path/file.html.avi",
		"http://example.com/path/file.html.wmv",
		"http://example.com/path/file.html.flv",
		"http://example.com/path/file.html.mkv",
	}
	for _, url := range nonExistingURLs {
		assert.False(t, s.IsURLInSlice(url, urls))
		assert.False(t, s.IsURLInSlice(url+"/", urls))
	}
}

func TestIsLinkScanned(t *testing.T) {
	s := scraper.New(&scraper.Conf{OldDomain: "http://example.com"})

	// Test scanned links
	links := []string{
		"http://example.com/path/file.html",
		"http://example.com/path/file.htm",
		"http://example.com/path/file.shtml",
		"http://example.com/path/file.php",
		"http://example.com/path/file.asp",
		"http://example.com/path/file.aspx",
		"http://example.com/path/file.jsp",
		"http://example.com/path/file.css",
		"http://example.com/path/file.js",
		"http://example.com/path/file.json",
		"http://example.com/path/file.xml",
		"http://example.com/path/file.svg",
		"http://example.com/path/file.gif",
		"http://example.com/path/file.jpg",
		"http://example.com/path/file.jpeg",
		"http://example.com/path/file.png",
		"http://example.com/path/file.ico",
		"http://example.com/path/file.txt",
		"http://example.com/path/file.pdf",
		"http://example.com/path/file.doc",
		"http://example.com/path/file.docx",
		"http://example.com/path/file.xls",
		"http://example.com/path/file.xlsx",
		"http://example.com/path/file.ppt",
		"http://example.com/path/file.pptx",
		"http://example.com/path/file.odt",
		"http://example.com/path/file.ods",
		"http://example.com/path/file.odp",
		"http://example.com/path/file.zip",
		"http://example.com/path/file.tar.gz",
		"http://example.com/path/file.tgz",
		"http://example.com/path/file.rar",
		"http://example.com/path/file.7z",
		"http://example.com/path/file.mp3",
		"http://example.com/path/file.ogg",
		"http://example.com/path/file.wav",
		"http://example.com/path/file.aiff",
		"http://example.com/path/file.mp4",
		"http://example.com/path/file.mpeg",
		"http://example.com/path/file.avi",
		"http://example.com/path/file.wmv",
		"http://example.com/path/file.flv",
		"http://example.com/path/file.mkv",
	}
	for _, link := range links {
		assert.True(t, s.IsLinkScanned(link, links))
	}

	// Test non-scanned links
	nonScannedLinks := []string{
		"http://example.com/path/file.html.html",
		"http://example.com/path/file.html.htm",
		"http://example.com/path/file.html.shtml",
		"http://example.com/path/file.html.php",
		"http://example.com/path/file.html.asp",
		"http://example.com/path/file.html.aspx",
		"http://example.com/path/file.html.jsp",
		"http://example.com/path/file.html.css",
		"http://example.com/path/file.html.js",
		"http://example.com/path/file.html.json",
		"http://example.com/path/file.html.xml",
		"http://example.com/path/file.html.svg",
		"http://example.com/path/file.html.gif",
		"http://example.com/path/file.html.jpg",
		"http://example.com/path/file.html.jpeg",
		"http://example.com/path/file.html.png",
		"http://example.com/path/file.html.ico",
		"http://example.com/path/file.html.txt",
		"http://example.com/path/file.html.pdf",
		"http://example.com/path/file.html.doc",
		"http://example.com/path/file.html.docx",
		"http://example.com/path/file.html.xls",
		"http://example.com/path/file.html.xlsx",
		"http://example.com/path/file.html.ppt",
		"http://example.com/path/file.html.pptx",
		"http://example.com/path/file.html.odt",
		"http://example.com/path/file.html.ods",
		"http://example.com/path/file.html.odp",
		"http://example.com/path/file.html.zip",
		"http://example.com/path/file.html.tar.gz",
		"http://example.com/path/file.html.tgz",
		"http://example.com/path/file.html.rar",
		"http://example.com/path/file.html.7z",
		"http://example.com/path/file.html.mp3",
		"http://example.com/path/file.html.ogg",
		"http://example.com/path/file.html.wav",
		"http://example.com/path/file.html.aiff",
		"http://example.com/path/file.html.mp4",
		"http://example.com/path/file.html.mpeg",
		"http://example.com/path/file.html.avi",
		"http://example.com/path/file.html.wmv",
		"http://example.com/path/file.html.flv",
		"http://example.com/path/file.html.mkv",
	}
	for _, link := range nonScannedLinks {
		assert.False(t, s.IsLinkScanned(link, links))
	}
}

func TestGetURLEmbedded(t *testing.T) {
	s := scraper.New(&scraper.Conf{OldDomain: "http://example.com"})

	// Test valid URLs
	validURLs := []string{
		`url("http://example.com/path/file.html")`,
		`url('http://example.com/path/file.htm')`,
		`url("http://example.com/path/file.shtml")`,
		`url('http://example.com/path/file.php')`,
		`url("http://example.com/path/file.asp")`,
		`url('http://example.com/path/file.aspx')`,
		`url("http://example.com/path/file.jsp")`,
		`url('http://example.com/path/file.css')`,
		`url("http://example.com/path/file.js")`,
		`url('http://example.com/path/file.json')`,
		`url("http://example.com/path/file.xml")`,
		`url('http://example.com/path/file.svg')`,
		`url("http://example.com/path/file.gif")`,
		`url('http://example.com/path/file.jpg')`,
		`url("http://example.com/path/file.jpeg")`,
		`url('http://example.com/path/file.png')`,
		`url("http://example.com/path/file.ico")`,
		`url('http://example.com/path/file.txt')`,
		`url("http://example.com/path/file.pdf")`,
		`url('http://example.com/path/file.doc')`,
		`url("http://example.com/path/file.docx")`,
		`url('http://example.com/path/file.xls')`,
		`url("http://example.com/path/file.xlsx")`,
		`url('http://example.com/path/file.ppt')`,
		`url("http://example.com/path/file.pptx")`,
		`url('http://example.com/path/file.odt')`,
		`url("http://example.com/path/file.ods")`,
		`url('http://example.com/path/file.odp')`,
		`url("http://example.com/path/file.zip")`,
		`url('http://example.com/path/file.tar.gz')`,
		`url("http://example.com/path/file.tgz")`,
		`url('http://example.com/path/file.rar')`,
		`url("http://example.com/path/file.7z")`,
		`url('http://example.com/path/file.mp3")`,
		`url("http://example.com/path/file.ogg")`,
		`url('http://example.com/path/file.wav")`,
		`url("http://example.com/path/file.aiff")`,
		`url('http://example.com/path/file.mp4")`,
		`url("http://example.com/path/file.mpeg")`,
		`url('http://example.com/path/file.avi")`,
		`url("http://example.com/path/file.wmv")`,
		`url('http://example.com/path/file.flv")`,
		`url("http://example.com/path/file.mkv")`,
	}
	for _, vu := range validURLs {
		u, err := url.Parse(s.GetURLEmbedded(vu))
		assert.NoError(t, err)
		assert.Equal(t, "http://example.com/path/file."+filepath.Ext(u.Path), u.String())
	}

	// Test invalid URLs
	invalidURLs := []string{
		`url("http://example.com/path/file.html.html")`,
		`url('http://example.com/path/file.html.htm")`,
		`url("http://example.com/path/file.html.shtml")`,
		`url('http://example.com/path/file.html.php")`,
		`url("http://example.com/path/file.html.asp")`,
		`url('http://example.com/path/file.html.aspx")`,
		`url("http://example.com/path/file.html.jsp")`,
		`url('http://example.com/path/file.html.css")`,
		`url("http://example.com/path/file.html.js")`,
		`url('http://example.com/path/file.html.json")`,
		`url("http://example.com/path/file.html.xml")`,
		`url('http://example.com/path/file.html.svg")`,
		`url("http://example.com/path/file.html.gif")`,
		`url('http://example.com/path/file.html.jpg")`,
		`url("http://example.com/path/file.html.jpeg")`,
		`url('http://example.com/path/file.html.png")`,
		`url("http://example.com/path/file.html.ico")`,
		`url('http://example.com/path/file.html.txt")`,
		`url("http://example.com/path/file.html.pdf")`,
		`url('http://example.com/path/file.html.doc")`,
		`url("http://example.com/path/file.html.docx")`,
		`url('http://example.com/path/file.html.xls")`,
		`url("http://example.com/path/file.html.xlsx")`,
		`url('http://example.com/path/file.html.ppt")`,
		`url("http://example.com/path/file.html.pptx")`,
		`url('http://example.com/path/file.html.odt")`,
		`url("http://example.com/path/file.html.ods")`,
		`url('http://example.com/path/file.html.odp")`,
		`url("http://example.com/path/file.html.zip")`,
		`url('http://example.com/path/file.html.tar.gz")`,
		`url("http://example.com/path/file.html.tgz")`,
		`url('http://example.com/path/file.html.rar")`,
		`url("http://example.com/path/file.html.7z")`,
		`url('http://example.com/path/file.html.mp3")`,
		`url("http://example.com/path/file.html.ogg")`,
		`url('http://example.com/path/file.html.wav")`,
		`url("http://example.com/path/file.html.aiff")`,
		`url('http://example.com/path/file.html.mp4")`,
		`url("http://example.com/path/file.html.mpeg")`,
		`url('http://example.com/path/file.html.avi")`,
		`url("http://example.com/path/file.html.wmv")`,
		`url('http://example.com/path/file.html.flv")`,
		`url("http://example.com/path/file.html.mkv")`,
	}
	for _, iu := range invalidURLs {
		u, err := url.Parse(s.GetURLEmbedded(iu))
		assert.Error(t, err)
		assert.Equal(t, "", u.String())
	}
}

func TestGetInsideAttachments(t *testing.T) {
	s := scraper.New(&scraper.Conf{OldDomain: "http://example.com"})

	// Test valid attachments
	validAttachments := []string{
		"http://example.com/path/file.zip",
		"http://example.com/path/file.tar.gz",
		"http://example.com/path/file.tgz",
		"http://example.com/path/file.rar",
		"http://example.com/path/file.7z",
		"http://example.com/path/file.mp3",
		"http://example.com/path/file.ogg",
		"http://example.com/path/file.wav",
		"http://example.com/path/file.aiff",
		"http://example.com/path/file.mp4",
		"http://example.com/path/file.mpeg",
		"http://example.com/path/file.avi",
		"http://example.com/path/file.wmv",
		"http://example.com/path/file.flv",
		"http://example.com/path/file.mkv",
	}
	for _, url := range validAttachments {
		resp, err := http.Get(url)
		assert.NoError(t, err)
		defer resp.Body.Close()

		attachments, err := s.GetInsideAttachments(url)
		assert.NoError(t, err)

		for _, attachment := range attachments {
			assert.True(t, s.IsValidAttachment(attachment))
		}
	}

	// Test invalid attachments
	invalidAttachments := []string{
		"http://example.com/path/file.html",
		"http://example.com/path/file.htm",
		"http://example.com/path/file.shtml",
		"http://example.com/path/file.php",
		"http://example.com/path/file.asp",
		"http://example.com/path/file.aspx",
		"http://example.com/path/file.jsp",
		"http://example.com/path/file.css",
		"http://example.com/path/file.js",
		"http://example.com/path/file.json",
		"http://example.com/path/file.xml",
		"http://example.com/path/file.svg",
		"http://example.com/path/file.gif",
		"http://example.com/path/file.jpg",
		"http://example.com/path/file.jpeg",
		"http://example.com/path/file.png",
		"http://example.com/path/file.ico",
		"http://example.com/path/file.txt",
		"http://example.com/path/file.pdf",
		"http://example.com/path/file.doc",
		"http://example.com/path/file.docx",
		"http://example.com/path/file.xls",
		"http://example.com/path/file.xlsx",
		"http://example.com/path/file.ppt",
		"http://example.com/path/file.pptx",
		"http://example.com/path/file.odt",
		"http://example.com/path/file.ods",
		"http://example.com/path/file.odp",
	}
	for _, url := range invalidAttachments {
		resp, err := http.Get(url)
		assert.NoError(t, err)
		defer resp.Body.Close()

		attachments, err := s.GetInsideAttachments(url)
		assert.NoError(t, err)

		assert.Empty(t, attachments)
	}
}
