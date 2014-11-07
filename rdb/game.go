package rdb

import (
	"fmt"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/vec"
	"strconv"
)

// Basic implementation of world.GameDB
type GameDB struct {
	*Database
}

// Information about the world. Blocktypes, biomes, layer size  and such
func (g *GameDB) GetWorldInfo() (map[string]string, ferr.FortiaError) {
	return g.GetHash("worldInfo")
}

func (g *GameDB) SetWorldInfo(info map[string]interface{}) ferr.FortiaError {
	return g.SetHash("worldInfo", info)
}

// Returns the specified user's info
func (g *GameDB) GetUserInfo(user string) (map[string]string, ferr.FortiaError) {
	return g.GetHash("u:" + user)
}

// Sets the specified users info fields from info map to whatever is in the info map
func (g *GameDB) SetUserInfo(user string, info map[string]interface{}) ferr.FortiaError {
	return g.SetHash("u:"+user, info)
}

/*
 - l:{xpos}:{ypos}:{zpos}
	json with info
 }
*/
func (g *GameDB) SetLayer(x, y, z int, layer []byte) ferr.FortiaError {
	_, err := g.Cmd("SET", fmt.Sprintf("l:%d:%d:%d", x, y, z), layer)
	return err
}

// Returns the specified layer
func (g *GameDB) GetLayer(x, y, z int) ([]byte, ferr.FortiaError) {
	reply, err := g.Cmd("GET", fmt.Sprintf("l:%d:%d:%d", x, y, z))
	if err != nil {
		return []byte{}, err
	}
	out, nErr := reply.Bytes()
	if nErr != nil {
		return []byte{}, ferr.Wrap(nErr, "")
	}
	return out, nil
}

// Returns multiple layers
func (g *GameDB) GetLayers(positions []*vec.Vec3I) ([][]byte, ferr.FortiaError) {
	out := make([][]byte, len(positions))

	keys := make([]interface{}, len(positions))
	for k, v := range positions {
		keys[k] = fmt.Sprintf("l:%d:%d:%d", v.X, v.Y, v.Z)
	}
	reply, err := g.Cmd("MGET", keys...)
	if err != nil {
		return out, err
	}

	for k, v := range reply.Elems {
		if v == nil {
			continue
		}
		b, err := v.Bytes()
		if err != nil {
			return out, ferr.Wrap(err, "")
		}
		out[k] = b
	}

	return out, nil
}

// Stores information about a chunk
// c:{x}:{y}
// Returns the chunkinfo, wether the chunk exists, and any errors
func (g *GameDB) GetChunkInfo(x, y int) ([]byte, bool, ferr.FortiaError) {
	reply, err := g.Cmd("GET", fmt.Sprintf("c:%d:%d", x, y))
	if err != nil {
		return []byte{}, false, err
	}

	out, nErr := reply.Bytes()
	if nErr != nil {
		return []byte{}, false, nil
	}
	return out, true, nil
}

func (g *GameDB) SetChunkInfo(x, y int, info []byte) ferr.FortiaError {
	_, err := g.Cmd("SET", fmt.Sprintf("c:%d:%d", x, y), info)
	return err
}

// Get and set entities
func (g *GameDB) GetEntity(id int) (map[string]string, ferr.FortiaError) {
	return g.GetHash("e:" + strconv.Itoa(id))
}

func (g *GameDB) SetEntity(id int, info map[string]interface{}) ferr.FortiaError {
	return g.SetHash("e:"+strconv.Itoa(id), info)
}
