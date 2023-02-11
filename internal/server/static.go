package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/orca-group/spirit/internal/util"
)

// serveFiles conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func serveFiles(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

// LoadStatic renders and loads static content and templates from the `web`
func (s *Server) LoadStatic() {
	// Serve static assets
	filesDir := http.Dir("./web/static")
	serveFiles(s.Router, "/static/", filesDir)

	// Serve homepage
	s.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile("./web/index.html")

		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		w.Write(file)
	})
}
