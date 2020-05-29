package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/cogolabs/beyond"
)

var (
	bind = flag.String("http", ":80", "listen address")
)

func main() {
	flag.Parse()

	if err := beyond.Setup(); err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(*bind, http.HandlerFunc(beyond.Handler)))
}
