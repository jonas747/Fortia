package rdb

import (
	"fmt"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/vec"
	"github.com/jonas747/fortia/world"
	"strconv"
)

// Basic implementation of world.GameDB
type GameDB struct {
	*Database
}

// Information about the world. Blocktypes, biomes, layer size  and such
func (g *GameDB) GetWorldInfo() (*world.WorldInfo, ferr.FortiaError) {
	var info world.WorldInfo
	err := g.GetJson("worldinfo", &info)
	return &info, err
}

func (g *GameDB) SetWorldInfo(info *world.WorldInfo) ferr.FortiaError {
	return g.SetJson("worldinfo", info)
}

// Returns the specified user's info
func (g *GameDB) GetUserInfo(user string) (*world.UserInfo, ferr.FortiaError) {
	infoHash, err := g.GetHash("u:" + user)
	if err != nil {
		return nil, err
	}
	return &world.UserInfo{
		Name: infoHash["name"],
	}, nil
}

// Sets the specified users info fields from info map to whatever is in the info map
func (g *GameDB) SetUserInfo(info *world.UserInfo) ferr.FortiaError {
	infoHash := map[string]interface{}{
		"name": info.Name,
	}
	return g.SetHash("u:"+info.Name, infoHash)
}

func (g *GameDB) GetUserEntities(user string) ([]int, ferr.FortiaError) {
	reply, err := g.Cmd("SMEMBERS", "ue:"+user)
	if err != nil {
		return []int{}, err
	}

	list, nErr := reply.List()
	if nErr != nil {
		return []int{}, ferr.Wrap(nErr, "")
	}

	intSlice := make([]int, len(list))
	for k, v := range list {
		intSlice[k], _ = strconv.Atoi(v)
	}
	return intSlice, nil
}

func (g *GameDB) EditUserEntities(user string, add, del []int) ferr.FortiaError {
	return g.EditSetI("ue:"+user, add, del)
}

/*
 - l:{xpos}:{ypos}:{zpos}
	json with info
 }
*/
func (g *GameDB) SetLayer(layer *world.Layer) ferr.FortiaError {
	return g.SetJson(fmt.Sprintf("l:%d:%d:%d", layer.Position.X, layer.Position.Y, layer.Position.Z), layer)
}

// Returns the specified layer
func (g *GameDB) GetLayer(pos vec.Vec3I) (*world.Layer, ferr.FortiaError) {
	var layer world.Layer
	err := g.GetJson(fmt.Sprintf("l:%d:%d:%d", pos.X, pos.Y, pos.Z), &layer)
	if err != nil {
		return nil, err
	}
	return &layer, nil
}
// I AM HERE!
// Returns multiple layers
func (g *GameDB) GetLayers(positions []*vec.Vec3I) ([]*], ferr.FortiaError) {
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
