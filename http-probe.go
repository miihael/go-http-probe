package http_probe

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

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
