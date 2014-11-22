package rdb

import (
	. "gopkg.in/check.v1"
	"testing"
)

type DBSuite struct {
	*Database
}

func (db *DBSuite) SetUpSuite(c *C) {
	database, err := NewDatabase(":6379")
	c.Assert(err, IsNil)

	db.Database = database
}

func (db *DBSuite) TearDownSuite(c *C) {
	db.Cmd("DEL", "testing")
}

var _ = Suite(&DBSuite{})

func Test(t *testing.T) { TestingT(t) }

func (db *DBSuite) TestCmd(c *C) {
	reply, err := db.Cmd("PING")
	c.Assert(err, IsNil)

	str, nErr := reply.Str()
	c.Assert(nErr, IsNil)
	c.Assert(str, Equals, "PONG")
}

func (db *DBSuite) TestPipelinedCmds(c *C) {
	cmds := []RedisCmd{
		RedisCmd{
			Cmd: "PING",
		},
		RedisCmd{
			Cmd:  "ECHO",
			Args: []interface{}{"lasagna"},
		},
	}

	replies, err := db.PipelinedCmds(cmds)
	c.Assert(err, IsNil)

	c.Assert(replies, HasLen, 2)

	rPing := replies[0]
	strPing, nErr := rPing.Str()
	c.Assert(nErr, IsNil)
	c.Assert(strPing, Equals, "PONG")

	rEcho := replies[1]
	strEcho, nErr := rEcho.Str()
	c.Assert(nErr, IsNil)
	c.Assert(strEcho, Equals, "lasagna")
}

func (db *DBSuite) TestGetSetHash(c *C) {
	hash := map[string]interface{}{
		"food":    "lasagna",
		"pet":     "dog",
		"country": "norway",
	}
	key := "testing"

	err := db.SetHash(key, hash)
	c.Assert(err, IsNil)

	hReply, err := db.GetHash(key)
	c.Assert(err, IsNil)

	for k, v := range hReply {
		c.Assert(v, Equals, hash[k])
	}
}
