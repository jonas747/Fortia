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
func handleLayers(w http.ResponseWriter, r *http.Request, body interface{}) {
	params := r.URL.Query()
	x, _ := strconv.Atoi(params.Get("x"))
	y, _ := strconv.Atoi(params.Get("y"))
	rawLayers := params.Get("layers")
	layerLocations := strings.Split(rawLayers, "-")

	layers := make([]*world.Layer, len(layerLocations))
	for k, v := range layerLocations {
		z, _ := strconv.Atoi(v)
		layer, err := gameWorld.GetLayer(vec.Vec3I{x, y, z})
		if server.HandleFortiaError(w, r, err) {
			return
		}
		layers[k] = layer
	}

	serialized, nErr := json.Marshal(layers)
	if server.HandleError(w, r, nErr) {
		return
	}
	w.Write(serialized)
}

func handleVisibleChunks(w http.ResponseWriter, r *http.Request, body interface{}) {}

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
