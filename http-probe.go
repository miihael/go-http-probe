package http_probe

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

type Probe struct {
	Url     string
	Elapsed time.Duration
}

func errTimeout() error {
	return fmt.Errorf("probe timeout: %w", os.ErrDeadlineExceeded)
}

// Select the fastest URL responded OK within timeout
func SelectWithContext(pctx context.Context, urls []string, timeout time.Duration, client *http.Client) (string, error) {
	ctx, cancel := context.WithCancel(pctx)
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
	case <-pctx.Done():
		return "", pctx.Err()
	case u := <-dst:
		return u, nil
	case <-time.After(timeout):
		return "", errTimeout()
	}
}

func Select(urls []string, timeout time.Duration, client *http.Client) (string, error) {
	return SelectWithContext(context.Background(), urls, timeout, client)
}

// SelectURLs returns the fastest URL as string from URL objects responded OK within timeout
func SelectURLs(urls []url.URL, timeout time.Duration, client *http.Client) (string, error) {
	strs := make([]string, len(urls))
	for i, u := range urls {
		strs[i] = u.String()
	}
	return Select(strs, timeout, client)
}

// SelectURLsIdx returns the index of fastest URL from URL objects responded OK within timeout or -1 if none
func SelectURLsIdx(urls []url.URL, timeout time.Duration, client *http.Client) (int, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if client == nil {
		client = &http.Client{
			Timeout: timeout - 100*time.Millisecond,
		}
	}

	dst := make(chan (int), len(urls))
	for i, u := range urls {
		go func(j int, ur *url.URL) {
			req := &http.Request{
				Method:     "HEAD",
				URL:        ur,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     make(http.Header),
				Body:       nil,
				Host:       ur.Host,
			}
			r, err := ctxhttp.Do(ctx, client, req)
			if err == nil {
				defer r.Body.Close()
				if r.StatusCode >= 200 && r.StatusCode < 400 {
					dst <- j
					cancel()
				}
			}
		}(i, &u)
	}

	select {
	case j := <-dst:
		return j, nil
	case <-time.After(timeout):
		return -1, errTimeout()
	}
}

// SelectAll returns URLs and elapsed times in order they have replied
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
			return res, errTimeout()
		}
	}
}

// SelectURLsIdxWithContext returns the index of fastest URL from URLs responded OK within timeout or -1 if none
func SelectURLsIdxWithContext(pctx context.Context, urls []url.URL, timeout time.Duration, client *http.Client) (int, error) {
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	if client == nil {
		client = &http.Client{
			Timeout: timeout - 100*time.Millisecond,
		}
	}

	dst := make(chan (int), len(urls))
	for i, u := range urls {
		go func(j int, ur url.URL) {
			req := &http.Request{
				Method:     "HEAD",
				URL:        &ur,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     make(http.Header),
				Body:       nil,
				Host:       ur.Host,
			}
			req = req.WithContext(ctx)
			r, err := client.Do(req)
			if err == nil {
				defer r.Body.Close()
				if r.StatusCode >= 200 && r.StatusCode < 400 {
					dst <- j
					cancel()
				}
			}
		}(i, u)
	}

	select {
	case <-pctx.Done():
		return -1, pctx.Err()
	case j := <-dst:
		return j, nil
	case <-time.After(timeout):
		return -1, errTimeout()
	}
}
