package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

var hostport = flag.String("l", "localhost:9090", "hostport to listen on")

func debugHandler(w http.ResponseWriter, r *http.Request) {
	b, err := httputil.DumpRequest(r, true)
	if err != nil {
		logErr(w, r, err, http.StatusInternalServerError)
		return
	}
	log.Printf("\u001b[35m>>> request dump follows\u001b[0m")
	mw := io.MultiWriter(os.Stderr, w)
	if _, err := mw.Write(append(b, []byte("\n")...)); err != nil {
		logErr(w, r, err, http.StatusInternalServerError)
		return
	}
	log.Printf("\u001b[34m>>> EOM\u001b[0m")
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
