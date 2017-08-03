package web

import "net/http"

// IndexHandler handles the / page.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Nothing to see here"))
}
