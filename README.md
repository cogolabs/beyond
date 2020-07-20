[![Build Status](https://travis-ci.org/cogolabs/beyond.svg?branch=master)](https://travis-ci.org/cogolabs/beyond)
[![codecov](https://codecov.io/gh/cogolabs/beyond/branch/master/graph/badge.svg)](https://codecov.io/gh/cogolabs/beyond)
[![Docker Build Status](https://img.shields.io/docker/cloud/build/cogolabs/beyond.svg)](https://hub.docker.com/r/cogolabs/beyond/)
[![Go Report Card](https://goreportcard.com/badge/github.com/cogolabs/beyond)](https://goreportcard.com/report/github.com/cogolabs/beyond)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# beyond
Control access to services beyond your perimeter network. Deploy with split-DNS to alleviate VPN in a zero-trust transition. Inspired by Google BeyondCorp research: https://research.google.com/pubs/pub45728.html

## Features
- Authenticate via:
  - OpenID Connect
  - OAuth2 Tokens
- Automate Configuration w/ https://your.json
- Customize Nexthop Learning (via Favorite Ports: 443, 80, ...)
- Supports WebSockets
- Supports GitHub Enterprise
- Supports Private Docker Registry APIs (v2)
- Analytics with ElasticSearch

## Install
```
$ docker pull cogolabs/beyond
```
or:
```
$ go get -u -x github.com/cogolabs/beyond
```
## Usage
```
$ docker run --rm -p 80:80 cogolabs/beyond httpd --help
Usage of ./httpd:
  -401-code int
    	status to respond when a user needs authentication (default 418)
  -404-message string
    	message to use for unlisted hosts when learning is disabled or fails (default "Please contact your network administrators to whitelist this system.")
  -client-id string
    	OIDC client ID (default "f8b8b020-4ec2-0135-6452-027de1ec0c4e43491")
  -client-secret string
    	OIDC client secret (default "cxLF74XOeRRFDJbKuJpZAOtL4pVPK1t2XGVrDbe5Rx0Uij1LS2e9k7opZI6jQzHC")
  -cookie-age int
    	MaxAge setting in seconds (default 21600)
  -cookie-domain string
    	session cookie domain (default ".colofoo.net")
  -cookie-key1 string
    	key1 of cookie crypto pair (example: "t8yG1gmeEyeb7pQpw544UeCTyDfPkE6u")
  -cookie-key2 string
    	key2 of cookie crypto pair (example: "Q599vrruZRhLFC144thCRZpyHM7qGDjt")
  -cookie-name string
    	session cookie name (default "beyond")
  -docker-auth-scheme string
    	(only for testing) (default "https")
  -docker-url string
    	 (default "https://docker.colofoo.net")
  -error-color string
    	css h1 color for errors (default "#69b342")
  -error-email string
    	address for help (eg. support@mycompany.com)
  -error-plain
    	disable html on error pages
  -fence-url string
    	URL to user fencing config (eg. https://pages.github.com/yourcompany/beyond-config/fence.json)
  -header-prefix string
    	prefix extra headers with this string (default "Beyond")
  -health-path string
    	URL of the health endpoint (default "/healthz/ping")
  -health-reply string
    	response body of the health endpoint (default "ok")
  -host string
    	hostname of self, eg. when generating OAuth redirect URLs (default "beyond.colofoo.net")
  -host-masq string
    	rewrite nexthop hosts (format: from1=to1,from2=to2)
  -http string
    	listen address (default ":80")
  -insecure-skip-verify
    	allow TLS backends without valid certificates
  -learn-dial-timeout duration
    	skip port after this connection timeout (default 5s)
  -learn-http-ports string
    	after HTTPS, try these HTTP ports (csv) (default "80,8080,6000,6060,7000,8000,9000,9200,15672")
  -learn-https-ports string
    	try learning these backend HTTPS ports (csv) (default "443,4443,6443,8443,9443")
  -learn-nexthops
    	set false to require explicit whitelisting (default true)
  -log-elastic string
    	csv of elasticsearch servers
  -log-elastic-interval duration
    	how often to commit bulk updates (default 1s)
  -log-elastic-prefix string
    	insert this on the front of elastic indexes (default "beyond")
  -log-elastic-workers int
    	bulk commit workers (default 3)
  -log-http
    	log HTTP requests to stdout
  -log-json
    	logrus json output
  -log-xff
    	include X-Forwarded-For in logs
  -oidc-issuer string
    	issuer URL provided by IdP (default "https://yourcompany.onelogin.com/oidc")
  -server-idle-timeout duration
    	maximum amount of time to wait for the next request when keep-alives are enabled (default 3m0s)
  -server-read-timeout duration
    	maximum duration for reading the entire request, including the body (default 1m0s)
  -server-write-timeout duration
    	maximum duration before timing out writes of the response (default 2m0s)
  -sites-url string
    	URL to allowed sites config (eg. https://pages.github.com/yourcompany/beyond-config/sites.json)
  -token-base string
    	token server URL prefix (eg. https://api.github.com/user?access_token=)
  -websocket-compression
    	allow websocket transport compression (gorilla/experimental)
  -whitelist-url string
    	URL to site whitelist (eg. https://pages.github.com/yourcompany/beyond-config/whitelist.json)
```
