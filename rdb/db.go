package rdb

/*
	This is a redis implementation of authserver.AuthDB and world.GameDB
*/
import (
	"encoding/json"
	"fmt"
	"github.com/fzzy/radix/extra/pool"
	"github.com/fzzy/radix/redis"
	"github.com/golang/protobuf/proto"
	"github.com/jonas747/fortia/db"
	"github.com/jonas747/fortia/errorcodes"
	"github.com/jonas747/fortia/errors"
	"strings"
)

var (
	emptyStrStrMap = make(map[string]string)
)

// Base database
type Database struct {
	Pool *pool.Pool
}

func NewDatabase(addr string) (*Database, error) {
	p, err := pool.NewPool("tcp", addr, 200)
	if err != nil {
		return nil, err
	}
	return &Database{p}, nil
}

func (rdb *Database) CheckReplyErrors(query string, reply *redis.Reply) errors.FortiaError {
	switch reply.Type {
	case redis.ErrorReply:
		// Error writing req or reading response
		return db.NewDBError(query, errorcodes.ErrorCode_RedisWriteReadErr, "Error writing or reading")
	case redis.NilReply:
		// Key not found
		return db.NewDBError(query, errorcodes.ErrorCode_RedisKeyNotFound, "Key not foudn")
	default:
		return nil
	}
}

// Same as redis.Client.Cmd but uses a connection from a pool
func (rdb *Database) Cmd(cmd string, args ...interface{}) (*redis.Reply, errors.FortiaError) {
	querySlice := make([]string, len(args)+1)
	for k, arg := range args {
		querySlice[k+1] = fmt.Sprint(arg)
	}
	querySlice[0] = cmd
	query := strings.Join(querySlice, " ")

	client, nErr := rdb.Pool.Get()
	if nErr != nil {
		err := db.NewDBError(query, errorcodes.ErrorCode_RedisDialError, nErr.Error())
		return nil, err
	}
	defer rdb.Pool.Put(client)

	reply := client.Cmd(cmd, args)
	if err := rdb.CheckReplyErrors(query, reply); err != nil {
		return reply, err
	}

	return reply, nil
}

// Represents a redis command
type RedisCmd struct {
	Cmd  string
	Args []interface{}
}

// Convenience methods

// Helper to get fetch value @ key and convert it to string
func (rdb *Database) GetString(cmd string, args ...interface{}) (string, errors.FortiaError) {
	reply, err := rdb.Cmd(cmd, args...)
	if err != nil {
		return "", err
	}
	str, errConv := reply.Str()
	if errConv != nil {
		return str, db.NewDBError(fmt.Sprint(cmd, args), errorcodes.ErrorCode_RedisReplyConversionErr, errConv.Error())
	}
	return str, nil
}

// Helper to get fetch value @ key and convert it to byte slice
func (rdb *Database) GetBytes(cmd string, args ...interface{}) ([]byte, errors.FortiaError) {
	reply, err := rdb.Cmd(cmd, args...)
	if err != nil {
		return []byte{}, err
	}
	bslice, errConv := reply.Bytes()
	if errConv != nil {
		return bslice, db.NewDBError(fmt.Sprint(cmd, args), errorcodes.ErrorCode_RedisReplyConversionErr, errConv.Error())
	}
	return bslice, nil
}

// Helper to get fetch value @ key and convert it to int64
func (rdb *Database) GetInt64(cmd string, args ...interface{}) (int64, errors.FortiaError) {
	reply, err := rdb.Cmd(cmd, args...)
	if err != nil {
		return 0, err
	}
	i64, errConv := reply.Int64()
	if errConv != nil {
		return i64, db.NewDBError(fmt.Sprint(cmd, args), errorcodes.ErrorCode_RedisReplyConversionErr, errConv.Error())
	}
	return i64, nil
}

// Same GetInt64 but converted to int
func (rdb *Database) GetInt(cmd string, args ...interface{}) (int, errors.FortiaError) {
	i64, err := rdb.GetInt64(cmd, args...)
	if err != nil {
		return 0, err
	}
	return int(i64), nil
}

// Helper to get fetch value @ key and convert it to bool
func (rdb *Database) GetBool(cmd string, args ...interface{}) (bool, errors.FortiaError) {
	reply, err := rdb.Cmd(cmd, args...)
	if err != nil {
		return false, err
	}
	rbool, errConv := reply.Bool()
	if errConv != nil {
		return rbool, db.NewDBError(fmt.Sprint(cmd, args), errorcodes.ErrorCode_RedisReplyConversionErr, errConv.Error())
	}
	return rbool, nil
}

