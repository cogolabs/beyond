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
	whitelistURL = flag.String("whitelist-url", "", "URL to site whitelist (eg. https://github.com/myorg/beyond-config/main/raw/whitelist.json)")

	fence     = concurrentMapMapBool{m: map[string]map[string]bool{}}
	sites     = concurrentMapMapBool{m: map[string]map[string]bool{}}
	whitelist = concurrentMapMapBool{m: map[string]map[string]bool{}}

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

func refreshWhitelist() error {
	if *whitelistURL == "" {
		return nil
	}

	resp, err := httpACL.Get(*whitelistURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	whitelist.Lock()
	defer whitelist.Unlock()
	return json.NewDecoder(resp.Body).Decode(&whitelist.m)
}

func whitelisted(r *http.Request) bool {
	whitelist.RLock()
	allow := whitelist.m["host"][r.Host]
	hostM := whitelist.m["host:method"][r.Host+":"+r.Method]
	paths := whitelist.m["path"]
	whitelist.RUnlock()
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
