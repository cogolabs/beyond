package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/koding/websocketproxy"
)

var (
	hostProxy = map[string]*httputil.ReverseProxy{}
)

func http2ws(r *http.Request) *url.URL {
	target := "wss://" + r.Host + r.URL.String()
	next, err := url.Parse(target)
	if err != nil {
		log.Printf("%s, parsing: %s", err, target)
	}
	return next
}

func newWebSocket(r *http.Request) *websocketproxy.WebsocketProxy {
	p := websocketproxy.NewProxy(http2ws(r))
	p.Director = websocketproxyDirector
	return p
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
