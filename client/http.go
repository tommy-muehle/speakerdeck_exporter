package client

import (
	"errors"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type httpClient interface {
	fetchData(url string) (*goquery.Document, error)
}

// RealHttpClient is a wrapper to make real HTTP requests.
type RealHttpClient struct {
	client http.Client
}

// NewRealHttpClient creates a RealHttpClient.
func NewRealHttpClient() *RealHttpClient {
	return &RealHttpClient{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *RealHttpClient) fetchData(url string) (*goquery.Document, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("fetchData: status code is not 200")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.New("fetchData: could not create goquery document")
	}

	return doc, nil
}
