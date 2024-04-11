package get

import (
	"bytes"
	"net/http"
)

type GetInterface interface {
	Get(link string) (resp *http.Response, buff *bytes.Buffer, err error)
}

type Get struct{}

func New() *Get {
	return &Get{}
}

func (g *Get) Get(link string) (resp *http.Response, buff *bytes.Buffer, err error) {
	resp, err = http.Get(link)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return
	}

	return resp, buf, nil
}
