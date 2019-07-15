package http_probe_test

import (
	"testing"
	"time"

	"github.com/miihael/go-http-probe"
)

func TestSelect(t *testing.T) {
	urls := []string{
		"http://debian.org",
		"http://google.com",
		"http://example.com",
	}
	url, err := http_probe.Select(urls, 10*time.Second, nil)
	if err != nil {
		t.Error(err)
	}
	t.Logf("selected %s", url)
}
