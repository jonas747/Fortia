/*
	Simple REST server package

	If a body is provided and is needed for the specific handler the muxer will decode
	the body for you

*/
package rest

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/jonas747/fortia/errorcodes"
	"net/http"
	"reflect"
	"strings"
)

const (
	ContentTypeJson = iota
	ContentTypeProtoBuf
)

type HandlerFunc func(*Request, interface{})

// Middleware should return false if the rest server should not call the handler
type Middleware func(*Request, interface{}) bool

// Simple container
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

// A rest handler
type Handler struct {
	Handler        HandlerFunc
	Method         string       // The metho ex: GET, PUT, PATCH etc..
	Path           string       // The path this handler takes action upon
	BodyType       reflect.Type // The type of the body
	BodyRequired   bool         // Wether a body is required or not
	AdditionalData Container    // Additional data used by middleware(required params for the reqparams middleware for example)
	MiddleWare     []Middleware // A slice of middleware to be executed before the handler itself
}

type Request struct {
	Server            *Server             `json:"-"` // The rest server processing this request
	Handler           *Handler            `json:"-"` // The handler
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
	header := r.RW.Header()
	var out []byte
	switch r.ResponseType {
	case ContentTypeProtoBuf:
		header.Set("Content-Type", "application/protobuf")
		encoded, err := proto.Marshal(msg)
		if r.Server.HandleError(r, err, errorcodes.ErrorCode_ProtoEncodeErr) {
			return
		}
		out = encoded
	default:
		header.Set("Content-Type", "application/json")
		encoded, err := json.Marshal(msg)
		if r.Server.HandleError(r, err, errorcodes.ErrorCode_JsonEncodeErr) {
			return
		}
		out = encoded
	}

	if r.Compressed {
		compressed, err := Compress(out, r)
		if err == nil {
			out = compressed
		} else {
			r.Server.logger.Error(err)
		}
	}

	r.RW.WriteHeader(statusCode)
	r.RW.Write(out)
}

// To be used with logging and such
func (r *Request) CreateDataMap() map[string]interface{} {
	dmap := make(map[string]interface{})
	dmap["requestUri"] = r.Request.RequestURI
	dmap["remoteAddr"] = r.Request.RemoteAddr
	dmap["header"] = r.Request.Header
	return dmap
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
