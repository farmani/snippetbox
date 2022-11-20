package main

import (
	"net/http"

	"github.com/justinas/alice" // New import
	"githun.com/farmani/snippetbox/internal/config"
)

// The routes() method returns a servemux containing our application routes.
func routes(cfgHost config.Host) http.Handler {

	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir(cfgHost.StaticDir))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register the other application routes as normal.
	mux.Handle("/", home())
	mux.Handle("/snippet/view", snippetView())
	mux.Handle("/snippet/create", snippetCreate())

	standard := alice.New(
		recoverPanic,
		logRequest,
		secureHeaders,
	)
	return standard.Then(mux)
}
