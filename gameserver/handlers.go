package gameserver

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/jonas747/fortia/messages"
	//ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/rest"
	"github.com/jonas747/fortia/vec"
	"net/http"
	"strconv"
	"strings"
)

func handleRegister(r *rest.Request, body interface{}) {}

func handleLogin(r *rest.Request, body interface{}) {}

func handleUpdate(r *rest.Request, body interface{}) {}

func serialize(data interface{}, out chan []byte) {
	serialized, err := json.Marshal(data)
	if err != nil {
		out <- []byte{}
	}
	out <- serialized
}

// /chunks
func handleChunks(r *rest.Request, body interface{}) {
	params := r.Request.URL.Query()
	xList := strings.Split(params.Get("x"), ",")
	yList := strings.Split(params.Get("y"), ",")

	positions := make([]vec.Vec2I, len(xList))
	for k, _ := range xList {
		x, _ := strconv.Atoi(xList[k])
		y, _ := strconv.Atoi(yList[k])
		positions[k] = vec.Vec2I{x, y}
	}
	chunks := make([]*messages.Chunk, len(positions))
	for i, v := range positions {
		chunk, err := gameWorld.GetChunk(v)
		if err != nil {
			if err.GetCode() == 404 {
				continue
			}
		}
		if server.HandleFortiaError(r, err) {
			return
		}
		if chunk == nil { // should not happend, but still make sure
			continue
		}
		chunks[i] = chunk.RawChunk
	}
	wireResp := &messages.ChunksResponse{
		Chunks: chunks,
	}
	r.WriteResponse(wireResp, http.StatusOK)
}

// /info
func handleInfo(r *rest.Request, body interface{}) {
	info, err := gameDb.GetWorldInfo()
	if server.HandleFortiaError(r, err) {
		return
	}

	wireInfo := &messages.WorldSettings{
		Size:        proto.Int(info.Size),
		ChunkWidth:  proto.Int(info.ChunkWidth),
		ChunkHeight: proto.Int(info.ChunkHeight),
	}
	wireResp := &messages.WorldSettingsResponse{
		Settings: wireInfo,
	}
	r.WriteResponse(wireResp, http.StatusOK)
}
