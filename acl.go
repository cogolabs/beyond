package beyond

import (
	"encoding/json"
	"flag"
	"net/http"
	"path"
	"sync"
)

var (
	fenceURL     = flag.String("fence-url", "", "URL to user fencing config (eg. https://github.com/myorg/beyond-config/main/raw/fence.json)")
	sitesURL     = flag.String("sites-url", "", "URL to allowed sites config (eg. https://github.com/myorg/beyond-config/main/raw/sites.json)")
	allowlistURL = flag.String("allowlist-url", "", "URL to site allowlist (eg. https://github.com/myorg/beyond-config/main/raw/allowlist.json)")

	fence     = concurrentMapMapBool{m: map[string]map[string]bool{}}
	sites     = concurrentMapMapBool{m: map[string]map[string]bool{}}
	allowlist = concurrentMapMapBool{m: map[string]map[string]bool{}}

	httpACL = &http.Client{}
)

type concurrentMapMapBool struct {
	sync.RWMutex
	m map[string]map[string]bool
}

func refreshFence() error {
	if *fenceURL == "" {
		return nil
	}

	resp, err := httpACL.Get(*fenceURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	d := map[string][]string{}
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return err
	}
	for k, v := range d {
		if _, ok := fence.m[k]; !ok {
			fence.m[k] = map[string]bool{}
		}
		for _, v := range v {
			fence.m[k][v] = true
		}
	}
	return nil
}

func refreshSites() error {
	if *sitesURL == "" {
		return nil
	}

	resp, err := httpACL.Get(*sitesURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	d := map[string][]string{}
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return err
	}
	sites.Lock()
	defer sites.Unlock()
	for k, v := range d {
		if _, ok := sites.m[k]; !ok {
			sites.m[k] = map[string]bool{}
		}
		for _, v := range v {
			sites.m[k][v] = true
		}
	}
	return nil
}

func refreshAllowlist() error {
	if *allowlistURL == "" {
		return nil
	}

	resp, err := httpACL.Get(*allowlistURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	allowlist.Lock()
	defer allowlist.Unlock()
	return json.NewDecoder(resp.Body).Decode(&allowlist.m)
}

func allowlisted(r *http.Request) bool {
	allowlist.RLock()
	allow := allowlist.m["host"][r.Host]
	hostM := allowlist.m["host:method"][r.Host+":"+r.Method]
	paths := allowlist.m["path"]
	allowlist.RUnlock()
	if allow || hostM {
		return true
	}
	p := path.Clean(r.URL.Path)
	for ; p != "/"; p = path.Dir(p) {
		if paths[p] {
			allow = true
		}
	}
	return allow
}

func deny(r *http.Request, user string) bool {
	fence.RLock()
	zones, ok := fence.m[user]
	fence.RUnlock()

	if !ok || len(zones) < 1 {
		return false
	}

	sites.RLock()
	m := sites.m
	sites.RUnlock()

	for k := range zones {
		if m[k]["https://"+r.Host] || m[k]["http://"+r.Host] {
			return false
		}
	}
	return true
}
