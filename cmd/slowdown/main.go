// slowdown will slow down printing.
//
// $ find ~ | slowdown -s 10ms
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	sleep     = flag.Duration("s", 10*time.Millisecond, "sleep after each char")
	sleepLine = flag.Duration("sl", 200*time.Millisecond, "sleep after line")
	br        = bufio.NewReader(os.Stdin)
)

func main() {
	flag.Parse()
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range line {
			fmt.Printf("%c", c)
			if c == '\n' {
				time.Sleep(*sleepLine)
			} else {
				time.Sleep(*sleep)
			}
		}
	}
}
