package main

import (
	"archive/tar"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/miku/parallel"
)

type Result struct {
	Filename string   `json:"f"`
	Errs     []string `json:"errs"`
	Contents []string `json:"c"`
}

func main() {
	pp := parallel.NewProcessor(os.Stdin, os.Stdout, func(p []byte) ([]byte, error) {
		filename := strings.TrimSpace(string(p))
		log.Println(filename)
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		result := Result{
			Filename: filename,
			Errs:     make([]string, 0),
			Contents: make([]string, 0),
		}
		tr := tar.NewReader(f)
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break // End of archive
			}
			if err != nil {
				result.Errs = append(result.Errs, err.Error())
				break
			}
			if hdr == nil {
				result.Errs = append(result.Errs, "unreadable tar")
				break
			}
			result.Contents = append(result.Contents, hdr.Name)
		}
		b, err := json.Marshal(result)
		b = append(b, []byte("\n")...)
		return b, err
	})
	pp.BatchSize = 10
	if err := pp.Run(); err != nil {
		log.Fatal(err)
	}
}
