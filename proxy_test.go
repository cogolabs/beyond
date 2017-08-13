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
