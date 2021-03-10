package beyond

import (
	"fmt"
	"net/http"
)

func NewMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(*healthPath, func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, *healthReply)
	})

	mux.HandleFunc(*host+"/federate", federate)
	mux.HandleFunc(*host+"/federate/verify", federateVerify)

	mux.HandleFunc(*host+"/launch", handleLaunch)
	mux.HandleFunc(*host+"/oidc", handleOIDC)

	for _, ds := range dockerServers {
		ds.RegisterHandlers(mux)
	}

	mux.HandleFunc("/", handler)

	return mux
}
