package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	httpACL.Transport = t
}

const (
	aclErrorBase = "http://localhost:9999"
)

func TestACL(t *testing.T) {
	*fenceURL = ""
	*sitesURL = ""
	*whitelistURL = ""

	assert.NoError(t, refreshFence())
	assert.NoError(t, refreshSites())
	assert.NoError(t, refreshWhitelist())

	*fenceURL = aclErrorBase
	*sitesURL = aclErrorBase
	*whitelistURL = aclErrorBase

	assert.Contains(t, refreshFence().Error(), "connection refused")
	assert.Contains(t, refreshSites().Error(), "connection refused")
	assert.Contains(t, refreshWhitelist().Error(), "connection refused")

	cwd, _ := os.Getwd()
	*fenceURL = "file://" + cwd + "/example/error.json"
	*sitesURL = "file://" + cwd + "/example/error.json"
	*whitelistURL = "file://" + cwd + "/example/error.json"
	assert.EqualError(t, refreshFence(), "unexpected EOF")
	assert.EqualError(t, refreshSites(), "unexpected EOF")
	assert.EqualError(t, refreshWhitelist(), "unexpected EOF")

	*fenceURL = "file://" + cwd + "/example/fence.json"
	*sitesURL = "file://" + cwd + "/example/sites.json"
	*whitelistURL = "file://" + cwd + "/example/whitelist.json"
	assert.NoError(t, setup())

	assert.NotEmpty(t, fence.m)
	assert.NotEmpty(t, sites.m["git"])
	assert.NotEmpty(t, whitelist.m["host"])
	assert.NotEmpty(t, whitelist.m["path"])
}
