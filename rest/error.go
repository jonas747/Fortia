package rest

import (
	"fmt"
	"github.com/jonas747/fortia/errorcodes"
	"github.com/jonas747/fortia/errors"
)

type RESTError struct {
	*errors.FortiaBaseError
	Request *Request // The request that caused this error
}

func NewRestError(request *Request, code errorcodes.ErrorCode, format string, a ...interface{}) *RESTError {
	baseErr := errors.New(code, format, a...).(*errors.FortiaBaseError)
	return &RESTError{
		FortiaBaseError: baseErr,
		Request:         request,
	}
}

func Wrap(err error, request *Request, code errorcodes.ErrorCode, format string, a ...interface{}) *RESTError {
	baseErr := errors.Wrap(err, code, format, a...).(*errors.FortiaBaseError)
	return &RESTError{
		FortiaBaseError: baseErr,
		Request:         request,
	}
}

func (a *RESTError) Error() string {
	header := fmt.Sprintf("RESP Error on [%s] From [%s]\n", a.Request.Request.RequestURI, a.Request.Request.RemoteAddr)
	body := errors.DefaultError(a)
	return header + body
}
