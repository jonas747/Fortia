package rdb

import (
	"code.google.com/p/go.crypto/bcrypt"
	ferr "github.com/jonas747/fortia/error"
	"math/rand"
	"time"
)

// Implements authserver.AuthDB
type AuthDB struct {
	*Database
}

var validSessionTokenChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// Creates a session token
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
func (a *AuthDB) GetUserInfo(user string) (map[string]string, ferr.FortiaError) {
	return a.GetHash("u:" + user)
}

// Sets the specified users info fields from info map to whatever is in the info map
func (a *AuthDB) SetUserInfo(user string, info map[string]interface{}) ferr.FortiaError {
	return a.SetHash("u:"+user, info)
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
		reply, err := a.Cmd("HGETALL", "world:"+name)
		if err != nil {
			return empty, err
		}
		hash, nerr := reply.Hash()
		if nerr != nil {
			return empty, ferr.Wrap(nerr, "")
		}

		return []map[string]string{hash}, nil
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

func (a *AuthDB) SetWorldInfo(name string, info map[string]interface{}) ferr.FortiaError {
	client, err := a.Pool.Get()
	if err != nil {
		return ferr.Wrap(err, "")
	}
	defer a.Pool.Put(client)

	reply := client.Cmd("SADD", "worlds", name)
	if reply.Err != nil {
		return ferr.Wrap(reply.Err, "")
	}
	return a.SetHash("world:"+name, info)
}
