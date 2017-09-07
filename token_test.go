package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/assert"
)

var (
	tokenTestLogin = "user1"
	tokenTestToken = "932928c0a4edf9878ee0257a1d8f4d06adaaffee"

	tokenServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("access_token") == "invalid" {
			io.WriteString(w, "{")
			return
		}
		if r.FormValue("access_token") != tokenTestToken {
			w.WriteHeader(403)
			return
		}
		json.NewEncoder(w).Encode(tokenUser{Login: tokenTestLogin})
	}))
)

func TestTokenError(t *testing.T) {
	*tokenBase = "https://foo.bar?"

	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	assert.Equal(t, "", tokenAuth(r))
	r.Header.Set("Authorization", "token test")
	assert.Equal(t, "", tokenAuth(r))

	*tokenBase = tokenServer.URL + "/?access_token="
	r.Header.Set("Authorization", "token invalid")
	assert.Equal(t, "", tokenAuth(r))
}

func TestTokenFederation(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	r.Header.Set("Authorization", "token test")
	assert.Equal(t, "", tokenAuth(r))

	r.Header.Set("Authorization", "token "+tokenTestToken)
	login1 := tokenAuth(r)
	login2 := tokenAuth(r)
	assert.Equal(t, "user1", login1)
	assert.Equal(t, "user1", login2)
}

func TestTokenSuccess(t *testing.T) {
	testflight.WithServer(h, func(r *testflight.Requester) {
		request, err := http.NewRequest("GET", "/ip", nil)
		assert.Nil(t, err)
		request.Header.Set("Authorization", "token "+tokenTestToken)
		request.Host = "httpbin.org"
		response := r.Do(request)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "{\n  \"origin\"", strings.Split(response.Body, ":")[0])
	})
}
