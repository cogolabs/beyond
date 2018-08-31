package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLearnHostScheme(t *testing.T) {
	assert.Equal(t, "https://golang.org", learnHostScheme("golang.org"))
}
