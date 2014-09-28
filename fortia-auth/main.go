package main

import (
	"github.com/jonas747/fortia/db"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/rest"
	//"net/http"
	//"net/url"
	"reflect"
)

var (
	logger *log.LogClient
	authDb *db.AuthDB
	gameDb *db.GameDB
	config *Config
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	c, err := loadConfig("config.json")
	panicErr(err)
	config = c

	l, nErr := log.NewLogClient(config.LogServer, -1, "authAPI")
	logger = l
	if nErr != nil {
		l.Error(ferr.Wrap(nErr, ""))
	}

	l.Info("(2/4) Log client init successful! Creating database connection pools...")

	adb, nErr := db.NewDatabase(config.AuthDb)
	if nErr != nil {
		l.Warn("Not connected to database..." + nErr.Error())
	}

	authDb = &db.AuthDB{adb}

	l.Info("(3/4) Initializing api handlers...")
	server := rest.NewSever(":8080", l)

	initApi(server)
	l.Info("(4/4) Starting http server...")
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
		Handler:         rest.RestHandlerFunc(handleGetInfo),
		Method:          "GET",
		Path:            "/getinfo",
		RequiredCookies: []string{"fortia-session"},
		RequiredParams:  []string{"user"},
	})

}
