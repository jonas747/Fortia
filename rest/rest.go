/*
	Simple REST server package

	If a body is provided and is needed for the specific handler the muxer will decode
	the body for you
*/
package rest

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/log"
	"io/ioutil"
	nativeLog "log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	ContentTypeJson = iota
	ContentTypeProtoBuf
)

type Server struct {
	hServer  http.Server         // The http server
	handlers map[string]*Handler // Map of handlers
	logger   *log.LogClient
}

func NewServer(addr string, logger *log.LogClient) *Server {
	server := &Server{
		handlers: make(map[string]*Handler),
		logger:   logger,
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
func (s *Server) Run() error {
	s.logger.Debug("Starting listen and serve on \"", s.hServer.Addr, "\"")
	return s.hServer.ListenAndServe()
}

// Implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	// Set some CORS headers
	// Simply setting the allow origin to * is not enough since that wont
	// Let requests be made with cookies if they are on a different origin
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

	var bodyDecoded interface{}
	if handler.BodyType != nil {
		bodyStruct := reflect.New(handler.BodyType)
		bodyDecoded = bodyStruct.Interface()
		body := r.Body
		bodyRaw, err := ioutil.ReadAll(body)
		// Decode the body json then...
		if err != nil {
			s.logger.Error(ferr.Wrapa(err, "", map[string]interface{}{"user-agent": r.UserAgent(), "remote": r.RemoteAddr}))
			HandleInternalServerError(fRequest, "Error decoding request body json")
			return
		}
		if len(bodyRaw) < 1 {
			if handler.BodyRequired {
				HandleBadRequest(fRequest, "Request body missing")
				return
			}
		} else {
			err = json.Unmarshal(bodyRaw, bodyDecoded)
			if err != nil {
				s.logger.Error(ferr.Wrapa(err, "", map[string]interface{}{"user-agent": r.UserAgent(), "remote": r.RemoteAddr, "body": bodyRaw}))
				HandleInternalServerError(fRequest, "Error decoding request json body")
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

type HandlerFunc func(*Request, interface{})
type Middleware func(*Request, interface{}) bool

type Container map[string]interface{}

func (c Container) GetInt(key string) (int, bool) {
	inter, ok := c[key]
	if !ok {
		return 0, false
	}

	val, okConv := inter.(int)
	return val, okConv
}

func (c Container) GetString(key string) (string, bool) {
	inter, ok := c[key]
	if !ok {
		return "", false
	}

	val, okConv := inter.(string)
	return val, okConv
}

func (c Container) GetSliceInt(key string) ([]int, bool) {
	inter, ok := c[key]
	if !ok {
		return []int{}, false
	}

	val, okConv := inter.([]int)
	return val, okConv
}

func (c Container) GetSliceString(key string) ([]string, bool) {
	inter, ok := c[key]
	if !ok {
		return []string{}, false
	}

	val, okConv := inter.([]string)
	return val, okConv
}

type Handler struct {
	Handler        HandlerFunc
	Method         string       // The metho ex: GET, PUT, PATCH etc..
	Path           string       // The path this handler takes action upon
	BodyType       reflect.Type // The type of the body
	BodyRequired   bool         // Wether a body is required or not
	AdditionalData Container
	MiddleWare     []Middleware
}

type Request struct {
	Server            *Server             // The rest server processing this request
	Handler           *Handler            // The handler
	Request           *http.Request       // The underlaying http.Request
	RW                http.ResponseWriter // The unerlaying http.ResponseWriter
	ResponseType      int                 // The content type of the response
	Compressed        bool                // Wether the response should be compressed or not
	AcceptedEncodings map[string]bool     // Accepted reponse encodings
	AdditionalData    Container           // Additional data, middleware can set data here
}

func NewRequest(w http.ResponseWriter, r *http.Request, server *Server, handler *Handler) *Request {
	acceptedEncodings := FindAcceptedEncodings(r)
	query := r.URL.Query()
	responseType := 0
	rtString := query.Get("api")
	switch rtString {
	case "protobuf":
		responseType = ContentTypeProtoBuf
	default:
		responseType = ContentTypeJson
	}

	req := &Request{
		Server:            server,
		Request:           r,
		RW:                w,
		ResponseType:      responseType,
		Compressed:        true,
		AcceptedEncodings: acceptedEncodings,
		Handler:           handler,
		AdditionalData:    Container(make(map[string]interface{})),
	}
	return req
}

func (r *Request) WriteResponse(msg proto.Message, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	r.RW.WriteHeader(statusCode)
	header := r.RW.Header()
	var out []byte
	switch r.ResponseType {
	case ContentTypeProtoBuf:
		header.Set("Content-Type", "application/protobuf")
		encoded, err := proto.Marshal(msg)
		if r.Server.HandleError(r, err) {
			return
		}
		out = encoded
	default:
		header.Set("Content-Type", "application/json")
		encoded, err := json.Marshal(msg)
		if r.Server.HandleError(r, err) {
			return
		}
		out = encoded
	}

	if r.Compressed {
		compressed, err := Compress(out, r.AcceptedEncodings)
		if err == nil {
			out = compressed
		} else {
			r.Server.logger.Error(err)
		}
	}

	r.RW.Write(out)
}

func FindAcceptedEncodings(r *http.Request) map[string]bool {
	encodings := make(map[string]bool)

	encodingsStr := r.Header.Get("Accept-Encoding")
	if encodingsStr == "" {
		return encodings
	}

	split := strings.Split(encodingsStr, ",")
	for _, v := range split {
		encodings[v] = true
	}

	return encodings
}
