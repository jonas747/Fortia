package db

import (
	"fmt"
	ferr "github.com/jonas747/fortia/error"
	"strconv"
)

type GameDB struct {
	*Database
}

func (g *GameDB) GetWorldInfo() (map[string]string, ferr.FortiaError) {
	return g.GetHash("worldInfo")
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
 - c:{xpos}:{ypos}:{zpos}
	json with info
 }
*/
func (g *GameDB) SetLayer(x, y, z int, layer []byte) ferr.FortiaError {
	_, err := g.Cmd("SET", fmt.Sprintf("l:%d:%d:%d", x, y, z), layer)
	return err
}

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

// // Returns a list of all chunks
// func (g *GameDB) GetChunkList() ([]string, ferr.FortiaError) {
// 	empty := []string{}
// 	reply, err := g.Cmd("KEYS", "c*")
// 	if err != nil {
// 		return empty, err
// 	}

// 	list, nerr := reply.List()
// 	if nerr != nil {
// 		return empty, ferr.Wrap(nerr, "")
// 	}
// 	return list, nil
// }

func (g *GameDB) GetEntity(id int) (map[string]string, ferr.FortiaError) {
	return g.GetHash("e:" + strconv.Itoa(id))
}

func (g *GameDB) SetEntity(id int, info map[string]interface{}) ferr.FortiaError {
	return g.SetHash("e:"+strconv.Itoa(id), info)
}
