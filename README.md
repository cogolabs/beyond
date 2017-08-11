# transcend

## Usage
```
$ go get -u github.com/cogolabs/transcend

$ transcend -help
Usage of ./transcend:
  -client-id string
    	OIDC client ID (default "f8b8b020-4ec2-0135-6452-027de1ec0c4e43491")
  -client-secret string
    	OIDC client secret (default "cxLF74XOeRRFDJbKuJpZAOtL4pVPK1t2XGVrDbe5Rx0Uij1LS2e9k7opZI6jQzHC")
  -cookie-age int
    	MaxAge setting in seconds (default 21600)
  -cookie-domain string
    	beyond cookie domain (default ".colofoo.net")
  -cookie-key1 string
    	key1 for cookie crypto pair (default "t8yG1gmeEyeb7pQpw544UeCTyDfPkE6u")
  -cookie-key2 string
    	key2 of cookie crypto pair (default "Q599vrruZRhLFC144thCRZpyHM7qGDjt")
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
