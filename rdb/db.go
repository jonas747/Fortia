package rdb

/*
	This is a redis implementation of authserver.AuthDB and world.GameDB
*/
import (
	"encoding/json"
	"github.com/fzzy/radix/extra/pool"
	"github.com/fzzy/radix/redis"
	ferr "github.com/jonas747/fortia/error"
)

var (
	emptyStrStrMap = make(map[string]string)
)

// Base database
type Database struct {
	Pool *pool.Pool
}

func NewDatabase(addr string) (*Database, error) {
	p, err := pool.NewPool("tcp", addr, 10)
	if err != nil {
		return nil, err
	}
	return &Database{p}, nil
}

// Same as redis.Client.Cmd but uses a connection from a pool
func (db *Database) Cmd(cmd string, args ...interface{}) (*redis.Reply, ferr.FortiaError) {
	reply, client, err := db.CmdChain(cmd, args)
	if client != nil {
		db.Pool.Put(client)
	}
	return reply, err
}

// Same as cmd but returns the connection after the command for more use
func (db *Database) CmdChain(cmd string, args ...interface{}) (*redis.Reply, *redis.Client, ferr.FortiaError) {
	client, err := db.Pool.Get()
	if err != nil {
		return nil, nil, ferr.Wrap(err, "Error Get db client")
	}
	reply := client.Cmd(cmd, args)
	if reply.Err != nil {
		db.Pool.Put(client)
		return nil, nil, ferr.Wrapa(reply.Err, "Redis cmd", map[string]interface{}{"cmd": cmd, "args": args})
	}
	return reply, client, nil
}

// Represents a redis command
type RedisCmd struct {
	Cmd  string
	Args []interface{}
}

// Todo use redis transaction
func (db *Database) MultiCmd(cmds []RedisCmd) ([]*redis.Reply, ferr.FortiaError) {
	client, err := db.Pool.Get()
	if err != nil {
		return nil, ferr.Wrap(err, "Error Get db client")
	}
	defer db.Pool.Put(client)
	replies := make([]*redis.Reply, 0)
	for _, cmd := range cmds {
		reply := client.Cmd(cmd.Cmd, cmd.Args...)
		replies = append(replies, reply)
	}
	return replies, nil
}

// Does HGETALL on "key"
func (db *Database) GetHash(key string) (map[string]string, ferr.FortiaError) {
	reply, err := db.Cmd("HGETALL", key)
	if err != nil {
		return emptyStrStrMap, err
	}
	hMap, errConv := reply.Hash()
	if errConv != nil {
		return emptyStrStrMap, ferr.Wrapa(errConv, "Redis reply conversion", map[string]interface{}{"key": key, "type": "hash"})
	}
	return hMap, err
}

// Does HMSET on "key" with all the fialds in info map
func (db *Database) SetHash(key string, info map[string]interface{}) ferr.FortiaError {
	_, err := db.Cmd("HMSET", key, info)
	return err
}

// Decodes the json at "key" info "val", returns any errors if any
func (db *Database) GetJson(key string, val interface{}) ferr.FortiaError {
	reply, err := db.Cmd("GET", key)
	if err != nil {
		return err
	}

	raw, nErr := reply.Bytes()
	if nErr != nil {
		return ferr.Wrap(nErr, "")
	}
	nErr = json.Unmarshal(raw, val)

	if nErr != nil {
		return ferr.Wrap(nErr, "")
	}
	return nil
}

// Sets "key" to json encoded "val"
func (db *Database) SetJson(key string, val interface{}) ferr.FortiaError {
	serialized, nErr := json.Marshal(val)
	if nErr != nil {
		return ferr.Wrap(nErr, "")
	}

	_, err := db.Cmd("SET", key, serialized)
	return err
}

// Find a better way to edit sets later
// Edit a set of integers
func (db *Database) EditSet(add []interface{}, del []interface{}) ferr.FortiaError {
	if len(add) > 1 || len(del) > 1 {
		client, err := db.Pool.Get()
		if err != nil {
			return ferr.Wrap(err, "")
		}
		defer db.Pool.Put(client)

		if len(add) > 1 {
			reply := client.Cmd("SADD", add...)
			if reply.Err != nil {
				return ferr.Wrap(reply.Err, "")
			}
		}

		if len(del) > 1 {
			reply := client.Cmd("SDEL", del...)
			if reply.Err != nil {
				return ferr.Wrap(reply.Err, "")
			}
		}
	}
	return nil
}

func (db *Database) EditSetI(key string, add, del []int) ferr.FortiaError {
	argAddSlice := make([]interface{}, len(add)+1)
	for k, v := range add {
		argAddSlice[k+1] = v
	}
	argAddSlice[0] = key

	argDelSlice := make([]interface{}, len(del)+1)
	for k, v := range del {
		argDelSlice[k+1] = v
	}
	argDelSlice[0] = key

	return db.EditSet(argAddSlice, argDelSlice)
}

// Edit a set of strings
func (db *Database) EditSetS(key string, add, del []string) ferr.FortiaError {
	argAddSlice := make([]interface{}, len(add)+1)
	for k, v := range add {
		argAddSlice[k+1] = v
	}
	argAddSlice[0] = key

	argDelSlice := make([]interface{}, len(del)+1)
	for k, v := range del {
		argDelSlice[k+1] = v
	}
	argDelSlice[0] = key

	return db.EditSet(argAddSlice, argDelSlice)
}
