// +build integration

package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FetchData(t *testing.T) {
	assert := assert.New(t)

	client := NewRealHttpClient()
	_, err := client.fetchData("https://raw.githubusercontent.com/tommy-muehle/speakerdeck_exporter/master/README.md")

	assert.NoError(err)
}

func Test_FetchDataWithTimeout(t *testing.T) {
	assert := assert.New(t)

	client := NewRealHttpClient()
	_, err := client.fetchData("http://httpstat.us/200?sleep=5000")

	assert.Error(err)
}

func Test_FetchDataWithInvalidResponse(t *testing.T) {
	assert := assert.New(t)

	client := NewRealHttpClient()
	_, err := client.fetchData("http://httpstat.us/500")

	assert.Error(err)
}
