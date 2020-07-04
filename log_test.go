package beyond

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	*logJSON = true
	*logXFF = true
}

func TestLogHTTP(t *testing.T) {
	*logHTTP = true

	req, err := http.NewRequest("GET", "/log", nil)
	assert.NoError(t, err)
	resp := &http.Response{Request: req}
	logRoundtrip(resp)
}
