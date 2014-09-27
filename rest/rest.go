/*
	Simple REST api package

	TODO: Use proper status codes
	TODO: Maybe more logging
	TODO: Add export function

	If a body is provided and is needed for the specific handler the muxer will decode
	the body for you
*/
package rest

import (
	"encoding/json"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/log"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

type Server struct {
	hServer  http.Server
	handlers map[string]*RestHandler
	logger   *log.LogClient
}

func NewSever(addr string, logger *log.LogClient) *Server {
	server := &Server{
		handlers: make(map[string]*RestHandler),
		logger:   logger,
	}
	hServer := http.Server{
		Addr:         addr,
		Handler:      server,
		ReadTimeout:  time.Duration(10) * time.Second,
		WriteTimeout: time.Duration(10) * time.Second,
	}
	server.hServer = hServer
	return server
}

// Starts serving
func (s *Server) Run() error {
	s.logger.Debug("Starting listen and serve on \"", s.hServer.Addr, "\"")
	return s.hServer.ListenAndServe()
}

// Implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check which handler is supposed to be used
	path := r.URL.Path
	if path == "" {
		path = "/"
	}
	s.logger.Debug("Incoming ", r.Method, " request for path \"", r.URL.Path, "\" From \"", r.RemoteAddr, "\"")
	handler, found := s.handlers[path]
	if !found {
		// 404
		HandleNotFound(w)
		return
	}

	// Maybe use method not allowed status instead?
	if r.Method != handler.Method {
		HandleNotFound(w)
		return
	}

	// Check if required params and cookies are there
	// Start with cookies
	for _, cookieName := range handler.RequiredCookies {
		_, err := r.Cookie(cookieName)
		if err != nil {
			// Maybe unathorized instead?
			HandleBadRequest(w, "Cookie \""+cookieName+"\" not found")
			return
		}
	}
	// Check url params
	params := r.URL.Query()
	for _, reqParam := range handler.RequiredParams {
		param := params.Get(reqParam)
		if param == "" {
			HandleBadRequest(w, "Param \""+reqParam+"\" Not found")
			return
		}
	}
	var bodyDecoded interface{}
	if handler.BodyType != nil {
		bodyStruct := reflect.New(handler.BodyType)
		bodyDecoded = bodyStruct.Interface()
		body := r.Body
		whole, err := ioutil.ReadAll(body)
		// Decode the body json then...
		if err != nil {
			s.logger.Error(ferr.Wrapa(err, "", map[string]interface{}{"user-agent": r.UserAgent(), "remote": r.RemoteAddr}))
			HandleInternalServerError(w, "Internal server error")
			return
		}
		if len(whole) < 1 {
			if handler.BodyRequired {
				HandleBadRequest(w, "Request body missing")
				return
			}
		} else {
			err = json.Unmarshal(whole, bodyDecoded)
			if err != nil {
				s.logger.Error(ferr.Wrapa(err, "", map[string]interface{}{"user-agent": r.UserAgent(), "remote": r.RemoteAddr, "body": whole}))
				HandleInternalServerError(w, "Error decoding request json body")
				return
			}
		}
	}
	// Finally call the handler
	handler.Handler(w, r, bodyDecoded)
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

// Registers a Handler
func (s *Server) RegisterHandler(r *RestHandler) {
	s.handlers[r.Path] = r
}

type RestHandlerFunc func(http.ResponseWriter, *http.Request, interface{})

type RestHandler struct {
	Handler         RestHandlerFunc
	Method          string       // The metho ex: GET, PUT, PATCH etc..
	RequiredParams  []string     // Required url parameters
	OptionalParams  []string     // Optional Url parameters
	RequiredCookies []string     // Required cookies
	Path            string       // The path this handler takes action upon
	BodyType        reflect.Type // The type of the body
	BodyRequired    bool         // Wther a body is required or not
}
