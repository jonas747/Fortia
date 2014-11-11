package gameserver

import (
	"encoding/json"
	"github.com/jonas747/fortia/vec"
	"github.com/jonas747/fortia/world"
	"net/http"
	"strconv"
	"strings"
)

func handleRegister(w http.ResponseWriter, r *http.Request, body interface{}) {
	w.Write([]byte("{\"ok\": true}"))
}

func handleLogin(w http.ResponseWriter, r *http.Request, body interface{}) {
	w.Write([]byte("{\"ok\": true}"))
}

func handleUpdate(w http.ResponseWriter, r *http.Request, body interface{}) {}

// /layers
// TODO add a world.getRawLayers function so i dont decode and then encode the json
func handleLayers(w http.ResponseWriter, r *http.Request, body interface{}) {
	params := r.URL.Query()
	xList := strings.Split(params.Get("x"), ",")
	yList := strings.Split(params.Get("y"), ",")
	zList := strings.Split(params.Get("z"), ",")

	positions := make([]vec.Vec3I, len(xList))
	for k, _ := range xList {
		x, _ := strconv.Atoi(xList[k])
		y, _ := strconv.Atoi(yList[k])
		z, _ := strconv.Atoi(zList[k])
		positions[k] = vec.Vec3I{x, y, z}
	}

	layers, err := gameDb.GetLayers(positions)
	if server.HandleFortiaError(w, r, err) {
		return
	}

	serialized, nErr := json.Marshal(layers)
	if server.HandleError(w, r, nErr) {
		return
	}
	w.Write(serialized)
}

// /chunks
func handleChunks(w http.ResponseWriter, r *http.Request, body interface{}) {
	params := r.URL.Query()
	xList := strings.Split(params.Get("x"), ",")
	yList := strings.Split(params.Get("y"), ",")

	positions := make([]vec.Vec2I, len(xList))
	for k, _ := range xList {
		x, _ := strconv.Atoi(xList[k])
		y, _ := strconv.Atoi(yList[k])
		positions[k] = vec.Vec2I{x, y}
	}
	chunks := make([]*world.Chunk, len(positions))
	for k, v := range positions {
		chunk, err := gameWorld.GetChunk(v.X, v.Y, true, true)
		if server.HandleFortiaError(w, r, err) {
			return
		}
		chunks[k] = chunk
	}
	serialized, nErr := json.Marshal(chunks)
	if server.HandleError(w, r, nErr) {
		return
	}
	w.Write(serialized)
}

// /info
func handleInfo(w http.ResponseWriter, r *http.Request, body interface{}) {
	infoHash, err := gameDb.GetWorldInfo()
	if server.HandleFortiaError(w, r, err) {
		return
	}
	serialized, nErr := json.Marshal(&infoHash)
	if server.HandleError(w, r, nErr) {
		return
	}
	w.Write(serialized)
}
