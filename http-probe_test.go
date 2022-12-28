package http_probe_test

import (
	"net/url"
	"testing"
	"time"

	http_probe "github.com/miihael/go-http-probe"
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

func TestSelectURLs(t *testing.T) {
	urls := []url.URL{
		url.URL{Scheme: "https", Host: "debian.org"},
		url.URL{Scheme: "https", Host: "google.com"},
		url.URL{Scheme: "https", Host: "ubuntu.com"},
	}
	url, err := http_probe.SelectURLs(urls, 10*time.Second, nil)
	if err != nil {
		t.Error(err)
	}
	t.Logf("selected %s", url)
}

func TestSelectURLsIdx(t *testing.T) {
	urls := []url.URL{
		url.URL{Scheme: "https", Host: "debian.org"},
		url.URL{Scheme: "https", Host: "google.com"},
		url.URL{Scheme: "https", Host: "ubuntu.com"},
	}
	j, err := http_probe.SelectURLsIdx(urls, 10*time.Second, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("selected %#v", urls[j])
}

func TestSmallTimeout(t *testing.T) {
	urls := []url.URL{
		url.URL{Scheme: "https", Host: "debian.org"},
		url.URL{Scheme: "https", Host: "google.com"},
		url.URL{Scheme: "https", Host: "ubuntu.com"},
	}
	url, err := http_probe.SelectURLs(urls, 10*time.Millisecond, nil)
	if err != nil {
		t.Logf("expected error: %s", err)
	}
	if url != "" {
		t.Errorf("selected %s, must be empty", url)
	}
}
