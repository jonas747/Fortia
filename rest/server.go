package rest

import (
	"encoding/json"
	"github.com/jonas747/fortia/errorcodes"
	"github.com/jonas747/fortia/errors"
	"github.com/jonas747/fortia/log"
	"io/ioutil"
	nativeLog "log"
	"net/http"
	"reflect"
	"time"
)

type Server struct {
	hServer  http.Server         // The http server
	handlers map[string]*Handler // Map of handlers
	logger   *log.LogClient
	Stop     chan bool
}

func NewServer(addr string, logger *log.LogClient) *Server {
	server := &Server{
		handlers: make(map[string]*Handler),
		logger:   logger,
		Stop:     make(chan bool),
	}
	hServer := http.Server{
		Addr:         addr,
		Handler:      server,
		ReadTimeout:  time.Duration(10) * time.Second,
		WriteTimeout: time.Duration(10) * time.Second,
		ErrorLog:     nativeLog.New(logger, "", 0),
	}
	server.hServer = hServer
	return server
}

// Start serving
func (s *Server) ListenAndServe() errors.FortiaError {
	s.logger.Info("Starting listen and serve on \"", s.hServer.Addr, "\"")
	listener, err := Listen(s.hServer.Addr)
	if err != nil {
		return err
	}
	nErr := s.hServer.Serve(listener)
	if nErr != nil {
		return errors.Wrap(nErr, errorcodes.ErrorCode_ServeErr, "")
	}
	return nil
}

// Implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	// Set some CORS headers
	// Simply setting the allow origin to * is not enough since that wont
	// Let requests be made with cookies if they are on a different origin
	// One may argue that this is bad security, but trusting the client to
	// provide the correct origin header is also a security flaw so this
	// isnt gonna change much
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Check which handler is supposed to be used
	path := r.URL.Path
	if path == "" {
		path = "/"
	}
	defer func() {
		taken := time.Since(started)
		s.logger.Debug("Handled ", r.Method, " request. Path: \"", r.URL.Path, "\", From \"", r.RemoteAddr, "\" Took: ", taken.String())
	}()

	handler, found := s.handlers[path]
	fRequest := NewRequest(w, r, s, handler)
	if !found {
		// 404
		HandleNotFound(fRequest)
		return
	}

	// Maybe use method not allowed status instead?
	if r.Method != handler.Method {
		HandleNotFound(fRequest)
		return
	}
	// Decode the body json if there exists one
	var bodyDecoded interface{}
	if handler.BodyType != nil {
		bodyStruct := reflect.New(handler.BodyType)
		bodyDecoded = bodyStruct.Interface()
		body := r.Body
		bodyRaw, err := ioutil.ReadAll(body) // <- Can run out of memory, so should fix this probably
		// Decode the body json then...
		if s.HandleError(fRequest, err, errorcodes.ErrorCode_IoReadErr) {
			return
		}
		if len(bodyRaw) < 1 {
			if handler.BodyRequired {
				HandleBadRequest(fRequest, "Request body is reuired", errorcodes.ErrorCode_RequestBodyRequired)
				// Should probably log this...
				return
			}
		} else {
			err = json.Unmarshal(bodyRaw, bodyDecoded)
			if s.HandleError(fRequest, err, errorcodes.ErrorCode_JsonDecodeErr) {
				return
			}
		}
	}
	// Call the middlewares
	for _, v := range handler.MiddleWare {
		cont := v(fRequest, bodyDecoded)
		if !cont {
			return
		}
	}

	// Finally call the handler
	handler.Handler(fRequest, bodyDecoded)
}

// Registers a Handler
func (s *Server) RegisterHandler(r *Handler) {
	s.handlers[r.Path] = r
}
