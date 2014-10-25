package main

import (
	"github.com/jonas747/fortia/db"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/rest"
	"github.com/jonas747/fortia/world"
	"reflect"
)

var (
	logger    *log.LogClient
	authDb    *db.AuthDB
	gameDb    *db.GameDB
	config    *Config
	gameWorld *world.World
	server    *rest.Server
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

	l, nErr := log.NewLogClient(config.LogServer, -1, "gameServer")
	logger = l
	if nErr != nil {
		l.Error(ferr.Wrap(nErr, ""))
	}

	l.Info("(2/4) Log client init done! Creating database connection pools...")

	adb, nErr := db.NewDatabase(config.AuthDb)
	if nErr != nil {
		l.Warn("Not connected to authentication database..." + nErr.Error())
	}
	authDb = &db.AuthDB{adb}

	gdb, nErr := db.NewDatabase(config.GameDb)
	if nErr != nil {
		l.Warn("Not connected to authentication database..." + nErr.Error())
	}
	gameDb = &db.GameDB{gdb}

	gameWorld = &world.World{
		Logger:      logger,
		Db:          gameDb,
		LayerSize:   50,
		LayerHeight: 100,
	}

	l.Info("(3/4) Initializing api handlers...")
	server = rest.NewSever(":8081", l)

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


func handleRegister(w http.ResponseWriter, r *http.Request, body interface{}) {}
func handleJoin(w http.ResponseWriter, r *http.Request, body interface{}) {}
func handleUpdate(w http.ResponseWriter, r *http.Request, body interface{}) {}
func handleChunk(w http.ResponseWriter, r *http.Request, body interface{}) {}
func handleVisibleChunks(w http.ResponseWriter, r *http.Request, body interface{}) {}
*/

func initApi(s *rest.Server) {
	s.RegisterHandler(&rest.RestHandler{
		Path:            "/register",
		Method:          "POST",
		Handler:         handleRegister,
		RequiredCookies: []string{"fortia-session"},
		BodyType:        reflect.TypeOf(BodyRegister{}),
		BodyRequired:    true,
	})
	s.RegisterHandler(&rest.RestHandler{
		Path:            "/login",
		Method:          "POST",
		Handler:         handleLogin,
		RequiredCookies: []string{"fortia-session"},
	})
	s.RegisterHandler(&rest.RestHandler{
		Path:            "/layers",
		Method:          "GET",
		Handler:         handleLayers,
		RequiredCookies: []string{"fortia-session"},
		RequiredParams:  []string{"x", "y", "layers"},
	})
	s.RegisterHandler(&rest.RestHandler{
		Path:            "/visiblechunks",
		Method:          "GET",
		Handler:         handleVisibleChunks,
		RequiredCookies: []string{"fortia-session"},
	})
}
