package beyond

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/securecookie"
)

// via https://docs.docker.com/registry/spec/auth/token/

var (
	dockerBase   = flag.String("docker-url", "https://docker.colofoo.net", "")
	dockerScheme = flag.String("docker-auth-scheme", "https", "(only for testing)")

	dockerHost string
	dockerRP   *httputil.ReverseProxy
)

func dockerSetup(u string) error {
	dockerURL, err := url.Parse(u)
	if err != nil {
		return err
	}

	dockerHost = dockerURL.Hostname()
	dockerRP = httputil.NewSingleHostReverseProxy(dockerURL)
	dockerRP.ModifyResponse = dockerModifyResponse
	return nil
}

func dockerModifyResponse(resp *http.Response) error {
	logRoundtrip(resp)

	if resp.Header.Get("WWW-Authenticate") != "" {
		resp.Header.Set("WWW-Authenticate", `Bearer realm="`+*dockerScheme+`://`+resp.Request.Host+`/v2/auth",service="`+dockerHost+`"`)
	}
	if resp.Request.URL.Path != "/v2/auth" || resp.StatusCode != 200 {
		return nil
	}

	// > GET /v2/auth?account=joe&client_id=docker&offline_token=true&service=docker.colofoo.net
	// < HTTP/1.1 200 OK
	// < {"token": "opaqueXYZ"}

	v := map[string]interface{}{}
	err := json.NewDecoder(resp.Body).Decode(&v)
	if err == nil {
		token, ok := v["token"].(string)
		if ok && strings.Contains(token, ".") {
			claim64 := strings.Split(token, ".")[1]
			data, err := base64.RawStdEncoding.DecodeString(claim64)
			if err == nil {
				claim := new(dockerClaimSet)
				err = json.Unmarshal(data, claim)
				if err == nil && claim.Context.Kind == "user" {
					v["token"], err = securecookie.EncodeMulti("token", v["token"], store.Codecs...)
					if err == nil {
						var buf bytes.Buffer
						err = json.NewEncoder(&buf).Encode(v)
						if err == nil {
							// < {"token": "beyondXYZ"}

							resp.Body = ioutil.NopCloser(&buf)
							resp.ContentLength = int64(buf.Len())
							resp.Header.Set("Content-Length", strconv.Itoa(buf.Len()))
							return nil
						}
					}
				}
			}
		}
	}

	return err
}

func dockerHandler(w http.ResponseWriter, r *http.Request) {
	allow := r.URL.Path == "/v2/auth" && len(r.Header.Get("Authorization")) > 0
	if !allow {
		token := strings.Split(r.Header.Get("Authorization"), " ")
		if len(token) > 1 && token[0] == "Bearer" {
			bearer := token[1]
			err := securecookie.DecodeMulti("token", bearer, &bearer, store.Codecs...)
			if err == nil {
				allow = true
				r.Header.Set("Authorization", "Bearer "+bearer)
			}
		}
	}
	if allow {
		dockerRP.ServeHTTP(w, r)
		return
	}

	w.Header().Set("Docker-Distribution-Api-Version", "registry/2.0")
	w.Header().Set("WWW-Authenticate", `Bearer realm="`+*dockerScheme+`://`+r.Host+`/v2/auth",service="`+dockerHost+`"`)
	w.WriteHeader(401)
}

// https://docs.docker.com/registry/spec/auth/jwt/
//
// {
// 	"context": {
// 	  "com.apostille.root": "$disabled"
// 	},
// 	"aud": "docker.colofoo.net",
// 	"exp": 1593910505,
// 	"iss": "quay",
// 	"iat": 1593906905,
// 	"nbf": 1593906905,
// 	"sub": "(anonymous)"
// }
//
// {
// 	"access": [
// 	  {
// 		"type": "repository",
// 		"name": "cogolabs/beyond",
// 		"actions": [
// 		  "pull"
// 		]
// 	  }
// 	],
// 	"context": {
// 	  "entity_kind": "appspecifictoken",
// 	  "kind": "user",
// 	  "version": 2,
// 	  "com.apostille.root": "$disabled",
// 	  "user": "joe",
// 	  "entity_reference": "4ac6f0e7-7bd2-4aea-9a77-738e1b98f22f"
// 	},
// 	"aud": null,
// 	"exp": 1593911101,
// 	"iss": "quay",
// 	"iat": 1593907501,
// 	"nbf": 1593907501,
// 	"sub": "joe"
// }
type dockerClaimSet struct {
	Access []struct {
		Type    string   `json:"type"`
		Name    string   `json:"name"`
		Actions []string `json:"actions"`
	} `json:"access"`
	Context struct {
		EntityKind       string `json:"entity_kind"`
		Kind             string `json:"kind"`
		Version          int    `json:"version"`
		ComApostilleRoot string `json:"com.apostille.root"`
		User             string `json:"user"`
		EntityReference  string `json:"entity_reference"`
	} `json:"context"`
	Aud string `json:"aud"`
	Exp int    `json:"exp"`
	Iss string `json:"iss"`
	Iat int    `json:"iat"`
	Nbf int    `json:"nbf"`
	Sub string `json:"sub"`
}
