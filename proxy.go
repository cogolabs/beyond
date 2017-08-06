package main

import (
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/koding/websocketproxy"
)

var (
	hostProxy = map[string]*httputil.ReverseProxy{}
	hostWS    = map[string]*websocketproxy.WebsocketProxy{}
)

func http2ws(u string) string {
	if strings.HasPrefix(u, "http:") {
		return "ws:" + u[5:]
	}
	if strings.HasPrefix(u, "https:") {
		return "wss:" + u[6:]
	}
	return u
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
			w, err := url.Parse(http2ws(x))
			if err != nil {
				return err
			}
			hostProxy[u.Host] = httputil.NewSingleHostReverseProxy(u)
			hostWS[u.Host] = websocketproxy.NewProxy(w)
			hostWS[u.Host].Director = websocketproxyDirector
		}
	}
	return nil
}
