package collector

import (
	"errors"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"

	"github.com/tommy-muehle/speakerdeck_exporter/client"
)

type fakeClient struct{}

func (c *fakeClient) GetResponses(users string) ([]*client.Response, error) {
	var responses []*client.Response

	responses = append(responses, &client.Response{
		Deck:   "foo",
		Author: "bar",
		Views:  1,
		Stars:  1,
	})

	return responses, nil
}

func TestSpeakerDeckCollector_Describe(t *testing.T) {
	assert := assert.New(t)

	ch := make(chan *prometheus.Desc, 10)
	done := make(chan bool, 1)

	defer func() {
		close(ch)
		close(done)
	}()

	go func() {
		collector := NewSpeakerDeckCollector(&fakeClient{}, "foo")
		collector.Describe(ch)

		done <- true
	}()

	<-done

	assert.Equal(4, len(ch))
}

func TestSpeakerDeckCollector_Collect(t *testing.T) {
	assert := assert.New(t)

	ch := make(chan prometheus.Metric, 10)
	done := make(chan bool, 1)

	defer func() {
		close(ch)
		close(done)
	}()

	go func() {
		collector := NewSpeakerDeckCollector(&fakeClient{}, "foo")
		collector.Collect(ch)

		done <- true
	}()

	<-done

	assert.Equal(4, len(ch))
}

type failClient struct{}

func (f failClient) GetResponses(users string) ([]*client.Response, error) {
	return nil, errors.New("no responses")
}

func TestSpeakerDeckCollector_CollectFailed(t *testing.T) {
	assert := assert.New(t)

	ch := make(chan prometheus.Metric, 10)
	done := make(chan bool, 1)

	defer func() {
		close(ch)
		close(done)
	}()

	go func() {
		collector := NewSpeakerDeckCollector(&failClient{}, "foo")
		collector.Collect(ch)

		done <- true
	}()

	<-done

	assert.Equal(2, len(ch))
}
