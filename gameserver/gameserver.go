package gameserver

import (
	"github.com/jonas747/fortia/db"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/rest"
	"github.com/jonas747/fortia/world"
	"reflect"
)

var (
	logger    *log.LogClient
	authDb    *db.AuthDB
	gameDb    *db.GameDB
	gameWorld *world.World
	server    *rest.Server
)

func Run(l *log.LogClient, gdb *db.GameDB, adb *db.AuthDB, addr string) {
	l.Info("Starting gameserver")
	logger = l
	authDb = adb
	gameDb = gdb

	gameWorld = &world.World{
		Logger:      logger,
		Db:          gameDb,
		LayerSize:   20,
		WorldHeight: 200,
	}

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
	s.RegisterHandler(&rest.RestHandler{
		Path:    "/info",
		Method:  "GET",
		Handler: handleInfo,
	})
}
