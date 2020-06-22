package beyond

import (
	"bytes"
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
	if err != nil {
		return err
	}
	v["token"], err = securecookie.EncodeMulti("token", v["token"], store.Codecs...)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(v)
	if err != nil {
		return err
	}
	resp.Body = ioutil.NopCloser(&buf)
	resp.ContentLength = int64(buf.Len())
	resp.Header.Set("Content-Length", strconv.Itoa(buf.Len()))

	// < {"token": "beyondXYZ"}
	return nil
}

func dockerHandler(w http.ResponseWriter, r *http.Request) bool {
	if r.Host != dockerHost {
		return false
	}

	ua := r.UserAgent()
	ua1 := strings.HasPrefix(ua, "docker/")
	ua2 := strings.HasPrefix(ua, "Go-")
	if !ua1 && !ua2 {
		return false
	}

	allow := r.URL.Path == "/v2/auth"
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
		return true
	}
	w.Header().Set("Docker-Distribution-Api-Version", "registry/2.0")
	w.Header().Set("WWW-Authenticate", `Bearer realm="`+*dockerScheme+`://`+r.Host+`/v2/auth",service="`+dockerHost+`"`)
	w.WriteHeader(401)
	return true
}
