package rdb

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jonas747/fortia/db"
	"github.com/jonas747/fortia/errorcodes"
	"github.com/jonas747/fortia/errors"
	"github.com/jonas747/fortia/messages"
	"github.com/jonas747/fortia/vec"
	"strconv"
)

// Basic implementation of world.GameDB
type GameDB struct {
	*Database
}

// Information about the world. Blocktypes, biomes, layer size  and such
func (g *GameDB) GetWorldSettings() (*messages.WorldSettings, errors.FortiaError) {
	var info messages.WorldSettings
	err := g.GetProto("worldinfo", &info)
	return &info, err
}

func (g *GameDB) SetWorldSettings(info *messages.WorldSettings) errors.FortiaError {
	return g.SetProto("worldinfo", info)
}

// Returns the specified user's info
func (g *GameDB) GetUserInfo(user string) (*db.GameUserInfo, errors.FortiaError) {
	infoHash, err := g.GetHash("u:" + user)
	if err != nil {
		return nil, err
	}
	return &db.GameUserInfo{
		Name: infoHash["name"],
	}, nil
}

// Sets the specified users info fields from info map to whatever is in the info map
func (g *GameDB) SetUserInfo(info *db.GameUserInfo) errors.FortiaError {
	infoHash := map[string]interface{}{
		"name": info.Name,
	}
	return g.SetHash("u:"+info.Name, infoHash)
}

func (g *GameDB) GetUserEntities(user string) ([]int, errors.FortiaError) {
	list, err := g.GetList("SMEMBERS", "ue:"+user)
	if err != nil {
		return []int{}, err
	}

	intSlice := make([]int, len(list))
	for k, v := range list {
		intSlice[k], _ = strconv.Atoi(v)
	}
	return intSlice, nil
}

func (g *GameDB) EditUserEntities(user string, add, del []int) errors.FortiaError {
	return g.EditSetI("ue:"+user, add, del)
}

// Returns the chunkinfo
func (g *GameDB) GetChunk(pos vec.Vec2I) (*messages.Chunk, errors.FortiaError) {
	var chunk messages.Chunk
	err := g.GetProto(fmt.Sprintf("c:%d:%d", pos.X, pos.Y), &chunk)
	return &chunk, err
}

func (g *GameDB) SetChunk(chunk *messages.Chunk) errors.FortiaError {
	return g.SetProto(fmt.Sprintf("c:%d:%d", chunk.GetX(), chunk.GetY()), chunk)
}

// Get and set entities
func (g *GameDB) GetEntity(id int) (map[string]string, errors.FortiaError) {
	return g.GetHash("e:" + strconv.Itoa(id))
}

func (g *GameDB) SetEntity(id int, info map[string]interface{}) errors.FortiaError {
	return g.SetHash("e:"+strconv.Itoa(id), info)
}

func (g *GameDB) PushAction(action *messages.Action, tick int) errors.FortiaError {
	serialized, nErr := proto.Marshal(action)
	if nErr != nil {
		return errors.New(errorcodes.ErrorCode_ProtoEncodeErr, nErr.Error(), nil)
	}
	_, err := g.Cmd("SADD", fmt.Sprintf("actionQueue:%d:%d", tick, action.GetKind()), serialized)
	return err
}

func (g *GameDB) PopAction(tick, kind int) (*messages.Action, errors.FortiaError) {
	raw, err := g.GetBytes("SPOP", fmt.Sprintf("actionQueue:%d:%d", tick, kind))
	if err != nil {
		return nil, err
	}

	var action *messages.Action
	nErr := proto.Unmarshal(raw, action)
	if nErr != nil {
		return nil, errors.New(errorcodes.ErrorCode_ProtoDecodeErr, nErr.Error(), nil)
	}

	return action, nil
}

func (g *GameDB) IncrTick() (int, errors.FortiaError) {
	n, err := g.GetInt("INCR", "tick")
	if err != nil {
		return n, err
	}
	return n, nil
}

func (g *GameDB) GetCurrentTick() (int, errors.FortiaError) {
	n, err := g.GetInt("GET", "tick")
	return n, err
}
