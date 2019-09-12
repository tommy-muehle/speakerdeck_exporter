package client

import (
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

type testHttpClient struct{}

func (t testHttpClient) fetchData(url string) (*goquery.Document, error) {
	file, err := os.Open("testdata/response.html")
	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromReader(file)
}

func TestSpeakerDeckClient_GetStats(t *testing.T) {
	assert := assert.New(t)

	client := NewSpeakerDeckClient(&testHttpClient{})
	stats, err := client.GetResponses("tommymuehle")

	assert.NoError(err)
	assert.Equal(5, len(stats))
}
