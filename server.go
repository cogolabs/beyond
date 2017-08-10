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
	bind = flag.String("http", ":80", "")
	host = flag.String("host", "beyond.colofoo.net", "")

	cookieAge  = flag.Int("cookie-age", 3600*6, "MaxAge setting in seconds")
	cookieDom  = flag.String("cookie-domain", ".colofoo.net", "")
	cookieKey1 = flag.String("cookie-key1", "t8yG1gmeEyeb7pQpw544UeCTyDfPkE6u", "keypair 1 for cookie crypto")
	cookieKey2 = flag.String("cookie-key2", "Q599vrruZRhLFC144thCRZpyHM7qGDjt", "keypair 2 for cookie crypto")

	skipVerify = flag.Bool("insecure-skip-verify", false, "allow TLS backends without valid certificates")
	wsCompress = flag.Bool("websocket-compression", false, "(gorilla/experimental)")

	store *sessions.CookieStore
)

func init() {
	flag.Parse()
	store = sessions.NewCookieStore([]byte(*cookieKey1), []byte(*cookieKey2))
	store.Config.Domain = *cookieDom
	store.Config.MaxAge = *cookieAge

	// allow insecure backends
	http.DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: *skipVerify},
	}
	if websocketproxy.DefaultDialer.TLSClientConfig == nil {
		websocketproxy.DefaultDialer.TLSClientConfig = &tls.Config{}
	}
	websocketproxy.DefaultDialer.TLSClientConfig.InsecureSkipVerify = *skipVerify
	websocketproxy.DefaultDialer.EnableCompression = *wsCompress
	websocketproxy.DefaultUpgrader.EnableCompression = *wsCompress
	websocketproxy.DefaultUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
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
	err := oidcSetup()
	if err == nil {
		err = refreshFence()
	}
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
