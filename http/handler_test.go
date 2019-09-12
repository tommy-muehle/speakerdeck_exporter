package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_DefaultHandler(t *testing.T) {
	ts := httptest.NewServer(DefaultHandler())
	defer ts.Close()

	time.Sleep(1 * time.Second)

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("status code isn't 200: %v", res.StatusCode)
	}
}
