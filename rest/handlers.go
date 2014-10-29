package rest

import (
	ferr "github.com/jonas747/fortia/error"
	"net/http"
)

// Returns true if there was an error, false is error == nil
func (s *Server) HandleFortiaError(w http.ResponseWriter, r *http.Request, err ferr.FortiaError) bool {
	// Add some additional details
	if err == nil {
		return false
	}
	err.SetData("remoteaddr", r.RemoteAddr)
	err.SetData("path", r.URL.Path)
	s.logger.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(ApiError(ErrServerError, "Something went wrong back in serverland :("))
	return true
}

// Same as fortiaerror
func (s *Server) HandleError(w http.ResponseWriter, r *http.Request, err error) bool {
	if err == nil {
		return false
	}
	return s.HandleFortiaError(w, r, ferr.Wrap(err, ""))
}

func (s *Server) HandleUnauthorized(w http.ResponseWriter, r *http.Request, user string) {
	s.logger.Warna("Unauthorized", map[string]interface{}{"remoteaddr": r.RemoteAddr, "user": user, "path": r.URL.Path})
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(ApiError(ErrInvalidSessionCookie, "Session expired"))
}

func HandleInternalServerError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("{\"error\": \"" + msg + "\"}"))
}

func HandleBadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"error\":\"" + msg + "\"}"))
}

func HandleNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("<h1>404 - Not found buddy</h1>"))
}
