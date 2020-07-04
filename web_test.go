package beyond

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/assert"
)

var (
	echoServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		switch string(b) {
		case "ping":
			fmt.Fprint(w, "pong")
		default:
			fmt.Fprint(w, string(b))
		}
	}))
)

func TestWebPOST(t *testing.T) {
	testflight.WithServer(h, func(rq *testflight.Requester) {
		request, err := http.NewRequest("POST", "/", strings.NewReader("ping"))
		request.Host = echoServer.URL[7:] // strip the http://
		assert.NoError(t, err)
		request.SetBasicAuth("", tokenTestUserTokens["user1"])
		response := rq.Do(request)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "pong", response.Body)

		request, err = http.NewRequest("POST", "/", strings.NewReader("aliens"))
		request.Host = echoServer.URL // including http://
		assert.NoError(t, err)
		request.SetBasicAuth("", tokenTestUserTokens["user1"])
		response = rq.Do(request)
		assert.Equal(t, 502, response.StatusCode)
		assert.Equal(t, "", response.Body)

		request, err = http.NewRequest("POST", "/", strings.NewReader("aliens"))
		request.Host = echoServer.URL[7:] // strip the http://
		assert.NoError(t, err)
		response = rq.Do(request)
		assert.Equal(t, *fouroOneCode, response.StatusCode)
	})
}
