package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var namespace = "speakerdeck"

func newStateMetric() *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_state"),
		"State of the last metric scrape",
		nil,
		nil,
	)
}

func newDurationMetric() *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_duration_seconds"),
		"Duration of the last metric scrape",
		nil,
		nil,
	)
}

func newViewsMetric() *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "views"),
		"Total number of views for given deck",
		[]string{"deck", "author"},
		nil,
	)
}

func newStarsMetric() *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "stars"),
		"Total number of stars for given deck",
		[]string{"deck", "author"},
		nil,
	)
}
