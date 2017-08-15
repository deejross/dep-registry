package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/deejross/dep-registry/gate"
)

// Router object.
type Router struct {
	gate *gate.Gate
}

// NewRouter returns a new Router.
func NewRouter(gate *gate.Gate) *Router {
	return &Router{
		gate: gate,
	}
}

// ServeHTTP handles the routing.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/api/v1/") {
		path := strings.Split(req.URL.Path[8:], "/")
		r.API(w, req, path)
	} else {
		r.Static(w, req)
	}
}

// WriteError writes an error to the response.
func (r *Router) WriteError(w http.ResponseWriter, status int, err string) {
	errM := map[string]string{
		"error": err,
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errM)
}

// API handles API requests.
func (r *Router) API(w http.ResponseWriter, req *http.Request, path []string) {
	switch path[0] {
	case "auth":
		if len(path) > 1 {
			switch path[1] {
			case "login":
				r.Login(w, req)
			}
		}
	case "projects":
		if len(path) > 1 {
			// TODO: finish
		}
	default:
		r.WriteError(w, http.StatusNotFound, "Not found")
	}
}

// Static handles static requests.
func (r *Router) Static(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Nothing to see here...yet"))
}
