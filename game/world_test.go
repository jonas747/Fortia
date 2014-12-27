package game

import (
	"github.com/golang/protobuf/proto"
	"github.com/jonas747/fortia/messages"
	"github.com/jonas747/fortia/vec"
	. "gopkg.in/check.v1"
	"testing"
)

type WorldSuite struct {
	W *World
}

func (w *WorldSuite) SetUpSuite(c *C) {
}

var _ = Suite(&WorldSuite{
	W: &World{
		Settings: &messages.WorldSettings{
			ChunkWidth:  proto.Int(25),
			ChunkHeight: proto.Int(200),
		},
	},
})

func Test(t *testing.T) { TestingT(t) }

func (w *WorldSuite) TestNDArray(c *C) {
	startCoords := vec.Vec3I{9, 23, 199}
	index := w.W.CoordsToIndex(startCoords)
	c.Assert(index, Equals, 124959)

	endCoords := w.W.IndexToCoords(index)
	c.Assert(endCoords.X, Equals, startCoords.X)
	c.Assert(endCoords.Y, Equals, startCoords.Y)
	c.Assert(endCoords.Z, Equals, startCoords.Z)
}
