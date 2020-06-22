package beyond

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/drewolson/testflight"

	"github.com/stretchr/testify/assert"
)

func init() {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		switch r.Header.Get("Authorization") {
		case "":
			w.WriteHeader(401)
			fmt.Fprint(w, "{\"errors\":[{\"code\":\"UNAUTHORIZED\",\"detail\":{},")
		default:
			w.WriteHeader(200)
			fmt.Fprint(w, `{"token":"secret1"}`)
		}
	}))
	*dockerBase = testServer.URL
	*dockerScheme = "http"
}

func TestDockerV2(t *testing.T) {
	err := dockerSetup(":")
	assert.Error(t, err)

	testflight.WithServer(h, func(r *testflight.Requester) {
		v2get := r.Get("/v2/auth")
		assert.Equal(t, 418, v2get.StatusCode)

		req, err := http.NewRequest("GET", "/v2/", nil)
		assert.NoError(t, err)
		req.Host = dockerHost
		req.Header.Set("User-Agent", "docker/1.12.6 go/go1.7.4")

		resp := r.Do(req)
		assert.Equal(t, 401, resp.StatusCode)
		assert.Equal(t, "", resp.Body)
		// assert.Equal(t, `Bearer realm="`+*dockerBase+`/v2/auth",service="docker.colofoo.net"`, resp.Header.Get("WWW-Authenticate"))
		assert.True(t, strings.HasPrefix(resp.Header.Get("WWW-Authenticate"), "Bearer realm="))

		req, err = http.NewRequest("GET", "/v2/auth?account=joe&client_id=docker&offline_token=true&service=docker.colofoo.net", nil)
		assert.NoError(t, err)
		req.Host = dockerHost
		req.SetBasicAuth("joe", "secret0")
		req.Header.Set("User-Agent", "docker/1.12.6 go/go1.7.4")

		resp = r.Do(req)
		assert.Equal(t, 200, resp.StatusCode)
		assert.True(t, strings.HasPrefix(resp.Body, "{\"token\":\""))

		v := map[string]interface{}{}
		err = json.Unmarshal([]byte(resp.Body), &v)
		assert.NoError(t, err)
		token := v["token"].(string)
		assert.NotZero(t, token)

		assert.True(t, strings.HasPrefix(token, "MTU5Mjg"))

		req, err = http.NewRequest("GET", "/v2/namespaces", nil)
		assert.NoError(t, err)
		req.Host = dockerHost
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("User-Agent", "docker/1.12.6 go/go1.7.4")

		resp = r.Do(req)
		assert.Equal(t, 200, resp.StatusCode)
		assert.True(t, strings.HasPrefix(resp.Body, "{\"token\":\""))

		token = token[:len(token)/2]

		req, err = http.NewRequest("GET", "/v2/namespaces", nil)
		assert.NoError(t, err)
		req.Host = dockerHost
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("User-Agent", "docker/1.12.6 go/go1.7.4")

		resp = r.Do(req)
		assert.Equal(t, 401, resp.StatusCode)
		assert.Equal(t, "", resp.Body)
	})
}
