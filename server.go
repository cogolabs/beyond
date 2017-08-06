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
	host   = flag.String("host", "beyond.colofoo.net", "")
	maxAge = flag.Int("max-age", 3600*6, "")

	cookiekey1 = flag.String("cookie-key1", "t8yG1gmeEyeb7pQpw544UeCTyDfPkE6u", "keypair 1 for cookie crypto")
	cookiekey2 = flag.String("cookie-key2", "Q599vrruZRhLFC144thCRZpyHM7qGDjt", "keypair 2 for cookie crypto")

	store *sessions.CookieStore
)

func init() {
	flag.Parse()
	store = sessions.NewCookieStore([]byte(*cookiekey1), []byte(*cookiekey2))
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
	err := setup()
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatal(http.ListenAndServe(*bind, http.HandlerFunc(handler)))
}

func setup() error {
	err := refreshFence()
	if err == nil {
		err = refreshSites()
	}
	if err == nil {
		err = refreshWhitelist()
	}
	if err == nil {
		err = reproxy()
	}
	return err
}
