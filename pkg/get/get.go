package get

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

type Get struct{}

func New() *Get {
	return &Get{}
}

// ParseURL parses a URL string and returns its components.
func (g *Get) ParseURL(baseURLString, relativeURLString string) (final string, err error) {
	baseURL, err := url.Parse(baseURLString)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Combine relative URL with base URL
	parsedURL, err := baseURL.Parse(relativeURLString)
	if err != nil {
		return "", fmt.Errorf("invalid relative URL: %w", err)
	}

	return parsedURL.String(), nil
}

func (g *Get) Get(link string) (final string, status int, buff *bytes.Buffer, err error) {
	resp, err := http.Get(link)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return
	}

	final = resp.Request.URL.String()
	return final, resp.StatusCode, buf, nil
}
