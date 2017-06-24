package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"

	"github.com/dghubble/sessions"
	"github.com/koding/websocketproxy"
)

var (
	bind   = flag.String("http", ":80", "")
	appid  = flag.Int64("appid", 442422, "")
	domain = flag.String("domain", ".colofoo.net", "")
	maxAge = flag.Int("max-age", 3600*6, "")
	token  = flag.String("token", "FC144thCRZpyHM7qGDjt", "")

	store = sessions.NewCookieStore([]byte("t8yG1gmeEyeb7pQpw544UeCTyDfPkE6u"), []byte("Q599vrruZRhLFC144thCRZpyHM7qGDjt"))
)

func init() {
	flag.Parse()
	store.Config.Domain = *domain
	store.Config.MaxAge = *maxAge

	// allow insecure backends
	http.DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if websocketproxy.DefaultDialer.TLSClientConfig == nil {
		websocketproxy.DefaultDialer.TLSClientConfig = &tls.Config{}
	}
	websocketproxy.DefaultDialer.TLSClientConfig.InsecureSkipVerify = true
	websocketproxy.DefaultDialer.EnableCompression = true
	websocketproxy.DefaultUpgrader.EnableCompression = true
}

func websocketproxyDirector(incoming *http.Request, out http.Header) {
	out.Set("User-Agent", incoming.UserAgent())
	out.Set("X-Forwarded-Proto", "https")
}

func main() {
	log.Fatal(http.ListenAndServe(*bind, http.HandlerFunc(handler)))
}
