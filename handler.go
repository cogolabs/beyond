package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

func beyond(w http.ResponseWriter, r *http.Request) {
	setCacheControl(w)
	switch r.URL.Path {

	case "/launch":
		session, err := store.Get(r, "beyond")
		if err != nil {
			session = store.New("beyond")
		}
		session.Values["next"] = r.FormValue("next")
		state := randhex32()
		session.Values["state"] = state
		session.Save(w)

		next := oidcConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
		fmt.Fprintf(w, `
<script type="text/javascript">
  window.location.replace("%s");
</script>
  `, next)

	case "/oidc":
		session, err := store.Get(r, "beyond")
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
			fmt.Fprintf(w, err.Error())
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

	session, err := store.Get(r, "beyond")
	if err != nil {
		session = store.New("beyond")
	}
	email, _ := session.Values["email"].(string)
	proxy := hostProxy[r.Host]

	// unconfigured
	if proxy == nil {
		setCacheControl(w)
		w.WriteHeader(404)
		fmt.Fprintln(w, "unknown URL")
		return
	}

	// allow
	if email != "" || whitelisted(r.Host, r.URL.Path) {
		if r.Header.Get("Upgrade") == "websocket" {
			hostWS[r.Host].ServeHTTP(w, r)
		} else {
			proxy.ServeHTTP(w, r)
		}
		return
	}

	// deny
	setCacheControl(w)

	// short-circuit WS+AJAX
	if r.Header.Get("Upgrade") != "" || r.Header.Get("X-Requested-With") != "" {
		w.WriteHeader(401)
		return
	}

	// interstitial landing to guarantee interactive before cookie save
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(401)
	fmt.Fprintf(w, `
<script type="text/javascript">
  window.location.replace("https://%s/launch?next=%s");
</script>
  `, *host, url.QueryEscape("https://"+r.Host+r.RequestURI))
}

func setCacheControl(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
}

func randhex32() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%x", b)
}
