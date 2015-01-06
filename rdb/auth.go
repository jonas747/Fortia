package rdb

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/jonas747/fortia/db"
	"github.com/jonas747/fortia/errorcodes"
	"github.com/jonas747/fortia/errors"
	"math/rand"
	"strconv"
	"time"
)

// Implements db.AuthDB
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
func (a *AuthDB) LoginUser(user string, duration int) (string, errors.FortiaError) {
	token := createSessionToken(64)
	token = user + ";" + token
	_, err := a.Cmd("SETEX", "t:"+token, duration, user)
	if err != nil {
		return token, err
	}
	return token, nil
}

// Returns true if the users password is correct
func (a *AuthDB) CheckUserPw(user, pw string) (bool, errors.FortiaError) {
	dbPw, err := a.GetBytes("HGET", "u:"+user, "pwHash")
	if err != nil {
		return false, err
	}
	berr := bcrypt.CompareHashAndPassword(dbPw, []byte(pw))
	if berr != nil {
		if berr == bcrypt.ErrMismatchedHashAndPassword {
			// incorrect password
			return false, nil
		}
		return false, errors.Wrap(berr, errorcodes.ErrorCode_BCryptErr, "")
	}
	return true, nil
}

// Returns the specified user's info
func (a *AuthDB) GetUserInfo(user string) (*db.AuthUserInfo, errors.FortiaError) {
	infoHash, err := a.GetHash("u:" + user)
	if err != nil {
		return nil, err
	}

	// Retrieve the worlds this user is on
	list, err := a.GetList("SMEMBERS", "uw:"+user)
	if err != nil {
		return nil, err
	}

	donorLvl, _ := strconv.Atoi(infoHash["donorLvl"])
	role, _ := strconv.Atoi(infoHash["role"])
	pwhash := []byte(infoHash["pwHash"])

	userInfo := &db.AuthUserInfo{
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
func (a *AuthDB) SetUserInfo(info *db.AuthUserInfo) errors.FortiaError {
	infoHash := map[string]interface{}{
		"name":     info.Name,
		"email":    info.Email,
		"pwHash":   info.PwHash,
		"role":     info.Role,
		"donorLvl": info.DonorLvl,
	}

	// Set the info hash
	_, err := a.Cmd("HMSET", "u:"+info.Name, infoHash)
	if err != nil {
		return err
	}

	// Set the worlds to a set if there are any worlds this player is in
	if len(info.Worlds) > 0 {
		argList := make([]interface{}, len(info.Worlds))
		for k, v := range info.Worlds {
			argList[k] = v
		}
		err := a.SetSet("uw:"+info.Name, argList)
		if err != nil {
			return err
		}
	}
	return nil
}

// implements db.AuthDB.EditUserInfoFields
func (a *AuthDB) EditUserInfoFields(user string, fields map[string]interface{}) errors.FortiaError {
	return a.SetHash("u:"+user, fields)
}

func (a *AuthDB) EditUserWorlds(user string, add []string, del []string) errors.FortiaError {
	return a.EditSetS("uw", add, del)
}

// Checks if the session token is existing, returning the user it belong to if found or "" if not
func (a *AuthDB) CheckSessionToken(token string) (string, errors.FortiaError) {
	owner, err := a.GetString("GET", "t:"+token)
	if err != nil {
		return "", err
	}
	return owner, nil
}

// Extends the session token for n seconds
func (a *AuthDB) ExtendSessionToken(token string, duration int) errors.FortiaError {
	_, err := a.Cmd("EXPIRE", "t:"+token, duration)
	return err
}

func (a *AuthDB) GetWorldListing() ([]*db.WorldInfo, errors.FortiaError) {
	listing, err := a.GetList("SMEMBERS", "worlds")
	if err != nil {
		return []*db.WorldInfo{}, err
	}

	worlds := make([]*db.WorldInfo, len(listing))
	for i, wname := range listing {
		info, err := a.GetWorldInfo(wname)
		if err != nil {
			return []*db.WorldInfo{}, err
		}
		worlds[i] = info
	}

	return worlds, nil
}

func (a *AuthDB) GetWorldInfo(world string) (*db.WorldInfo, errors.FortiaError) {
	infoHash, err := a.GetHash("world:" + world)
	if err != nil {
		return nil, err
	}

	players, _ := strconv.Atoi(infoHash["players"])
	size, _ := strconv.Atoi(infoHash["size"])
	started, _ := strconv.Atoi(infoHash["started"])

	info := &db.WorldInfo{
		Name:    world,
		Started: started,
		Players: players,
		Size:    size,
	}

	return info, nil
}

func (a *AuthDB) SetWorldInfo(info *db.WorldInfo) errors.FortiaError {
	infoHash := map[string]interface{}{
		"name":    info.Name,
		"players": info.Players,
		"size":    info.Size,
		"started": info.Started,
	}

	_, err := a.Cmd("SADD", "worlds", info.Name)
	if err != nil {
		return err
	}
	return a.SetHash("world:"+info.Name, infoHash)
}

// Overwrites the stored fields with fields provided
func (a *AuthDB) EditWorldInfo(world string, fields map[string]interface{}) errors.FortiaError {
	return a.SetHash("world:"+world, fields)
}
