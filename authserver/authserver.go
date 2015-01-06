package authserver

/*
	Fortia authorisation server
	Serves a resp api
*/

import (
	"github.com/jonas747/fortia/db"
	//ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/rest"
	"reflect"
)

var (
	logger *log.LogClient
	authDb db.AuthDB
	server *rest.Server
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Runs the server witht the specified db and address
func Run(l *log.LogClient, adb db.AuthDB, addr string) {
	l.Info("Starting Authserver")
	logger = l
	authDb = adb
	server = rest.NewServer(addr, logger)
	initApi(server)
	server.ListenAndServe()
}

/*
type Handler struct {
	Handler        HandlerFunc
	Method         string       // The metho ex: GET, PUT, PATCH etc..
	Path           string       // The path this handler takes action upon
	BodyType       reflect.Type // The type of the body
	BodyRequired   bool         // Wther a body is required or not
	AdditionalData Container
	MiddleWare     []Middleware
}
*/
func initApi(s *rest.Server) {

	s.RegisterHandler(&rest.Handler{
		Handler:      rest.HandlerFunc(handleRegister),
		Method:       "POST",
		Path:         "/register",
		BodyRequired: true,
		BodyType:     reflect.TypeOf(BodyRegister{}),
	})

	s.RegisterHandler(&rest.Handler{
		Handler:      rest.HandlerFunc(handleLogin),
		Method:       "POST",
		Path:         "/login",
		BodyRequired: true,
		BodyType:     reflect.TypeOf(BodyLogin{}),
	})

	s.RegisterHandler(&rest.Handler{
		Handler: rest.HandlerFunc(handleMe),
		Method:  "GET",
		Path:    "/me",
	})

	s.RegisterHandler(&rest.Handler{
		Handler: rest.HandlerFunc(handleWorlds),
		Method:  "GET",
		Path:    "/worlds",
	})

}
