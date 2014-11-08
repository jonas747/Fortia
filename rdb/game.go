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
	return &layer, err
}

// Returns multiple layers
func (g *GameDB) GetLayers(positions []vec.Vec3I) ([]*world.Layer, ferr.FortiaError) {
	out := make([]*world.Layer, len(positions))

	for k, v := range positions {
		layer, err := g.GetLayer(v)
		if err != nil {
			return out, err
		}
		out[k] = layer
	}
	return out, nil
}

// Saves multiple layers
func (g *GameDB) SetLayers(layers []*world.Layer) ferr.FortiaError {
	for _, v := range layers {
		err := g.SetLayer(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Returns the chunkinfo
func (g *GameDB) GetChunkInfo(pos vec.Vec2I) (*world.Chunk, ferr.FortiaError) {
	var chunk world.Chunk
	err := g.GetJson(fmt.Sprintf("c:%d:%d", pos.X, pos.Y), &chunk)
	return &chunk, err
}

func (g *GameDB) SetChunkInfo(chunk *world.Chunk) ferr.FortiaError {
	return g.SetJson(fmt.Sprintf("c:%d:%d", chunk.Position.X, chunk.Position.Y), chunk)
}

// Get and set entities
func (g *GameDB) GetEntity(id int) (map[string]string, ferr.FortiaError) {
	return g.GetHash("e:" + strconv.Itoa(id))
}

func (g *GameDB) SetEntity(id int, info map[string]interface{}) ferr.FortiaError {
	return g.SetHash("e:"+strconv.Itoa(id), info)
}
