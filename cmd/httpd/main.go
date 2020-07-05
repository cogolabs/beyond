package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/cogolabs/beyond"
)

var (
	bind = flag.String("http", ":80", "listen address")

	srvReadTimeout  = flag.Duration("server-read-timeout", 1*time.Minute, "maximum duration for reading the entire request, including the body")
	srvWriteTimeout = flag.Duration("server-write-timeout", 2*time.Minute, "maximum duration before timing out writes of the response")
	srvIdleTimeout  = flag.Duration("server-idle-timeout", 3*time.Minute, "maximum amount of time to wait for the next request when keep-alives are enabled")
)

func main() {
	flag.Parse()

	if err := beyond.Setup(); err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    *bind,
		Handler: http.HandlerFunc(beyond.Handler),

		// https://blog.cloudflare.com/exposing-go-on-the-internet/
		ReadTimeout:  *srvReadTimeout,
		WriteTimeout: *srvWriteTimeout,
		IdleTimeout:  *srvIdleTimeout,
	}
	log.Fatal(srv.ListenAndServe())
}
