package web

import (
	"encoding/json"
	"io"
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

// GetBinary gets the binary for the given version, or latest version if version string is empty.
func (r *Router) GetBinary(w http.ResponseWriter, req *http.Request, importURL, version string) {
	token := r.GetToken(req)
	reader, err := r.gate.GetVersionBinary(token, importURL, version)
	if err != nil {
		r.WriteError(w, 401, err.Error())
	}

	buf := make([]byte, 512)
	_, err = reader.Read(buf)
	if err != nil {
		r.WriteError(w, 500, err.Error())
	}

	ctype := http.DetectContentType(buf)
	w.Header().Add("Content-Type", ctype)
	w.Write(buf)
	io.Copy(w, reader)
}

// DeleteDisableImport decides if an import should be deleted or disabled.
func (r *Router) DeleteDisableImport(w http.ResponseWriter, req *http.Request, importURL string, delete bool) {
	token := r.GetToken(req)
	if delete {
		if err := r.gate.DeleteImport(token, importURL); err != nil {
			r.WriteError(w, 401, err.Error())
		} else {
			r.WriteOK(w)
		}
	} else {
		if err := r.gate.DisableImport(token, importURL); err != nil {
			r.WriteError(w, 401, err.Error())
		} else {
			r.WriteOK(w)
		}
	}
}

// DeleteDisableVersion decides if a version should be deleted or disabled.
func (r *Router) DeleteDisableVersion(w http.ResponseWriter, req *http.Request, importURL, version string, delete bool) {
	token := r.GetToken(req)
	if delete {
		if err := r.gate.DeleteVersion(token, importURL, version); err != nil {
			r.WriteError(w, 401, err.Error())
		} else {
			r.WriteOK(w)
		}
	} else {
		if err := r.gate.DisableVersion(token, importURL, version); err != nil {
			r.WriteError(w, 401, err.Error())
		} else {
			r.WriteOK(w)
		}
	}
}
