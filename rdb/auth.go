package rdb

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/jonas747/fortia/authserver"
	ferr "github.com/jonas747/fortia/error"
	"math/rand"
	"strconv"
	"time"
)

// Implements authserver.AuthDB
type AuthDB struct {
	*Database
}

var validSessionTokenChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// Creates a session token
// TODO: Maybe not seed every time?
func createSessionToken(length int) string {
	seed := time.Now().Nanosecond()
	rand.Seed(int64(seed))
	finalString := ""
	for i := 0; i < length; i++ {
		rNum := rand.Intn(len(validSessionTokenChars))
		char := string(validSessionTokenChars[rNum])
		finalString += char
	}
	return finalString
}

// Logs the specified user in returning a session token
func (a *AuthDB) LoginUser(user string, duration int) (string, ferr.FortiaError) {
	token := createSessionToken(64)
	token = user + ";" + token
	_, err := a.Cmd("SETEX", "t:"+token, duration, user)
	if err != nil {
		return token, err
	}
	return token, nil
}

// Returns true if the users password is correct
func (a *AuthDB) CheckUserPw(user, pw string) (bool, ferr.FortiaError) {
	pwReply, err := a.Cmd("HGET", "u:"+user, "pw")
	if err != nil {
		return false, err
	}
	correctPwHash := pwReply.String()
	if correctPwHash == "" {
		return false, nil
	}

	berr := bcrypt.CompareHashAndPassword([]byte(correctPwHash), []byte(pw))
	if berr != nil {
		return false, nil
	}
	return true, nil
}

// Returns the specified user's info
func (a *AuthDB) GetUserInfo(user string) (*authserver.UserInfo, ferr.FortiaError) {
	infoHash, err := a.GetHash("u:" + user)
	if err != nil {
		return nil, err
	}

	// Retrieve the servers this user is on
	reply, err := a.Cmd("SMEMBERS", "uw:"+user)
	if err != nil {
		return nil, err
	}

	list, nErr := reply.List()
	if nErr != nil {
		return nil, ferr.Wrap(nErr, "")
	}

	donorLvl, _ := strconv.Atoi(infoHash["donorLvl"])
	role, _ := strconv.Atoi(infoHash["role"])
	pwhash := []byte(infoHash["pwHash"])

	userInfo := &authserver.UserInfo{
		Name:   user,
		Email:  infoHash["email"],
		PwHash: pwhash,

		Worlds:   list,
		Role:     role,
		DonorLvl: donorLvl,
	}

	return userInfo, nil
}

// Sets the specified users info fields from info map to whatever is in the info map
func (a *AuthDB) SetUserInfo(info *authserver.UserInfo) ferr.FortiaError {
	infoHash := map[string]interface{}{
		"name":     info.Name,
		"email":    info.Email,
		"pwHash":   info.PwHash,
		"role":     info.Role,
		"donorLvl": info.DonorLvl,
	}

	conn, nErr := a.Database.Pool.Get()
	if nErr != nil {
		return ferr.Wrap(nErr, "")
	}
	defer a.Database.Pool.Put(conn)

	// Set the info hash
	reply := conn.Cmd("HMSET", "u:"+info.Name, infoHash)
	if reply.Err != nil {
		return ferr.Wrap(reply.Err, "")
	}

	// Set the worlds to a set if there are any worlds this player is in
	if len(info.Worlds) > 0 {
		argList := []interface{}{"uw:" + info.Name}
		for _, v := range info.Worlds {
			argList = append(argList, v)
		}

		reply = conn.Cmd("SADD", argList...)
		if reply.Err != nil {
			return ferr.Wrap(reply.Err, "")
		}
	}
	return nil
}

// implements authserver.AuthDB.EditUserInfoFields
func (a *AuthDB) EditUserInfoFields(user string, fields map[string]interface{}) ferr.FortiaError {
	return a.SetHash("u:"+user, fields)
}

func (a *AuthDB) EditUserWorlds(user string, add []string, del []string) ferr.FortiaError {
	conn, nErr := a.Database.Pool.Get()
	if nErr != nil {
		return ferr.Wrap(nErr, "")
	}
	defer a.Pool.Put(conn)

	if len(add) > 0 {
		argList := []interface{}{"uw:" + user}
		for _, v := range add {
			argList = append(argList, v)
		}
		conn.Cmd("SADD", argList...)
	}

	if len(del) > 0 {
		argList := []interface{}{"uw:" + user}
		for _, v := range del {
			argList = append(argList, v)
		}
		conn.Cmd("SREM", argList...)
	}

	return nil
}

// Checks if the session token is existing, returning the user it belong to if found or "" if not
func (a *AuthDB) CheckSessionToken(token string) (string, ferr.FortiaError) {
	reply, err := a.Cmd("GET", "t:"+token)
	if err != nil {
		return "", err
	}
	owner := reply.String()
	return owner, nil
}

// Extends the session token for n seconds
func (a *AuthDB) ExtendSessionToken(token string, duration int) ferr.FortiaError {
	_, err := a.Cmd("EXPIRE", "t:"+token, duration)
	return err
}

func (a *AuthDB) GetWorldListing() ([]*authserver.WorldInfo, ferr.FortiaError) {
	reply, err := a.Cmd("SMEMBERS", "worlds")
	if err != nil {
		return []*authserver.WorldInfo{}, err
	}
	listing, nErr := reply.List()
	if nErr != nil {
		return []*authserver.WorldInfo{}, ferr.Wrap(nErr, "")
	}

	worlds := make([]*authserver.WorldInfo, len(listing))
	for i, wname := range listing {
		info, err := a.GetWorldInfo(wname)
		if err != nil {
			return []*authserver.WorldInfo{}, err
		}
		worlds[i] = info
	}

	return worlds, nil
}

func (a *AuthDB) GetWorldInfo(world string) (*authserver.WorldInfo, ferr.FortiaError) {
	infoHash, err := a.GetHash("world:" + world)
	if err != nil {
		return nil, err
	}

	players, _ := strconv.Atoi(infoHash["players"])
	size, _ := strconv.Atoi(infoHash["size"])
	started, _ := strconv.Atoi(infoHash["started"])

	info := &authserver.WorldInfo{
		Name:    world,
		Started: started,
		Players: players,
		Size:    size,
	}

	return info, nil
}

func (a *AuthDB) SetWorldInfo(info *authserver.WorldInfo) ferr.FortiaError {
	infoHash := map[string]interface{}{
		"name":    info.Name,
		"players": info.Players,
		"size":    info.Size,
		"started": info.Started,
	}

	client, err := a.Pool.Get()
	if err != nil {
		return ferr.Wrap(err, "")
	}
	defer a.Pool.Put(client)

	reply := client.Cmd("SADD", "worlds", info.Name)
	if reply.Err != nil {
		return ferr.Wrap(reply.Err, "")
	}
	return a.SetHash("world:"+info.Name, infoHash)
}

// Overwrites the stored fields with fields provided
func (a *AuthDB) EditWorldInfo(world string, fields map[string]interface{}) ferr.FortiaError {
	return a.SetHash("world:"+world, fields)
}
