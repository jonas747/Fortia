package db

import (
	"code.google.com/p/go.crypto/bcrypt"
	ferr "github.com/jonas747/fortia/error"
	"math/rand"
	"time"
)

type AuthDB struct {
	*Database
}

var validSessionTokenChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func createSessionToken(length int) string {
	seed := time.Now().Nanosecond()
	rand.Seed(int64(seed))
	/*
		arr := make([]byte, length)
		for i := 0; i < length; i++ {
			arr[i] = byte(rand.Intn(255))
		}
		return string(arr)
	*/
	finalString := ""
	for i := 0; i < length; i++ {
		rNum := rand.Intn(len(validSessionTokenChars))
		char := string(validSessionTokenChars[rNum])
		finalString += char
	}
	return finalString
}

// Logs the specified user in returning a session token
func (a *AuthDB) LoginUser(user string) (string, ferr.FortiaError) {
	token := createSessionToken(64)

	_, err := a.Cmd("SETEX", "t:"+user+":"+token, 3600 /*An hour*/, 1)
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
func (a *AuthDB) GetUserInfo(user string) (map[string]string, ferr.FortiaError) {
	reply, err := a.Cmd("HGETALL", "u:"+user)
	if err != nil {
		return emptyStrStrMap, nil
	}
	hash, convErr := reply.Hash()
	if convErr != nil {
		return emptyStrStrMap, ferr.Wrapa(convErr, "Error converting to map", map[string]interface{}{"user": user})
	}
	return hash, nil
}

// Sets the specified users info fields from info map to whatever is in the info map
func (a *AuthDB) SetUserInfo(user string, info map[string]interface{}) ferr.FortiaError {
	_, err := a.Cmd("HMSET", "u:"+user, info)
	return err
}

// Checks if the session token is existing, returning the ttl if found or -1 if not
func (a *AuthDB) CheckSessionToken(user, token string) (int, ferr.FortiaError) {
	reply, err := a.Cmd("TTL", "t:"+user+":"+token)
	if err != nil {
		return -1, err
	}
	duration, convErr := reply.Int()
	if convErr != nil {
		return -1, ferr.Wrapa(err, "Error converting to int", map[string]interface{}{"user": user, "token": token})
	}
	return duration, nil
}

// Extends the session token for n seconds
func (a *AuthDB) ExtendSessionToken(user, token string, duration int) ferr.FortiaError {
	_, err := a.Cmd("EXPIRE", "t:"+user+":"+token, duration)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthDB) GetWorldInfo(name string) ([]map[string]string, ferr.FortiaError) {
	empty := make([]map[string]string, 0)
	if name != "" {
		// TODO
		return empty, nil
	}
	conn, err := a.Pool.Get()
	if err != nil {
		return empty, ferr.Wrap(err, "")
	}
	defer a.Pool.Put(conn)

	replyList := conn.Cmd("SMEMBERS", "worlds")
	list, err := replyList.List()
	if err != nil {
		return empty, ferr.Wrap(err, "")
	}
	servers := make([]map[string]string, 0)
	for _, v := range list {
		reply := conn.Cmd("HGETALL", "world:"+v)
		hash, err := reply.Hash()
		if err != nil {
			continue
		}
		hash["name"] = v
		servers = append(servers, hash)
	}
	return servers, nil
}
