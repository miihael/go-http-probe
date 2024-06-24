package http_probe_test

import (
	"context"
	"net/http"
	"net/url"
	"os"
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
	u, err := http_probe.Select(urls, 10*time.Second, nil)
	if err != nil {
		t.Error(err)
	}
	t.Logf("selected %s", u)
}

func TestSelectWithContext(t *testing.T) {
	const stopAfter = 30
	ctx, cancel := context.WithTimeout(context.Background(), stopAfter*time.Millisecond)
	defer cancel()
	urls := []string{
		"http://debian.org",
		"http://google.com",
		"http://example.com",
	}
	tt := time.Now()
	u, err := http_probe.SelectWithContext(ctx, urls, 1*time.Second, &http.Client{})
	if err != nil && !os.IsTimeout(err) {
		t.Error(err)
	}
	elapsed := time.Since(tt)
	t.Logf("selected %q, elapsed %s, error %v", u, elapsed, err)
	if elapsed.Milliseconds() > stopAfter*1.1 {
		t.Errorf("too much time passed: %s, expected %dms", elapsed, stopAfter)
	}
}

func TestSelectURLs(t *testing.T) {
	urls := []url.URL{
		url.URL{Scheme: "https", Host: "debian.org"},
		url.URL{Scheme: "https", Host: "google.com"},
		url.URL{Scheme: "https", Host: "ubuntu.com"},
	}
	u, err := http_probe.SelectURLs(urls, 10*time.Second, nil)
	if err != nil {
		t.Error(err)
	}
	t.Logf("selected %s", u)
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

func TestSelectURLsIdxCtx(t *testing.T) {
	urls := []url.URL{
		url.URL{Scheme: "https", Host: "debian.org"},
		url.URL{Scheme: "https", Host: "google.com"},
		url.URL{Scheme: "https", Host: "ubuntu.com"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	tt := time.Now()
	j, err := http_probe.SelectURLsIdxWithContext(ctx, urls, 10*time.Second, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("selected %#v, elapsed: %s", urls[j], time.Since(tt))
}
