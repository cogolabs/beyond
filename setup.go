package beyond

import (
	"crypto/tls"
	"flag"
	"net/http"

	"github.com/dghubble/sessions"
	"github.com/koding/websocketproxy"
)

var (
	host = flag.String("host", "beyond.colofoo.net", "hostname of self, eg. when generating OAuth redirect URLs")

	healthPath  = flag.String("health-path", "/healthz/ping", "URL of the health endpoint")
	healthReply = flag.String("health-reply", "ok", "response body of the health endpoint")

	cookieAge  = flag.Int("cookie-age", 3600*6, "MaxAge setting in seconds")
	cookieDom  = flag.String("cookie-domain", ".colofoo.net", "session cookie domain")
	cookieKey1 = flag.String("cookie-key1", "t8yG1gmeEyeb7pQpw544UeCTyDfPkE6u", "key1 of cookie crypto pair")
	cookieKey2 = flag.String("cookie-key2", "Q599vrruZRhLFC144thCRZpyHM7qGDjt", "key2 of cookie crypto pair")
	cookieName = flag.String("cookie-name", "beyond", "session cookie name")

	fouroFourMessage = flag.String("404-message", "Please contact your network administrators to whitelist this system.", "message to use for unlisted hosts when learning is disabled or fails")
	fouroOneCode     = flag.Int("401-code", 418, "status to respond when a user needs authentication")
	headerPrefix     = flag.String("header-prefix", "Beyond", "prefix extra headers with this string")

	skipVerify = flag.Bool("insecure-skip-verify", false, "allow TLS backends without valid certificates")
	wsCompress = flag.Bool("websocket-compression", false, "allow websocket transport compression (gorilla/experimental)")

	store *sessions.CookieStore
)

func Setup() error {
	// setup encrypted cookies
	store = sessions.NewCookieStore([]byte(*cookieKey1), []byte(*cookieKey2))
	store.Config.Domain = *cookieDom
	store.Config.MaxAge = *cookieAge

	// setup backend encryption
	http.DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: *skipVerify,
		},
	}

	// setup websockets
	if websocketproxy.DefaultDialer.TLSClientConfig == nil {
		websocketproxy.DefaultDialer.TLSClientConfig = &tls.Config{}
	}
	websocketproxy.DefaultDialer.TLSClientConfig.InsecureSkipVerify = *skipVerify
	websocketproxy.DefaultDialer.EnableCompression = *wsCompress
	websocketproxy.DefaultUpgrader.EnableCompression = *wsCompress
	websocketproxy.DefaultUpgrader.CheckOrigin = websocketproxyCheckOrigin

	err := dockerSetup(*dockerBase)
	if err == nil {
		err = hostMasqSetup(*hostMasq)
	}
	if err == nil {
		err = logSetup()
	}
	if err == nil {
		err = oidcSetup(*oidcIssuer)
	}
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
