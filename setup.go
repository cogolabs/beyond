package beyond

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/dghubble/sessions"
	"github.com/koding/websocketproxy"
)

var (
	host = flag.String("beyond-host", "beyond.myorg.net", "hostname of self")

	healthPath  = flag.String("health-path", "/healthz/ping", "URL of the health endpoint")
	healthReply = flag.String("health-reply", "ok", "response body of the health endpoint")

	cookieAge  = flag.Int("cookie-age", 3600*6, "MaxAge setting in seconds")
	cookieDom  = flag.String("cookie-domain", ".myorg.net", "session cookie domain")
	cookieKey1 = flag.String("cookie-key1", "", `key1 of cookie crypto pair (example: "t8yG1gmeEyeb7pQpw544UeCTyDfPkE6u")`)
	cookieKey2 = flag.String("cookie-key2", "", `key2 of cookie crypto pair (example: "Q599vrruZRhLFC144thCRZpyHM7qGDjt")`)
	cookieName = flag.String("cookie-name", "beyond", "session cookie name")

	fouroFourMessage = flag.String("404-message", "Please contact the application administrators to setup access.", "message to use when backend apps do not respond")
	fouroOneCode     = flag.Int("401-code", 418, "status to respond when a user needs authentication")
	headerPrefix     = flag.String("header-prefix", "Beyond", "prefix extra headers with this string")

	skipVerify = flag.Bool("insecure-skip-verify", false, "allow TLS backends without valid certificates")
	wsCompress = flag.Bool("websocket-compression", false, "allow websocket transport compression (gorilla/experimental)")

	store *sessions.CookieStore
)

// Setup initializes all configured modules
func Setup() error {
	if len(*cookieKey1) == 0 {
		return fmt.Errorf("missing cookie key")
	}

	// setup encrypted cookies
	store = sessions.NewCookieStore([]byte(*cookieKey1), []byte(*cookieKey2))
	store.Config.Domain = *cookieDom
	store.Config.MaxAge = *cookieAge
	store.Config.HTTPOnly = true
	store.Config.SameSite = http.SameSiteNoneMode
	store.Config.Secure = true

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

	dURLs := []string{*dockerBase}
	if len(*dockerURLs) > 0 {
		dURLs = append(dURLs, strings.Split(*dockerURLs, ",")...)
	}

	err := dockerSetup(dURLs...)
	if err == nil {
		err = federateSetup()
	}
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
		err = samlSetup()
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
