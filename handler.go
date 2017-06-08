package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
)

func consume(w http.ResponseWriter, r *http.Request) {
	var (
		firstname = r.FormValue("firstname")
		lastname  = r.FormValue("lastname")
		email     = r.FormValue("email")
		timestamp = r.FormValue("timestamp")
		signature = r.FormValue("signature")
	)
	key := sha1.New()
	key.Write([]byte(firstname + lastname + email + timestamp + *token))
	keyx := hex.EncodeToString(key.Sum(nil))

	if keyx != signature {
		w.WriteHeader(403)
		fmt.Fprint(w, "invalid signature")
		return
	}

	session, err := store.Get(r, "beyond")
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, err.Error())
	}
	session.Values["email"] = email
	next, _ := session.Values["next"].(string)
	session.Values["next"] = ""
	session.Save(w)

	http.Redirect(w, r, next, 302)
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/ping":
		fmt.Fprint(w, "OK")
		return
	case "/onelogin/consume":
		consume(w, r)
		return
	}

	session, err := store.Get(r, "beyond")
	if err != nil {
		session = store.New("beyond")
	}
	email, _ := session.Values["email"].(string)
	next, _ := session.Values["next"].(string)
	proxy := hostProxy[r.Host]
	if proxy == nil {
		w.WriteHeader(404)
		fmt.Fprint(w, "unknown URL")
		return
	}

	if email != "" {
		proxy.ServeHTTP(w, r)
	}

	if next == "" {
		session.Values["next"] = r.URL.String()
		session.Save(w)
	}
	http.Redirect(w, r, fmt.Sprintf("https://app.onelogin.com/launch/%d", *appid), 302)
}
