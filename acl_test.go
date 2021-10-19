package beyond

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
	*allowlistURL = ""

	assert.NoError(t, refreshFence())
	assert.NoError(t, refreshSites())
	assert.NoError(t, refreshAllowlist())

	*fenceURL = aclErrorBase
	*sitesURL = aclErrorBase
	*allowlistURL = aclErrorBase

	assert.Contains(t, refreshFence().Error(), "connection refused")
	assert.Contains(t, refreshSites().Error(), "connection refused")
	assert.Contains(t, refreshAllowlist().Error(), "connection refused")

	cwd, _ := os.Getwd()
	*fenceURL = "file://" + cwd + "/example/error.json"
	*sitesURL = "file://" + cwd + "/example/error.json"
	*allowlistURL = "file://" + cwd + "/example/error.json"
	assert.EqualError(t, refreshFence(), "unexpected EOF")
	assert.EqualError(t, refreshSites(), "unexpected EOF")
	assert.EqualError(t, refreshAllowlist(), "unexpected EOF")

	*fenceURL = "file://" + cwd + "/example/fence.json"
	*sitesURL = "file://" + cwd + "/example/sites.json"
	*allowlistURL = "file://" + cwd + "/example/allowlist.json"
	assert.NoError(t, Setup())

	assert.NotEmpty(t, fence.m)
	assert.NotEmpty(t, sites.m["git"])
	assert.NotEmpty(t, allowlist.m["host"])
	assert.NotEmpty(t, allowlist.m["path"])

	reqDeny, _ := http.NewRequest("GET", "https://deny", nil)
	assert.True(t, deny(reqDeny, "consultant@gmail.com"))
	reqAllow, _ := http.NewRequest("GET", "https://github.com/test", nil)
	assert.False(t, deny(reqAllow, "consultant@gmail.com"))
}
