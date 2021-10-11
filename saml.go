package beyond

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"net/http"
	"net/url"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
)

var (
	samlCert = flag.String("saml-cert-file", "example/myservice.cert", "Path to SP cert.pem")
	samlKey  = flag.String("saml-key-file", "example/myservice.key", "Path to SP key.pem")
	samlIDP  = flag.String("saml-metadata-url", "", "Metadata URL for IdP (blank disables SAML)")
	samlSign = flag.Bool("saml-sign-requests", true, "Sign Requests to IdP")

	samlSP *samlsp.Middleware
)

func samlSetup() error {
	if *samlIDP == "" {
		return nil
	}

	keyPair, err := tls.LoadX509KeyPair(*samlCert, *samlKey)
	if err != nil {
		return err
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return err
	}

	idpMetadataURL, err := url.Parse(*samlIDP)
	if err != nil {
		return err
	}
	idpMetadata, err := samlsp.FetchMetadata(
		context.Background(), http.DefaultClient,
		*idpMetadataURL)
	if err != nil {
		return err
	}

	rootURL, _ := url.Parse("https://" + *host)
	if err != nil {
		return err
	}
	samlSP, err = samlsp.New(samlsp.Options{
		EntityID:    *host,
		SignRequest: *samlSign,
		URL:         *rootURL,

		Certificate: keyPair.Leaf,
		IDPMetadata: idpMetadata,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),

		AllowIDPInitiated: true,
	})
	if err != nil {
		return err
	}

	samlSP.ServiceProvider.AuthnNameIDFormat = saml.PersistentNameIDFormat
	return nil
}

func samlFilter(w http.ResponseWriter, r *http.Request) bool {
	samlSession, _ := samlSP.Session.GetSession(r)
	if _, ok := samlSession.(samlsp.SessionWithAttributes); !ok {
		return false
	}
	samlAttributes := samlSession.(samlsp.SessionWithAttributes).GetAttributes()
	email := samlAttributes.Get("email")
	if email == "" {
		return false
	}

	session, err := store.Get(r, *cookieName)
	if err != nil {
		session = store.New(*cookieName)
	}
	session.Values["user"] = email
	session.Save(w)
	samlSP.Session.DeleteSession(w, r)
	return true
}
