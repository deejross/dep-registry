package web

import (
	"encoding/json"
	"net/http"
)

// Login and generate a token.
func (r *Router) Login(w http.ResponseWriter, req *http.Request) {
	username, password, ok := req.BasicAuth()
	if !ok {
		r.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	token, err := r.gate.Login(username, password)
	if err != nil {
		r.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// GetToken gets the token from the Request.
func (r *Router) GetToken(req *http.Request) string {
	a := req.Header.Get("Authorization")
	expected := "Bearer"
	if len(a) > len(expected)+2 {
		return a[len(expected)+1:]
	}
	return ""
}
