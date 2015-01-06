package rest

import (
	"github.com/golang/protobuf/proto"
	"github.com/jonas747/fortia/errorcodes"
	"github.com/jonas747/fortia/errors"
	"github.com/jonas747/fortia/messages"

	"net/http"
)

// Returns true if there was an error, false is error == nil
func (s *Server) HandleFortiaError(r *Request, err errors.FortiaError) bool {
	if err == nil {
		return false
	}
	// Check if we need to wrap it or not
	if _, ok := err.(*RESTError); !ok {
		err = Wrap(err, r, err.GetCode(), "")
	}
	HandleInternalServerError(r, "Something went wrong: "+errorcodes.ErrorCode_name[int32(err.GetCode())], err.GetCode())

	s.logger.Error(err)

	return true
}

// Same as fortiaerror
func (s *Server) HandleError(r *Request, err error, code errorcodes.ErrorCode) bool {
	if err == nil {
		return false
	}
	return s.HandleFortiaError(r, errors.Wrap(err, code, ""))
}

func NewPlainResponse(code errorcodes.ErrorCode, msg string) *messages.PlainResponse {
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
	r.WriteResponse(NewPlainResponse(errorcodes.ErrorCode_Unauthorized, "Session cookie expired"), http.StatusUnauthorized)
}

func HandleInternalServerError(r *Request, msg string, code errorcodes.ErrorCode) {
	r.WriteResponse(NewPlainResponse(code, msg), http.StatusInternalServerError)
}

func HandleBadRequest(r *Request, msg string, code errorcodes.ErrorCode) {
	r.WriteResponse(NewPlainResponse(code, msg), http.StatusBadRequest)
}

func HandleNotFound(r *Request) {
	r.WriteResponse(NewPlainResponse(errorcodes.ErrorCode_PageNotFound, "404 - Not found"), http.StatusNotFound)
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
			HandleBadRequest(r, "Required parameter \""+reqParam+"\" Not found", errorcodes.ErrorCode_MissingRequiredParams)
			return false
		}
	}
	return true
}