// Helper to get fetch value @ key and convert it to string slice
func (rdb *Database) GetList(cmd string, args ...interface{}) ([]string, errors.FortiaError) {
	reply, err := rdb.Cmd(cmd, args...)
	if err != nil {
		return []string{}, err
	}

	list, errConv := reply.List()
	if errConv != nil {
		return list, db.NewDBError(fmt.Sprint(cmd, args), errorcodes.ErrorCode_RedisReplyConversionErr, errConv.Error())
	}
	return list, nil
}

// Does HGETALL on "key"
func (rdb *Database) GetHash(key string) (map[string]string, errors.FortiaError) {
	reply, err := rdb.Cmd("HGETALL", key)
	if err != nil {
		return emptyStrStrMap, err
	}
	hMap, errConv := reply.Hash()
	if errConv != nil {
		return emptyStrStrMap, db.NewDBError(fmt.Sprint("HGETALL", key), errorcodes.ErrorCode_RedisReplyConversionErr, errConv.Error())
	}
	return hMap, err
}

// Does HMSET on "key" with all the fialds in info map
func (rdb *Database) SetHash(key string, info map[string]interface{}) errors.FortiaError {
	_, err := rdb.Cmd("HMSET", key, info)
	return err
}

// Decodes the json at "key" into "val", returns any errors if any
// val and err is nil if not found
func (rdb *Database) GetJson(key string, val interface{}) errors.FortiaError {
	raw, err := rdb.GetBytes("GET", key)
	if err != nil {
		return err
	}

	nErr := json.Unmarshal(raw, val)
	if nErr != nil {
		return errors.Wrap(nErr, errorcodes.ErrorCode_JsonDecodeErr, "", nil)
	}
	return nil
}

// Sets "key" to json encoded "val"
func (rdb *Database) SetJson(key string, val interface{}) errors.FortiaError {
	serialized, nErr := json.Marshal(val)
	if nErr != nil {
		return errors.Wrap(nErr, errorcodes.ErrorCode_JsonEncodeErr, "", nil)
	}

	_, err := rdb.Cmd("SET", key, serialized)
	return err
}

func (rdb *Database) GetProto(key string, out proto.Message) errors.FortiaError {
	raw, err := rdb.GetBytes("GET", key)
	if err != nil {
		return err
	}
	nErr := proto.Unmarshal(raw, out)
	if nErr != nil {
		return errors.Wrap(nErr, errorcodes.ErrorCode_ProtoDecodeErr, "", nil)
	}
	return nil
}

func (rdb *Database) SetProto(key string, pb proto.Message) errors.FortiaError {
	serialized, nErr := proto.Marshal(pb)
	if nErr != nil {
		return errors.Wrap(nErr, errorcodes.ErrorCode_ProtoEncodeErr, "", nil)
	}
	_, err := rdb.Cmd("SET", key, serialized)
	return err
}

// Convenience methods for modifying sets
// Replace a set with the provided one
func (rdb *Database) SetSet(key string, set []interface{}) errors.FortiaError {
	if len(set) < 1 {
		return nil // Perhaps return an error here?
	}

	// Prepare argument slice
	args := make([]interface{}, 1)
	args[0] = key
	args = append(args, set...)

	// Delete it first
	_, err := rdb.Cmd("DEL", key)
	if err != nil {
		return err
	}

	_, err = rdb.Cmd("SADD", args...)
	return err
}

// Edit a set
func (rdb *Database) EditSet(add []interface{}, del []interface{}) errors.FortiaError {
	if len(add) > 1 || len(del) > 1 {
		if len(add) > 1 {
			_, err := rdb.Cmd("SADD", add...)
			if err != nil {
				return err
			}
		}

		if len(del) > 1 {
			_, err := rdb.Cmd("SDEL", del...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Edit a set of integers
func (rdb *Database) EditSetI(key string, add, del []int) errors.FortiaError {
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

	return rdb.EditSet(argAddSlice, argDelSlice)
}

// Edit a set of strings
func (rdb *Database) EditSetS(key string, add, del []string) errors.FortiaError {
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

	return rdb.EditSet(argAddSlice, argDelSlice)
}
