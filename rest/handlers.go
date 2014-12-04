package rest

import (
	"code.google.com/p/goprotobuf/proto"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/messages"

	"net/http"
)

// Returns true if there was an error, false is error == nil
func (s *Server) HandleFortiaError(r *Request, err ferr.FortiaError) bool {
	// Add some additional details
	if err == nil {
		return false
	}
	err.SetData("remoteaddr", r.Request.RemoteAddr)
	err.SetData("path", r.Request.URL.Path)
	s.logger.Error(err)
	HandleBadRequest(r, "Something went wrong... (Contact admin)")
	return true
}

// Same as fortiaerror
func (s *Server) HandleError(r *Request, err error) bool {
	if err == nil {
		return false
	}
	return s.HandleFortiaError(r, ferr.Wrap(err, ""))
}

func NewPlainResponse(code int, msg string) *messages.PlainResponse {
	r := &messages.PlainResponse{}

	if code != 0 || msg != "" {
		err := &messages.Error{
			Code: proto.Int32(int32(code)),
			Text: proto.String(msg),
		}
		r.Error = err
	}
	return r
}

func (s *Server) HandleUnauthorized(r *Request, user string) {
	s.logger.Warna("Unauthorized", map[string]interface{}{"remoteaddr": r.Request.RemoteAddr, "user": user, "path": r.Request.URL.Path})
	r.WriteResponse(NewPlainResponse(ErrInvalidSessionCookie, "Session cookie expired"), http.StatusUnauthorized)
}

func HandleInternalServerError(r *Request, msg string) {
	r.WriteResponse(NewPlainResponse(0, msg), http.StatusInternalServerError)
}

func HandleBadRequest(r *Request, msg string) {
	r.WriteResponse(NewPlainResponse(0, msg), http.StatusBadRequest)
}

func HandleNotFound(r *Request) {
	r.WriteResponse(NewPlainResponse(0, "404 - Not found"), http.StatusNotFound)
}

func RequiredParamsMiddleWare(r *Request, b interface{}) bool {
	params := r.Request.URL.Query()

	requiredParams, ok := r.Handler.AdditionalData.GetSliceString("requiredParams")
	if !ok {
		return true
	}

	for _, reqParam := range requiredParams {
		param := params.Get(reqParam)
		if param == "" {
			HandleBadRequest(r, "Required parameter \""+reqParam+"\" Not found")
			return false
		}
	}
	return true
}
