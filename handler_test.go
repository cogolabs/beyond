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

func TestHandlerTest(t *testing.T) {
	t.SkipNow()
	testflight.WithServer(h, func(r *testflight.Requester) {
		request, err := http.NewRequest("GET", "/test", nil)
		request.Host = "alachart.colofoo.net"
		assert.Nil(t, err)
		response := r.Do(request)
		assert.Equal(t, 403, response.StatusCode)
	})
}
