package client

import (
	"errors"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type httpClient interface {
	fetchData(url string) (*goquery.Document, error)
}

// RealHttpClient is a wrapper to make real HTTP requests.
type RealHttpClient struct{}

func (c *RealHttpClient) fetchData(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("query: status code is not 200")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.New("")
	}

	return doc, nil
}
