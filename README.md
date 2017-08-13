[![Build Status](https://travis-ci.org/cogolabs/transcend.svg?branch=master)](https://travis-ci.org/cogolabs/transcend)
[![Coverage Status](https://img.shields.io/coveralls/cogolabs/transcend.svg)](https://coveralls.io/github/cogolabs/transcend)
[![Docker Build Status](https://img.shields.io/docker/build/cogolabs/transcend.svg)](https://hub.docker.com/r/cogolabs/transcend/)
[![Go Report Card](https://goreportcard.com/badge/github.com/cogolabs/transcend)](https://goreportcard.com/report/github.com/cogolabs/transcend)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# transcend
## Install
```
$ docker pull cogolabs/transcend
```
or:
```
$ go get -u -x github.com/cogolabs/transcend
```
## Usage
```
$ docker run --rm -p 80:80 cogolabs/transcend transcend --help
Usage of ./transcend:
  -client-id string
    	OIDC client ID (default "f8b8b020-4ec2-0135-6452-027de1ec0c4e43491")
  -client-secret string
    	OIDC client secret (default "cxLF74XOeRRFDJbKuJpZAOtL4pVPK1t2XGVrDbe5Rx0Uij1LS2e9k7opZI6jQzHC")
  -cookie-age int
    	MaxAge setting in seconds (default 21600)
  -cookie-domain string
    	session cookie domain (default ".colofoo.net")
  -cookie-key1 string
    	key1 for cookie crypto pair (default "t8yG1gmeEyeb7pQpw544UeCTyDfPkE6u")
  -cookie-key2 string
    	key2 of cookie crypto pair (default "Q599vrruZRhLFC144thCRZpyHM7qGDjt")
  -cookie-name string
    	session cookie name (default "transcend")
  -fence-url string
    	URL to user fencing config (default "https://pages.github.com/yourcompany/beyond-config/fence.json")
  -host string
    	hostname of self, eg. when generating OAuth redirect URLs (default "beyond.colofoo.net")
  -http string
    	listen address (default ":80")
  -insecure-skip-verify
    	allow TLS backends without valid certificates
  -oidc-issuer string
    	issuer URL provided by IdP (default "https://yourcompany.onelogin.com/oidc")
  -sites-url string
    	URL to allowed sites config (default "https://pages.github.com/yourcompany/beyond-config/sites.json")
  -websocket-compression
    	allow websocket transport compression (gorilla/experimental)
  -whitelist-url string
    	URL to site whitelist (default "https://pages.github.com/yourcompany/beyond-config/whitelist.json")
```
