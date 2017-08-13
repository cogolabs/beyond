package main

import (
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

func newWebSocket(r *http.Request) (*websocketproxy.WebsocketProxy, error) {
	ws, err := http2ws(r)
	p := websocketproxy.NewProxy(ws)
	p.Director = websocketproxyDirector
	return p, err
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
