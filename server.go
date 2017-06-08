package main

import (
	"flag"
	"log"
	"net/http"

	httpgzip "github.com/daaku/go.httpgzip"
	"github.com/dghubble/sessions"
)

var (
	bind  = flag.String("http", ":80", "")
	appid = flag.Int64("appid", 442422, "")
	token = flag.String("token", "FC144thCRZpyHM7qGDjt", "")

	store = sessions.NewCookieStore([]byte("t8yG1gmeEyeb7pQpw544UeCTyDfPkE6u"), []byte("Q599vrruZRhLFC144thCRZpyHM7qGDjt"))
)

func init() {
	flag.Parse()
	store.Config.MaxAge = 3600 * 6
}

func main() {
	log.Fatal(http.ListenAndServe(*bind, httpgzip.NewHandler(http.HandlerFunc(handler))))
}
