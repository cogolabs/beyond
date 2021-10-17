package beyond

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	homeURL = flag.String("home-url", "https://google.com", "redirect users here from root")
)

// NewMux mounts all configured web handlers
func NewMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(*healthPath, func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, *healthReply)
	})

	mux.HandleFunc(*host+"/federate", federate)
	mux.HandleFunc(*host+"/federate/verify", federateVerify)

	mux.HandleFunc(*host+"/launch", handleLaunch)
	mux.HandleFunc(*host+"/oidc", handleOIDC)
	if samlSP != nil {
		mux.HandleFunc(*host+"/saml/", samlSP.ServeHTTP)
	}
	mux.Handle(*host+"/", http.RedirectHandler(*homeURL, http.StatusTemporaryRedirect))

	for _, ds := range dockerServers {
		ds.RegisterHandlers(mux)
	}

	mux.HandleFunc("/", handler)

	return mux
}
