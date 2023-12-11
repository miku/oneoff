// Basic URL check by trying to parse the string.
package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"runtime"

	"github.com/miku/parallel"
)

var (
	numWorkers = flag.Int("w", runtime.NumCPU(), "number of workers")
	batchSize  = flag.Int("b", 20000, "batch size")
)

func main() {
	flag.Parse()
	pp := parallel.NewProcessor(os.Stdin, os.Stdout, func(p []byte) ([]byte, error) {
		_, err := url.Parse(string(p))
		if err != nil {
			return p, nil
		}
		return nil, nil
	})
	pp.BatchSize = *batchSize
	pp.NumWorkers = *numWorkers
	if err := pp.Run(); err != nil {
		log.Fatal(err)
	}
}
