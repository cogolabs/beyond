package beyond

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	*cookieKey1 = "t8yG1gmeEyeb7pQpw544UeCTyDfPkE6u"
	*cookieKey2 = "Q599vrruZRhLFC144thCRZpyHM7qGDjt"
}

func TestSetupErr(t *testing.T) {
	prev := *cookieKey1
	*cookieKey1 = ""
	err := Setup()
	assert.Error(t, err)
	*cookieKey1 = prev
}
