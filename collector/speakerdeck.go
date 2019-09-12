package collector

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/tommy-muehle/speakerdeck_exporter/client"
)

// SpeakerDeckCollector collects SpeakerDeck metrics. It implements prometheus.Collector interface.
type SpeakerDeckCollector struct {
	client client.Client
	users  string

	state    *prometheus.Desc
	duration *prometheus.Desc
	views    *prometheus.Desc
	stars    *prometheus.Desc

	mutex sync.Mutex
}

// NewSpeakerDeckCollector creates a SpeakerDeckCollector.
func NewSpeakerDeckCollector(client client.Client, users string) *SpeakerDeckCollector {
	return &SpeakerDeckCollector{
		client:   client,
		users:    users,
		duration: newDurationMetric(),
		state:    newStateMetric(),
		stars:    newStarsMetric(),
		views:    newViewsMetric(),
	}
}

// Describe sends the super-set of all possible descriptors of SpeakerDeck metrics
// to the provided channel.
func (c *SpeakerDeckCollector) Describe(channel chan<- *prometheus.Desc) {
	channel <- c.duration
	channel <- c.state
	channel <- c.stars
	channel <- c.views
}

// Collect fetches metrics from SpeakerDeck and sends them to the provided channel.
func (c *SpeakerDeckCollector) Collect(channel chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	begin := time.Now()
	stats, err := c.client.GetResponses(c.users)
	duration := time.Since(begin)

	if err != nil {
		channel <- prometheus.MustNewConstMetric(c.state, prometheus.GaugeValue, 0)
		channel <- prometheus.MustNewConstMetric(c.duration, prometheus.GaugeValue, duration.Seconds())
		return
	}

	channel <- prometheus.MustNewConstMetric(c.state, prometheus.GaugeValue, 1)
	channel <- prometheus.MustNewConstMetric(c.duration, prometheus.GaugeValue, duration.Seconds())

	for _, stat := range stats {
		channel <- prometheus.MustNewConstMetric(c.stars, prometheus.CounterValue, stat.Stars, stat.Deck, stat.Author)
		channel <- prometheus.MustNewConstMetric(c.views, prometheus.CounterValue, stat.Views, stat.Deck, stat.Author)
	}
}
