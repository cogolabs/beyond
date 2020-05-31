package beyond

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	oidc "github.com/coreos/go-oidc"
	"github.com/drewolson/testflight"
	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var (
	oidcServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			err := json.NewEncoder(w).Encode(oidcWK)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			return

		case "/token":
			fmt.Fprint(w, `{"typ":"JWT","alg":"HS256"}`)
			return

		case "/next":
			fmt.Fprint(w, "NEXT")
			return

		default:
			WithField("path", r.URL.Path).Error("Invalid OIDC Request")
			return

		}
	}))
	oidcWK = struct {
		Issuer                 string
		Authorization_endpoint string
		Token_endpoint         string
	}{
		"/issuer",
		"/authorize",
		"/token",
	}
)

type oidcMock struct{}

func (o *oidcMock) AuthCodeURL(state string, opt ...oauth2.AuthCodeOption) string {
	return oidcServer.URL + "/AuthCodeURL"
}

func (o *oidcMock) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token := &oauth2.Token{}
	token.AccessToken = "AccessToken"
	token.Expiry = time.Now().Add(time.Hour)
	token = token.WithExtra(map[string]interface{}{"id_token": "IDToken"})
	return token, nil
}

func (o *oidcMock) Verify(ctx context.Context, raw string) (*oidc.IDToken, error) {
	if raw == "err" {
		return nil, http.ErrHijacked
	}
	token := &oidc.IDToken{}

	// requires gcflags=-l
	monkey.PatchInstanceMethod(reflect.TypeOf(token), "Claims", func(t *oidc.IDToken, v interface{}) error {
		v.(*oidcClaims).Email = "user3@domain3.com"
		return nil
	})
	return token, nil
}

func init() {
	// *oidcIssuer = oidcServer.URL
	oidcWK.Issuer = oidcServer.URL + oidcWK.Issuer
	oidcWK.Token_endpoint = oidcServer.URL + oidcWK.Token_endpoint
	oidcWK.Authorization_endpoint = oidcServer.URL + oidcWK.Authorization_endpoint
}

func TestOIDCSetup(t *testing.T) {
	assert.Contains(t, oidcSetup("ftp://localhost").Error(), "unsupported protocol scheme")
}

func TestOIDCSuccess(t *testing.T) {
	mock := &oidcMock{}
	oidcConfig = mock
	oidcVerifier = mock

	testflight.WithServer(h, func(r *testflight.Requester) {
		request, err := http.NewRequest("GET", "/oidc?state=barbaz&next=localhost/next", nil)
		assert.Nil(t, err)

		vals := map[string]interface{}{"state": "barbaz", "next": oidcServer.URL + "/next"}
		cookieValue, err := securecookie.EncodeMulti(*cookieName, &vals, store.Codecs...)
		assert.NoError(t, err)
		request.AddCookie(&http.Cookie{Name: *cookieName, Value: cookieValue})

		request.Host = *host
		response := r.Do(request)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "NEXT", response.Body)
	})
}

func TestOIDCVerifyToken(t *testing.T) {
	token := &oauth2.Token{}
	s, err := oidcVerifyToken(context.TODO(), token)
	assert.Empty(t, s)
	assert.Equal(t, "missing ID token", err.Error())
}

func TestOIDCVerifyTokenID(t *testing.T) {
	email, err := oidcVerifyTokenID(context.TODO(), "err")
	assert.Equal(t, "", email)
	assert.NoError(t, err)
}
