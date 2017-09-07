package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestH2W(t *testing.T) {
	req, err := http.NewRequest("GET", "/foo?key=bar", nil)
	assert.NoError(t, err)
	req.Host = "websocketserver:9443"
	actual, err := http2ws(req)
	assert.NoError(t, err)
	assert.Equal(t, "wss://websocketserver:9443/foo?key=bar", actual.String())

	req.Host = "websock etserver:9443"
	actual, err = http2ws(req)
	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestWebsocketNew(t *testing.T) {
	r, err := http.NewRequest("GET", "https://socketio.fire.base", nil)
	assert.NoError(t, err)
	p, err := websocketproxyNew(r)
	assert.NoError(t, err)

	assert.Equal(t, "wss:"+r.URL.String()[6:], p.Backend(r).String())
}
