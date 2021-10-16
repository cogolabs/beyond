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
	"github.com/pkg/errors"
)

var (
	samlCert = flag.String("saml-cert-file", "example/myservice.cert", "SAML SP path to cert.pem")
	samlKey  = flag.String("saml-key-file", "example/myservice.key", "SAML SP path to key.pem")

	samlID  = flag.String("saml-entity-id", "", "SAML SP entity ID (blank defaults to beyond-host)")
	samlIDP = flag.String("saml-metadata-url", "", "SAML metadata URL from IdP (blank disables SAML)")

	samlNIDF = flag.String("saml-nameid-format", "unspecified", "SAML SP NameID format: {unspecified, email, persistent, transient}")
	samlSign = flag.Bool("saml-sign-requests", true, "SAML SP signs authentication requests")

	samlSP *samlsp.Middleware
)

func samlSetup() error {
	if *samlIDP == "" {
		return nil
	}
	if *samlID == "" {
		*samlID = *host
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
		EntityID:    *samlID,
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

	switch *samlNIDF {
	case "email":
		samlSP.ServiceProvider.AuthnNameIDFormat = saml.EmailAddressNameIDFormat
	case "persistent":
		samlSP.ServiceProvider.AuthnNameIDFormat = saml.PersistentNameIDFormat
	case "transient":
		samlSP.ServiceProvider.AuthnNameIDFormat = saml.TransientNameIDFormat
	case "unspecified":
		samlSP.ServiceProvider.AuthnNameIDFormat = saml.UnspecifiedNameIDFormat
	case "":
	default:
		return errors.Errorf("invalid nameid format: \"%s\"", *samlNIDF)
	}
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
