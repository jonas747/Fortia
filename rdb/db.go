package rdb

/*
	This is a redis implementation of authserver.AuthDB and world.GameDB
*/

import (
	"github.com/fzzy/radix/extra/pool"
	"github.com/fzzy/radix/redis"
	ferr "github.com/jonas747/fortia/error"
)

var (
	emptyStrStrMap = make(map[string]string)
)

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
	client, err := db.Pool.Get()
	if err != nil {
		return nil, ferr.Wrap(err, "Error Get db client")
	}
	defer db.Pool.Put(client)
	reply := client.Cmd(cmd, args)
	if reply.Err != nil {
		return nil, ferr.Wrapa(reply.Err, "Redis cmd", map[string]interface{}{"cmd": cmd, "args": args})
	}
	return reply, nil
}

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

func (db *Database) SetHash(key string, info map[string]interface{}) ferr.FortiaError {
	_, err := db.Cmd("HMSET", key, info)
	return err
}
