package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

func beyond(w http.ResponseWriter, r *http.Request) {
	setCacheControl(w)
	switch r.URL.Path {

	case "/launch":
		session, err := store.Get(r, *cookieName)
		if err != nil {
			session = store.New(*cookieName)
		}
		session.Values["next"] = r.FormValue("next")
		state, _ := randhex32()
		session.Values["state"] = state
		session.Save(w)

		next := oidcConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
		jsRedirect(w, next)

	case "/oidc":
		session, err := store.Get(r, *cookieName)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintln(w, err.Error())
			return
		}
		if state, ok := session.Values["state"].(string); !ok || state != r.FormValue("state") {
			w.WriteHeader(403)
			fmt.Fprintln(w, "invalid state")
			return
		}
		email, err := oidcVerify(r.FormValue("code"))
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintln(w, err.Error())
			return
		}
		session.Values["email"] = email
		next, _ := session.Values["next"].(string)
		session.Values["next"] = ""
		session.Values["state"] = ""
		session.Save(w)

		http.Redirect(w, r, next, 302)

	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/ping" {
		fmt.Fprintln(w, "OK")
		return
	}
	if r.Host == *host {
		beyond(w, r)
		return
	}
	if whitelisted(r) {
		nexthop(w, r)
		return
	}

	session, err := store.Get(r, *cookieName)
	if err != nil {
		session = store.New(*cookieName)
	}
	email, _ := session.Values["email"].(string)
	switch email {
	case "":
		login(w, r)
	default:
		nexthop(w, r)
	}
}

func nexthop(w http.ResponseWriter, r *http.Request) {
	var proxy http.Handler
	proxy, ok := hostProxy[r.Host]
	if !ok {
		setCacheControl(w)
		w.WriteHeader(404)
		fmt.Fprintln(w, "invalid URL")
		return
	}

	if r.Header.Get("Upgrade") == "websocket" {
		proxy, _ = websocketproxyNew(r)
	}
	proxy.ServeHTTP(w, r)
}

func login(w http.ResponseWriter, r *http.Request) {
	setCacheControl(w)
	w.WriteHeader(401)

	// short-circuit WS+AJAX
	if r.Header.Get("Upgrade") != "" || r.Header.Get("X-Requested-With") != "" {
		return
	}

	jsRedirect(w, "https://"+*host+"/launch?next="+url.QueryEscape("https://"+r.Host+r.RequestURI))
}

func jsRedirect(w http.ResponseWriter, next string) {
	// hack to guarantee interactive session
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
<script type="text/javascript">
window.location.replace("%s");
</script>
`, next)
}

func setCacheControl(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
}

func randhex32() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return fmt.Sprintf("%x", b), err
}
