// Basic URL check by trying to parse the string.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/miku/parallel"
)

var (
	numWorkers = flag.Int("w", runtime.NumCPU(), "number of workers")
	batchSize  = flag.Int("b", 20000, "batch size")
	debug      = flag.Bool("v", false, "output skipped URLs to stderr")
)

func main() {
	flag.Parse()
	pp := parallel.NewProcessor(os.Stdin, os.Stdout, func(p []byte) ([]byte, error) {
		s := strings.TrimSpace(string(p))
		_, err := url.Parse(string(p))
		if err == nil {
			return p, nil
		}
		if *debug {
			fmt.Fprintf(os.Stderr, "failed: %s\n", s)
		}
		return nil, nil
	})
	pp.BatchSize = *batchSize
	pp.NumWorkers = *numWorkers
	if err := pp.Run(); err != nil {
		log.Fatal(err)
	}
}
