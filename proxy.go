package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/koding/websocketproxy"
)

var (
	hostProxy = map[string]*httputil.ReverseProxy{}
)

func http2ws(r *http.Request) (*url.URL, error) {
	target := "wss://" + r.Host + r.URL.RequestURI()
	return url.Parse(target)
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

func reproxy() error {
	sites.RLock()
	defer sites.RUnlock()
	for _, v := range sites.m {
		for x := range v {
			u, err := url.Parse(x)
			if err != nil {
				return err
			}
			hostProxy[u.Host] = httputil.NewSingleHostReverseProxy(u)
		}
	}
	return nil
}

func websocketproxyDirector(incoming *http.Request, out http.Header) {
	out.Set("User-Agent", incoming.UserAgent())
	out.Set("X-Forwarded-Proto", "https")
}

func websocketproxyNew(r *http.Request) (*websocketproxy.WebsocketProxy, error) {
	ws, err := http2ws(r)
	p := websocketproxy.NewProxy(ws)
	p.Director = websocketproxyDirector
	return p, err
}
