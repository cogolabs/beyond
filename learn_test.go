package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLearnHostScheme(t *testing.T) {
	ports := *learnHTTPSPorts
	*learnHTTPSPorts = ""

	assert.Equal(t, "http://neverssl.com", learnHostScheme("neverssl.com"))

	*learnHTTPSPorts = ports
	assert.Equal(t, "https://golang.org", learnHostScheme("golang.org"))
	assert.NotNil(t, learn("golang.org"))
}
