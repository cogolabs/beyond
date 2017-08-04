package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestACLRefresh(t *testing.T) {
	assert.NoError(t, refreshFence())
	assert.NoError(t, refreshSites())
	assert.NoError(t, refreshWhitelist())

	assert.NotEmpty(t, fence.m)
	assert.NotEmpty(t, sites.m["git"])
	assert.NotEmpty(t, whitelist.m["host"])
	assert.NotEmpty(t, whitelist.m["path"])
}
