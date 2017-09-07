package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	tokenBase = flag.String("token-base", "", "token server URL prefix (eg. https://api.github.com/user?access_token=)")

	tokenCache = cache.New(10*time.Minute, 10*time.Minute)
)

func tokenAuth(r *http.Request) string {
	if *tokenBase == "" {
		return ""
	}

	token := ""
	parts := strings.Split(r.Header.Get("Authorization"), " ")
	if len(parts) > 1 {
		if strings.ToLower(parts[0]) == "token" {
			token = parts[1]
		}
	}
	if token == "" {
		return ""
	}

	if v, ex := tokenCache.Get(token); ex {
		if v, ok := v.(string); ok {
			return v
		}
	}
	resp, err := http.Get(*tokenBase + token)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		tokenCache.Set(token, "", cache.DefaultExpiration)
		return ""
	}

	v := &tokenUser{}
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		log.Println(err)
		return ""
	}
	tokenCache.Set(token, v.Login, cache.DefaultExpiration)
	return v.Login
}

type tokenUser struct {
	Login string
	Email string
}
