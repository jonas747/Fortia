package gameserver

import (
	"github.com/jonas747/fortia/authserver"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/rest"
	"github.com/jonas747/fortia/world"
	"reflect"
)

var (
	logger    *log.LogClient
	authDb    authserver.AuthDB
	gameDb    world.GameDB
	gameWorld *world.World
	server    *rest.Server
)

func Run(l *log.LogClient, gdb world.GameDB, adb authserver.AuthDB, addr string) ferr.FortiaError {
	l.Info("Starting gameserver")
	logger = l
	authDb = adb
	gameDb = gdb

	gameWorld = &world.World{
		Logger: logger,
		Db:     gameDb,
	}
	err := gameWorld.LoadSettingsFromDb()
	if err != nil {
		l.Error(err)
		return err
	}

	server = rest.NewServer(addr, logger)
	initApi(server)
	server.Run()
	return nil
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
	s.RegisterHandler(&rest.Handler{
		Path:         "/register",
		Method:       "POST",
		Handler:      handleRegister,
		BodyType:     reflect.TypeOf(BodyRegister{}),
		BodyRequired: true,
	})
	s.RegisterHandler(&rest.Handler{
		Path:    "/login",
		Method:  "POST",
		Handler: handleLogin,
	})

	s.RegisterHandler(&rest.Handler{
		Path:    "/chunks",
		Method:  "GET",
		Handler: handleChunks,
		AdditionalData: map[string]interface{}{
			"requiredParams": []string{"x", "y"},
		},
		MiddleWare: []rest.Middleware{rest.RequiredParamsMiddleWare},
	})
	s.RegisterHandler(&rest.Handler{
		Path:    "/info",
		Method:  "GET",
		Handler: handleInfo,
	})
}
