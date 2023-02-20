package httpext

import (
	"io"
	"net/http"
	"sync"
	"time"
)

var client *http.Client
var cSetUp sync.Once

func GetClient() *http.Client {
	cSetUp.Do(func() {
		client = &http.Client{
			Timeout: time.Minute,
		}
	})
	return client
}

func Get(url string) (*http.Response, error) {
	return GetClient().Get(url)
}

func Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	return GetClient().Post(url, contentType, body)
}

func Do(req *http.Request) (*http.Response, error) {
	return GetClient().Do(req)
}
