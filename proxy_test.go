package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func init() {
	hostProxy.Store("test.com", nil)
}

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

func TestWebsocketEcho(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(nexthop))
	defer server.Close()

	h := http.Header{"Host": []string{"echo.websocket.org"}}
	c, _, err := websocket.DefaultDialer.Dial("ws:"+server.URL[5:], h)
	assert.NotNil(t, c)
	assert.NoError(t, err)
	assert.NoError(t, c.WriteJSON(map[string]string{"test": "123"}))
	v := map[string]string{}
	assert.NoError(t, c.ReadJSON(&v))
	assert.Equal(t, "123", v["test"])
}

func TestWebsocketNew(t *testing.T) {
	r, err := http.NewRequest("GET", "https://socketio.fire.base", nil)
	assert.NoError(t, err)
	p, err := websocketproxyNew(r)
	assert.NoError(t, err)

	assert.Equal(t, "wss:"+r.URL.String()[6:], p.Backend(r).String())
}
