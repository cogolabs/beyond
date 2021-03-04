package beyond

import (
	"encoding/json"
	"flag"
	"net/http"

	"github.com/gorilla/securecookie"
)

var (
	federateAccessKey = flag.String("federate-access", "", "shared secret, 64 chars, enables federation")
	federateSecretKey = flag.String("federate-secret", "", "internal secret, 64 chars")

	federateAccessCodec []securecookie.Codec
	federateSecretCodec []securecookie.Codec
)

func federateSetup() error {
	if *federateAccessKey == "" {
		return nil
	}

	federateAccessCodec = securecookie.CodecsFromPairs([]byte(*federateAccessKey)[0:31], []byte(*federateAccessKey)[32:64])
	federateSecretCodec = securecookie.CodecsFromPairs([]byte(*federateSecretKey)[0:31], []byte(*federateSecretKey)[32:64])
	return nil
}

func federate(w http.ResponseWriter, r *http.Request) {
	setCacheControl(w)

	// authenticate relying party
	next := r.URL.Query().Get("next")
	err := securecookie.DecodeMulti("next", next, &next, federateAccessCodec...)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}

	// authenticate end user
	session, err := store.Get(r, *cookieName)
	if err != nil {
		session = store.New(*cookieName)
	}
	user, _ := session.Values["user"].(string)

	// 401
	if user == "" {
		login(w, r)
		return
	}

	// issue token
	token, err := securecookie.EncodeMulti("user", user, federateSecretCodec...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 302
	http.Redirect(w, r, next+token, http.StatusFound)
}

func federateVerify(w http.ResponseWriter, r *http.Request) {
	// authenticate relying party
	token := r.URL.Query().Get("token")
	err := securecookie.DecodeMulti("user", token, &token, federateSecretCodec...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	v := map[string]string{"email": token}
	json.NewEncoder(w).Encode(v)
}
