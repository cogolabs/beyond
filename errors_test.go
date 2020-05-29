package beyond

import (
	"net/http"
	"testing"

	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/assert"
)

func init() {
	*errorEmail = "support@mycompany.com"
}

type testIntString struct {
	code int
	text string
}

func TestErrorQuery(t *testing.T) {
	for errorQuery, expectedValues := range map[string]testIntString{
		"invalid_request":           {400, "400 - Bad Request"},
		"access_denied":             {403, "403 - Forbidden"},
		"invalid_resource":          {404, "404 - Not Found"},
		"unknown":                   {500, "500 - Internal Server Error"},
		"server_error":              {500, "500 - Internal Server Error"},
		"unsupported_response_type": {501, "501 - Not Implemented"},
		"temporarily_unavailable":   {503, "503 - Service Unavailable"},
	} {
		testflight.WithServer(h, func(r *testflight.Requester) {
			request, err := http.NewRequest("GET", "/oidc?error="+errorQuery, nil)
			assert.Nil(t, err)
			request.Host = *host
			response := r.Do(request)
			assert.Equal(t, expectedValues.code, response.StatusCode)
			assert.Contains(t, response.Body, expectedValues.text)
		})
	}
}

func TestErrorPlain(t *testing.T) {
	*errorPlain = true

	testflight.WithServer(h, func(r *testflight.Requester) {
		request, err := http.NewRequest("GET", "/oidc?error=server_error&error_description=Foo+Biz", nil)
		assert.Nil(t, err)
		request.Host = *host
		response := r.Do(request)
		assert.Equal(t, 500, response.StatusCode)
		assert.Contains(t, response.Body, "Foo Biz")
	})

	*errorPlain = false
}
