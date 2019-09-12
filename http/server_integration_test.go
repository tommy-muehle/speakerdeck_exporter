package http

import (
	"net/http"
	"testing"
	"time"
)

func TestServer_ListenAndServe(t *testing.T) {
	s := NewServer(":1235")
	s.AddHandler("/foo", http.NotFoundHandler())

	defer s.Shutdown()

	go func() {
		s.ListenAndServe()
	}()

	time.Sleep(1 * time.Second)

	res, err := http.Get("http://localhost:1235/foo")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status code isn't 404: %v", res.StatusCode)
	}
}
