package http

import (
	"log"
	"net/http"
)

// DefaultHandler returns a simple request handler
// that replies to each request with an info page.
func DefaultHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(
			`<html>
				<head><title>SpeakerDeck Exporter</title></head>
				<body>
					<h1>SpeakerDeck Exporter</h1>
					<p><a href='/metrics'>Metrics</a></p>
				</body>
			</html>`))
		if err != nil {
			log.Printf("error while sending a response for the '/' path: %v", err)
		}
	})
}
