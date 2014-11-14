package gameserver

import (
	"encoding/json"
	//ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/vec"
	//"github.com/jonas747/fortia/world"
	"bytes"
	"compress/gzip"
	"net/http"
	"strconv"
	"strings"
)

func writeCompressed(w http.ResponseWriter, r *http.Request, body []byte) {
	encodingsStr := r.Header.Get("Accept-Encoding")
	split := strings.Split(encodingsStr, ",")
	var buffer bytes.Buffer
	encoding := ""
	for _, v := range split {
		if v == "gzip" {
			writer := gzip.NewWriter(&buffer)
			total := 0
			for total < len(body) {
				n, _ := writer.Write(body[total:])
				total += n
			}
			writer.Close()
			encoding = "gzip"
			break
		}
	}
	// Write uncompressed if no compression is supoprted
	if encoding == " " {
		w.Write(body)
		return
	}
	w.Header().Set("Content-Encoding", encoding)
	w.Header().Set("Content-Type", "application/json")
	logger.Debug("Serialized length: ", buffer.Len()/1000, "k")
	total := int64(0)
	for total < int64(buffer.Len()) {
		n, _ := buffer.WriteTo(w)
		total += n
	}
}

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
	writeCompressed(w, r, serialized)
}

func serialize(data interface{}, out chan []byte) {
	serialized, err := json.Marshal(data)
	if err != nil {
		out <- []byte{}
	}
	out <- serialized
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
	dataChan := make(chan []byte, 5)
	numChunks := 0
	for _, v := range positions {
		chunk, err := gameWorld.GetChunk(v.X, v.Y, true, true)
		if err != nil {
			if err.GetMessage() == "404" {
				continue
			}
		}
		if server.HandleFortiaError(w, r, err) {
			return
		}
		dmap := map[string]interface{}{
			"Layers":        chunk.Layers,
			"Position":      chunk.Position,
			"Biome":         chunk.Biome.Id,
			"VisibleLayers": chunk.VisibleLayers,
		}
		numChunks++
		go serialize(dmap, dataChan)
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")
	firstElen := true
	for i := 0; i < numChunks; i++ {
		chunk := <-dataChan
		if len(chunk) < 1 {
			continue
		}
		if !firstElen {
			buffer.WriteString(",")
		}
		firstElen = false
		buffer.Write(chunk)
	}
	buffer.WriteString("]")
	// serialized, nErr := json.Marshal(chunks)
	// if server.HandleError(w, r, nErr) {
	// 	return
	// }
	serialized := buffer.Bytes()
	writeCompressed(w, r, serialized)
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
