package main

import (
	"net/http"
	"testing"

	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/assert"
)

var h = http.HandlerFunc(handler)

func TestHandlerPing(t *testing.T) {
	testflight.WithServer(h, func(r *testflight.Requester) {
		request, err := http.NewRequest("GET", "/ping", nil)
		assert.Nil(t, err)
		response := r.Do(request)
		assert.Equal(t, 200, response.StatusCode)
	})
}

func TestHandlerGo(t *testing.T) {
	testflight.WithServer(h, func(r *testflight.Requester) {
		request, err := http.NewRequest("GET", "/test?a=1", nil)
		request.Host = "github.com"
		assert.Nil(t, err)
		response := r.Do(request)
		assert.Equal(t, 401, response.StatusCode)
		assert.Equal(t, "", response.Header.Get("Set-Cookie"))
		assert.Equal(t, "\n<script type=\"text/javascript\">\n  window.location.replace(\"https://beyond.colofoo.net/launch?next=https%3A%2F%2Fgithub.com%2Ftest%3Fa%3D1\");\n</script>\n  ", response.Body)
	})
}

func TestHandlerLaunch(t *testing.T) {
	testflight.WithServer(h, func(r *testflight.Requester) {
		request, err := http.NewRequest("GET", "/launch?next=https%3A%2F%2Falachart.colofoo.net%2Ftest%3Fa%3D1", nil)
		request.Host = "beyond.colofoo.net"
		assert.Nil(t, err)
		response := r.Do(request)
		assert.Equal(t, 200, response.StatusCode)
		assert.NotEqual(t, "", response.Header.Get("Set-Cookie"))
	})
}

func TestHandlerWhitelist(t *testing.T) {
	testflight.WithServer(h, func(r *testflight.Requester) {
		request, err := http.NewRequest("GET", "/favicon.ico", nil)
		request.Host = "assets.github.com"
		assert.Nil(t, err)
		response := r.Do(request)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "", response.Header.Get("Set-Cookie"))
		assert.Equal(t, []byte{00, 00, 01, 00, 2, 0, 0x10, 0x10, 0, 0}, response.RawBody[:10])
	})
	testflight.WithServer(h, func(r *testflight.Requester) {
		request, err := http.NewRequest("GET", "/.well-known/acme-challenge/test", nil)
		request.Host = "github.com"
		assert.Nil(t, err)
		response := r.Do(request)
		assert.Equal(t, 404, response.StatusCode)
		assert.NotEqual(t, "", response.Header.Get("Set-Cookie"))
		assert.Contains(t, response.Body, "Page not found")
	})
}
