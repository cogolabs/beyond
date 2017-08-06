package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	cwd, _ := os.Getwd()
	*fenceURL = "file://" + cwd + "/example/fence.json"
	*sitesURL = "file://" + cwd + "/example/sites.json"
	*whitelistURL = "file://" + cwd + "/example/whitelist.json"

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	httpACL.Transport = t
}

func TestACLSetup(t *testing.T) {
	assert.NoError(t, setup())

	assert.NotEmpty(t, fence.m)
	assert.NotEmpty(t, sites.m["git"])
	assert.NotEmpty(t, whitelist.m["host"])
	assert.NotEmpty(t, whitelist.m["path"])
}
