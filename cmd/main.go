package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	http_probe "github.com/miihael/go-http-probe"
)

func main() {
	fTimeout := flag.Duration("timeout", 30*time.Second, "selection timeout")
	flag.Parse()
	if flag.NArg() < 1 {
		panic("no URLs provided")
	}
	res, err := http_probe.SelectAll(flag.Args(), *fTimeout, nil)
	fmt.Fprintf(os.Stderr, "%d responses: error=%v\n", len(res), err)
	j := json.NewEncoder(os.Stdout)
	j.Encode(res)
}
