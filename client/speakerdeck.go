package client

import (
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const speakerDeckUrl = "https://speakerdeck.com/"

// Response represents provided data by SpeakerDeck.
type Response struct {
	Deck   string
	Author string
	Views  float64
	Stars  float64
}

// Client defines an interface for SpeakerDeck clients.
type Client interface {
	GetResponses(users string) ([]*Response, error)
}

// SpeakerDeckClient allows you to fetch SpeakerDeck metrics.
type SpeakerDeckClient struct {
	httpClient httpClient
}

// NewSpeakerDeckClient creates a SpeakerDeckClient.
func NewSpeakerDeckClient(client httpClient) *SpeakerDeckClient {
	return &SpeakerDeckClient{
		httpClient: client,
	}
}

// GetResponses returns responses for the given users.
func (c *SpeakerDeckClient) GetResponses(users string) ([]*Response, error) {
	urls := make([]string, 0, len(users))
	names := strings.Split(users, ",")

	for _, name := range names {
		urls = append(urls, speakerDeckUrl+name)
	}

	return c.parallelGet(urls)
}

func (c *SpeakerDeckClient) parallelGet(urls []string) ([]*Response, error) {
	var wg sync.WaitGroup
	var items []*Response

	for _, u := range urls {
		wg.Add(1)

		go func(u string) {
			defer wg.Done()

			doc, err := c.httpClient.fetchData(u)
			if err != nil {
				return
			}
			doc.Url = &url.URL{Path: strings.TrimPrefix(u, speakerDeckUrl)}

			items = append(items, c.grabStats(doc)...)
		}(u)
	}

	wg.Wait()

	return items, nil
}

func (c *SpeakerDeckClient) grabStats(doc *goquery.Document) []*Response {
	items := make([]*Response, 0)

	doc.Find(".row > .col-12").Each(func(_ int, s *goquery.Selection) {
		item := &Response{}
		item.Author = doc.Url.Path

		title, ok := s.Find("a").Attr("title")
		if !ok {
			return
		}

		item.Deck = title

		s.Find(".py-3").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 1: // stars
				item.Stars = c.unhumanizeCount(strings.TrimSpace(s.Text()))
			case 2: // views
				item.Views = c.unhumanizeCount(strings.TrimSpace(s.Text()))
			}
		})

		items = append(items, item)
	})

	return items
}

func (c *SpeakerDeckClient) unhumanizeCount(count string) float64 {
	var multiply bool = false

	if strings.HasSuffix(count, "k") {
		count = strings.TrimSuffix(count, "k")
		multiply = true
	}

	v, err := strconv.ParseFloat(count, 64)
	if err != nil {
		return 0
	}

	if multiply {
		return v * 1000
	}

	return v
}
