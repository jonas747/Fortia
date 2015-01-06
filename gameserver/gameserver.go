package gameserver

import (
	"github.com/jonas747/fortia/db"
	"github.com/jonas747/fortia/errors"
	"github.com/jonas747/fortia/game"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/messages"
	"github.com/jonas747/fortia/rest"
	"reflect"
)

var (
	logger    *log.LogClient
	authDb    db.AuthDB
	gameDb    db.GameDB
	gameWorld *game.World
	server    *rest.Server
)

func Run(l *log.LogClient, gdb db.GameDB, adb db.AuthDB, addr string) errors.FortiaError {
	l.Info("Starting gameserver")
	logger = l
	authDb = adb
	gameDb = gdb

	gameWorld = &game.World{
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
	server.ListenAndServe()
	return nil
}

func initApi(s *rest.Server) {
	s.RegisterHandler(&rest.Handler{
		Path:         "/register",
		Method:       "POST",
		Handler:      handleRegister,
		BodyType:     reflect.TypeOf(messages.BodyGameRegister{}),
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

	s.RegisterHandler(&rest.Handler{
		Path:         "/placeunit",
		Method:       "POST",
		Handler:      handlePlaceUnit,
		BodyRequired: true,
		BodyType:     reflect.TypeOf(messages.BodyPlaceUnit{}),
	})
}
