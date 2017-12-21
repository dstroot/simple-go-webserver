package health

import (
	"io"
	"net/http"
	"sync/atomic"
)

// Handler supports a liveness probe. It is a simple handler which
// always return response code 200
func Handler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

// Ready supports a readiness probe.  For the readiness probe we might
// need to wait for some event (e.g. the database is ready) to be able
// to serve traffic. We return 200 only if the variable "isReady" is true.
func Ready(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"ready": true}`)
	}
}

// HandlerFunc returns the info HTTP Handler.
func HandlerFunc() http.Handler {
	return http.HandlerFunc(Handler)
}

// ReadyFunc returns the info HTTP Handler.
func ReadyFunc(isReady *atomic.Value) http.Handler {
	return http.HandlerFunc(Ready(isReady))
}