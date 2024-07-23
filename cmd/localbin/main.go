package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sync/atomic"
)

var (
	hostport = flag.String("l", "localhost:9090", "hostport to listen on")
	dumpBody = flag.Bool("b", false, "dump body")

	count              int64 // count the number of requests
	totalContentLength int64 // sum of content length
	numErrors          int64 // how many failures
)

func debugHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&count, 1)
	atomic.AddInt64(&totalContentLength, r.ContentLength)
	b, err := httputil.DumpRequest(r, *dumpBody)
	if err != nil {
		logErr(w, r, err, http.StatusInternalServerError)
		atomic.AddInt64(&numErrors, 1)
		return
	}
	log.Printf("\u001b[35m--------8<-------- request dump follows [%d/%d/%d] --------8<--------\u001b[0m",
		atomic.LoadInt64(&count),
		atomic.LoadInt64(&totalContentLength),
		atomic.LoadInt64(&numErrors))
	mw := io.MultiWriter(os.Stderr, w)
	if _, err := mw.Write(append(b, []byte("\n")...)); err != nil {
		logErr(w, r, err, http.StatusInternalServerError)
		atomic.AddInt64(&numErrors, 1)
		return
	}
	log.Printf("\u001b[34m--------8<-------- EOM --------8<--------\u001b[0m")
}

func logErr(w http.ResponseWriter, r *http.Request, err error, status int) {
	log.Println(err)
	w.WriteHeader(status)
}

func main() {
	flag.Parse()
	log.Printf("point your client to http://%s", *hostport)
	log.Fatal(http.ListenAndServe(*hostport, http.HandlerFunc(debugHandler)))
}
