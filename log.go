package beyond

import (
	"flag"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	Error      = log.Error
	WithError  = log.WithError
	WithField  = log.WithField
	WithFields = log.WithFields

	logHTTP = flag.Bool("log-http", false, "log HTTP requests to stdout")
	logJSON = flag.Bool("log-json", false, "logrus json output")
	logXFF  = flag.Bool("log-xff", false, "include X-Forwarded-For in logs")
)

func logSetup() error {
	if *logJSON {
		log.SetFormatter(&log.JSONFormatter{})
	}
	return nil
}

func logRoundtrip(resp *http.Response) {
	d := map[string]interface{}{
		"date": time.Now().Format(time.RFC3339),
		"user": resp.Request.Header.Get(*headerPrefix + "-User"),

		"method": resp.Request.Method,
		"host":   resp.Request.Host,
		"path":   resp.Request.URL.Path,
		"query":  resp.Request.URL.RawQuery,

		"code":     resp.StatusCode,
		"len":      resp.ContentLength,
		"location": resp.Header.Get("Location"),
		"proto":    resp.Proto,
		"server":   resp.Header.Get("Server"),
		"type":     resp.Header.Get("Content-Type"),

		// "req.header":  resp.Request.Header,
		// "resp.header": resp.Header,
	}
	if *logXFF {
		d["xff"] = resp.Request.Header.Get("X-Forwarded-For")
	}
	for k, v := range d {
		if v == "" {
			delete(d, k)
		}
	}

	if *logHTTP {
		WithFields(d).Info("HTTP")
	}
}
