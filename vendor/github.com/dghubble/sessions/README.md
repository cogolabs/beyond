# sessions [![Build Status](https://travis-ci.org/dghubble/sessions.svg?branch=master)](https://travis-ci.org/dghubble/sessions) [![GoDoc](https://godoc.org/github.com/dghubble/sessions?status.svg)](https://godoc.org/github.com/dghubble/sessions)

Package `sessions` provides minimalist Go sessions, backed by `securecookie` or database stores.

### Features

* `Store` provides a predicatable interface for dealing with *individual* sessions.
    * `New` returns a new named `Session`. 
    * `Get` returns the named `Session` from the `http.Request` iff it was correctly verified and decoded. Otherwise the error is non-nil.
    * `Save` encodes and signs Session.Value data.
    * `Destroy` removes (expires) the session cookie of a given name.
* Each `Session` provides `Save` and `Destroy` convenience methods.
* Provides `CookieStore` for managing client-side secure cookies.
* Extensible for custom session database backends.

## Install

```
go get github.com/dghubble/sessions
```

## Documentation

Read [GoDoc](https://godoc.org/github.com/dghubble/sessions)

### Differences from gorilla/sessions

* Gorilla stores a context map of Requests to Sessions to abstract multiple sessions. `dghubble/sessions` provides individual sessions, leaving multiple sessions to a `multisessions` package. No Registry is needed.
* Gorilla has a depedency on `gorilla/context`, a non-standard context.
* Gorilla requires all handlers be wrapped in `context.ClearHandler` to avoid memory leaks.
* Gorilla's `Store` interface is surprising. `New` and `Get` can both possibly return a new session, a field check is needed. Some use cases expect developers to [ignore an error](https://github.com/gorilla/sessions/blob/master/doc.go#L32). `Destroy` isn't provided.

## License

[MIT License](LICENSE)
