package authserver

/*
	Fortia authorisation server
	Serves a resp api
*/

import (
	//"github.com/jonas747/fortia/rdb"
	//ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/rest"
	"reflect"
)

var (
	logger *log.LogClient
	authDb AuthDB
	server *rest.Server
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Runs the server witht the specified db and address
func Run(l *log.LogClient, adb AuthDB, addr string) {
	l.Info("Starting Authserver")
	logger = l
	authDb = adb
	server = rest.NewServer(addr, logger)
	initApi(server)
	server.Run()
}

/*
	Handler         RestHandlerFunc
	Method          string       // The metho ex: GET, POST, PUT, PATCH etc..
	RequiredParams  []string     // Required url parameters
	OptionalParams  []string     // Optional Url parameters
	RequiredCookies []string     // Required cookies
	Path            string       // The path this handler takes action upon
	BodyType        reflect.Type // The type of the body
	BodyRequired    bool         // Wther a body is required or not
*/
func initApi(s *rest.Server) {

	s.RegisterHandler(&rest.RestHandler{
		Handler:      rest.RestHandlerFunc(handleRegister),
		Method:       "POST",
		Path:         "/register",
		BodyRequired: true,
		BodyType:     reflect.TypeOf(BodyRegister{}),
	})

	s.RegisterHandler(&rest.RestHandler{
		Handler:      rest.RestHandlerFunc(handleLogin),
		Method:       "POST",
		Path:         "/login",
		BodyRequired: true,
		BodyType:     reflect.TypeOf(BodyLogin{}),
	})

	s.RegisterHandler(&rest.RestHandler{
		Handler:         rest.RestHandlerFunc(handleMe),
		Method:          "GET",
		Path:            "/me",
		RequiredCookies: []string{"fortia-session"},
	})

	s.RegisterHandler(&rest.RestHandler{
		Handler: rest.RestHandlerFunc(handleWorlds),
		Method:  "GET",
		Path:    "/worlds",
	})

}
