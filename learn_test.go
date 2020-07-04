package beyond

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLearnHostScheme(t *testing.T) {
	ports1 := *learnHTTPSPorts
	ports2 := *learnHTTPPorts

	learnTest1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	learnTest2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	*learnHTTPSPorts = strings.Split(learnTest1.URL, ":")[2]
	*learnHTTPPorts = strings.Split(learnTest2.URL, ":")[2]
	base1 := learnBase("127.0.0.1")
	learnTest1.Close()
	base2 := learnBase("127.0.0.1")
	learnTest2.Close()
	assert.Equal(t, strings.Replace(learnTest1.URL, "http://", "https://", 1), base1)
	assert.Equal(t, learnTest2.URL, base2)

	*learnHTTPSPorts = ""
	*learnHTTPPorts = "80"
	assert.Equal(t, "http://neverssl.com", learnBase("neverssl.com"))

	*learnHTTPSPorts = ports1
	*learnHTTPPorts = ports2
	assert.Equal(t, "https://golang.org", learnBase("golang.org"))
	assert.NotNil(t, learn("golang.org"))

	assert.Equal(t, "https://golang.org:443", learnBase("golang.org:443"))
}
