package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"path"
	"sync"
)

var (
	fenceURL     = flag.String("fence-url", "https://pages.github.com/yourcompany/beyond-config/fence.json", "")
	sitesURL     = flag.String("sites-url", "https://pages.github.com/yourcompany/beyond-config/sites.json", "")
	whitelistURL = flag.String("whitelist-url", "https://pages.github.com/yourcompany/beyond-config/whitelist.json", "")

	fence     = concurrentMapMapBool{m: map[string]map[string]bool{}}
	sites     = concurrentMapMapBool{m: map[string]map[string]bool{}}
	whitelist = concurrentMapMapBool{m: map[string]map[string]bool{}}
)

type concurrentMapMapBool struct {
	sync.RWMutex
	m map[string]map[string]bool
}

func refreshFence() error {
	resp, err := http.Get(*fenceURL)
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
	resp, err := http.Get(*sitesURL)
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
	resp, err := http.Get(*whitelistURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	whitelist.Lock()
	defer whitelist.Unlock()
	return json.NewDecoder(resp.Body).Decode(&whitelist.m)
}

func init() {
	err := refreshFence()
	if err != nil {
		log.Fatalln(err)
	}
	err = refreshSites()
	if err != nil {
		log.Fatalln(err)
	}
	err = refreshWhitelist()
	if err != nil {
		log.Fatalln(err)
	}
}

func whitelisted(host, urlpath string) bool {
	whitelist.RLock()
	allow := whitelist.m["host"][host]
	paths := whitelist.m["path"]
	whitelist.RUnlock()
	if allow {
		return true
	}
	p := path.Clean(urlpath)
	for ; p != "/"; p = path.Dir(p) {
		if paths[p] {
			allow = true
		}
	}
	return allow
}
