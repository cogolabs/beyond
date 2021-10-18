package beyond

import (
	"encoding/json"
	"flag"
	"net/http"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var (
	tokenBase = flag.String("token-base", "", "token server URL prefix (eg. https://api.github.com/user?access_token=)")
	tokenPost = flag.Bool("token-post", true, "POST token as Bearer instead of query string")

	tokenCache = cache.New(10*time.Minute, 10*time.Minute)

	tokenTypes = map[string]bool{
		"bearer": true,
		"token":  true,
	}
)

func tokenAuth(r *http.Request) string {
	if *tokenBase == "" {
		return ""
	}

	u, token, ok := r.BasicAuth()
	if ok && (token == "x-oauth-basic" || token == "") {
		token = u
	}
	if token == "" {
		parts := strings.Split(r.Header.Get("Authorization"), " ")
		if len(parts) > 1 && tokenTypes[strings.ToLower(parts[0])] {
			token = parts[1]
		}
	}
	if token == "" {
		token = r.URL.Query().Get("access_token")
	}
	if token == "" {
		return ""
	}

	if v, ex := tokenCache.Get(token); ex {
		if v, ok := v.(string); ok {
			return v
		}
	}
	var err error
	var resp *http.Response
	if *tokenPost {
		req, err := http.NewRequest("POST", *tokenBase, nil)
		if err != nil {
			Error(err)
			return ""
		}
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			Error(err)
			return ""
		}
	} else {
		resp, err = http.Get(*tokenBase + token)
		if err != nil {
			Error(err)
			return ""
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		tokenCache.Set(token, "", cache.DefaultExpiration)
		return ""
	}

	v := &tokenUser{}
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		Error(err)
		return ""
	}
	tokenCache.Set(token, v.Login, cache.DefaultExpiration)
	return v.Login
}

type tokenUser struct {
	Login string
	Email string
}
