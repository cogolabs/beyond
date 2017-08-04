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
		request.Host = "alachart.colofoo.net"
		assert.Nil(t, err)
		response := r.Do(request)
		assert.Equal(t, 401, response.StatusCode)
		assert.Equal(t, "", response.Header.Get("Set-Cookie"))
		assert.Equal(t, "\n<script type=\"text/javascript\">\n  window.location.replace(\"https://beyond.colofoo.net/launch?next=https%3A%2F%2Falachart.colofoo.net%2Ftest%3Fa%3D1\");\n</script>\n  ", response.Body)
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
		// assert.Equal(t, "\n<script type=\"text/javascript\">\n  window.location.replace(\"https://app.onelogin.com/launch/"+strconv.FormatInt(*appid, 10)+"\");\n</script>\n  ", response.Body)
	})
}

func TestHandlerWhitelist(t *testing.T) {
	testflight.WithServer(h, func(r *testflight.Requester) {
		request, err := http.NewRequest("GET", "/", nil)
		request.Host = "assets.git.colofoo.net"
		assert.Nil(t, err)
		response := r.Do(request)
		assert.Equal(t, 200, response.StatusCode)
		assert.NotEqual(t, "", response.Header.Get("Set-Cookie"))
		assert.Contains(t, response.Body, "Recent Repositories")
	})
	testflight.WithServer(h, func(r *testflight.Requester) {
		request, err := http.NewRequest("GET", "/.well-known/acme-challenge/test", nil)
		request.Host = "git.colofoo.net"
		assert.Nil(t, err)
		response := r.Do(request)
		assert.Equal(t, 404, response.StatusCode)
		assert.NotEqual(t, "", response.Header.Get("Set-Cookie"))
		assert.Contains(t, response.Body, "Page not found")
	})
}
