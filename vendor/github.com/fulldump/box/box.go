package box

import (
	"fmt"
	"net/http"
	"os"
)

type B struct {
	// R is the root resource in box
	*R
	HttpHandler            http.Handler
	HandleResourceNotFound any
	HandleMethodNotAllowed any
}

func (b *B) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.HttpHandler.ServeHTTP(w, r)
}

func NewBox() *B {
	b := &B{
		R:                      NewResource(),
		HandleResourceNotFound: DefaultHandlerResourceNotFound,
		HandleMethodNotAllowed: DefaultHandlerMethodNotAllowed,
	}
	b.HttpHandler = Box2Http(b)

	return b
}

// Serve is a helper that creates a new http.Server and call to its method
// ListenAndServe on address :8080
// Deprecated: Use the analogous name b.ListenAndServe
func (b *B) Serve() {
	b.ListenAndServe()
}

// ListenAndServe is a helper that creates a new http.Server and call to its
// method ListenAndServe on address :8080
func (b *B) ListenAndServe() error {

	server := &http.Server{
		Addr:    ":8080",
		Handler: b,
	}
	fmt.Fprintln(os.Stderr, "Listening to", server.Addr)

	return server.ListenAndServe()
}
