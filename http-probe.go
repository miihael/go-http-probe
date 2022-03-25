package http_probe

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

type Probe struct {
	Url     string
	Elapsed time.Duration
}

//Select the fastest URL responded OK within timeout
func Select(urls []string, timeout time.Duration, client *http.Client) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if client == nil {
		client = &http.Client{
			Timeout: timeout - 100*time.Millisecond,
		}
	}

	dst := make(chan (string), len(urls))
	for _, u := range urls {
		go func(url string) {
			r, err := ctxhttp.Head(ctx, client, url)
			if err == nil {
				defer r.Body.Close()
				if r.StatusCode >= 200 && r.StatusCode < 400 {
					dst <- url
					cancel()
				}
			}
		}(u)
	}

	select {
	case u := <-dst:
		return u, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("Timeout")
	}
}

//SelectURLs returns the fastest URL from URL objects responded OK within timeout
func SelectURLs(urls []url.URL, timeout time.Duration, client *http.Client) (string, error) {
	strs := make([]string, len(urls))
	for i, u := range urls {
		strs[i] = u.String()
	}
	return Select(strs, timeout, client)
}

//SelectAll returns URLs and elapsed times in order they have replied
func SelectAll(urls []string, timeout time.Duration, client *http.Client) ([]Probe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if client == nil {
		client = &http.Client{
			Timeout: timeout - 100*time.Millisecond,
		}
	}

	dst := make(chan (Probe), len(urls))
	for _, u := range urls {
		go func(url string) {
			t := time.Now()
			r, err := ctxhttp.Head(ctx, client, url)
			if err == nil {
				defer r.Body.Close()
				if r.StatusCode >= 200 && r.StatusCode < 400 {
					dst <- Probe{url, time.Since(t)}
				}
			}
		}(u)
	}

	cnt := len(urls)
	res := make([]Probe, 0, cnt)
	for {
		select {
		case p := <-dst:
			res = append(res, p)
			cnt--
			if cnt < 1 {
				return res, nil
			}
		case <-ctx.Done():
			return res, fmt.Errorf("Timeout")
		}
	}
	return res, nil
}
