package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostMasq(t *testing.T) {
	assert.NoError(t, hostMasqSetup(""))
	assert.Equal(t, "test1.com", hostRewrite("test1.com"))

	assert.NoError(t, hostMasqSetup("test1.com=test1.net,test2.com=test2.org"))
	assert.Equal(t, "test1.net", hostRewrite("test1.com"))
	assert.Equal(t, "test2.org", hostRewrite("test2.com"))
}
