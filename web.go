package beyond

import (
	"fmt"
	"net/http"
	"strings"
)

func NewMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(*healthPath, func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, *healthReply)
	})

	mux.HandleFunc(*host+"/launch", handleLaunch)
	mux.HandleFunc(*host+"/oidc", handleOIDC)

	mux.HandleFunc(dockerHost+"/", func(rw http.ResponseWriter, r *http.Request) {
		ua := r.UserAgent()
		ua1 := strings.HasPrefix(ua, "docker/")
		ua2 := strings.HasPrefix(ua, "Go-")
		if !ua1 && !ua2 {
			handler(rw, r)
			return
		}
		dockerHandler(rw, r)
	})

	mux.HandleFunc("/", handler)

	return mux
}
