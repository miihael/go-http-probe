# go-http-probe
Run HEAD request to several URLs to select the fastest

Example:

```
package main


import (
	"flag"
	"log"
	"time"

	"github.com/miihael/go-http-probe"
)

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatalln("Usage: %s url url [url...]")
	}

	url, err := http_probe.Select(flag.Args(), 15*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("selected %s", url)
}
```

Run it:

`go get github.com/miihael/go-http-probe`

`go run example.go`
