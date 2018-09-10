package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/koding/websocketproxy"
)

var (
	hostProxy = sync.Map{}
)

func http2ws(r *http.Request) (*url.URL, error) {
	target := "wss://" + hostRewrite(r.Host) + r.URL.RequestURI()
	return url.Parse(target)
}

func nexthop(w http.ResponseWriter, r *http.Request) {
	var (
		nextHost  = hostRewrite(r.Host)
		nextProxy http.Handler
	)

	v, ok := hostProxy.Load(nextHost)
	if ok {
		nextProxy, ok = v.(*httputil.ReverseProxy)
	}
	if !ok && *learnNexthops {
		nextProxy = learn(nextHost)
		if nextProxy != nil {
			hostProxy.Store(nextHost, nextProxy)
			ok = true
		}
	}

	if !ok || nextProxy == nil {
		// unconfigured
		setCacheControl(w)
		w.WriteHeader(404)
		fmt.Fprintln(w, *fouroFourMessage)
		return
	}

	if r.Header.Get("Upgrade") == "websocket" {
		nextProxy, _ = websocketproxyNew(r)
	}
	nextProxy.ServeHTTP(w, r)
}

func reproxy() error {
	cleanup := map[string]bool{}
	hostProxy.Range(func(key interface{}, value interface{}) bool {
		if key, ok := key.(string); ok {
			cleanup[key] = true
		}
		return true
	})
	var lerr error
	sites.RLock()
	for _, v := range sites.m {
		for x := range v {
			u, err := url.Parse(x)
			if err != nil {
				lerr = err
			} else {
				delete(cleanup, u.Host)
				hostProxy.Store(u.Host, httputil.NewSingleHostReverseProxy(u))
			}
		}
	}
	sites.RUnlock()
	for key := range cleanup {
		hostProxy.Delete(key)
	}
	return lerr
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
