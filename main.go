package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	internalClient "github.com/tommy-muehle/speakerdeck_exporter/client"
	internalCollector "github.com/tommy-muehle/speakerdeck_exporter/collector"
	internalHttp "github.com/tommy-muehle/speakerdeck_exporter/http"
)

var (
	version   string = "latest"
	gitCommit string = "undefined"
)

var (
	addr  = flag.String("addr", ":9887", "An address to listen on for web interface and telemetry.")
	users = flag.String("users", "tommymuehle,fabpot", "Comma separated list of Speakerdeck users to watch.")
	v     = flag.Bool("version", false, "Prints current version")
)

func main() {
	flag.Parse()

	if *v {
		fmt.Println(version + "-" + gitCommit)
		os.Exit(0)
	}

	client := internalClient.NewSpeakerDeckClient(&internalClient.RealHttpClient{})
	collector := internalCollector.NewSpeakerDeckCollector(client, *users)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-signalChan
		log.Println("Shutdown exporter ...")
		os.Exit(0)
	}()

	server := internalHttp.NewServer(*addr)
	server.AddHandler("/", internalHttp.DefaultHandler())
	server.AddHandler("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	server.AddHandler("/favicon", http.NotFoundHandler())

	log.Printf("SpeakerDeck Prometheus exporter has successfully started")
	server.ListenAndServe()
}
