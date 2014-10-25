package main

import (
	"errors"
	"github.com/jonas747/fortia/world"
	"strconv"
)

var (
	ErrPropertyNotFound = errors.New("Property not found")
)

type BlockType struct {
	Id         int
	Name       string
	Flags      []string
	Properties map[string]string
}

func (j *BlockType) GetPropertyInt(key string) (value int, err error) {
	value = -1
	err = nil

	strVal, ok := j.Properties[key]
	if !ok {
		err = ErrPropertyNotFound
		return
	}

	value, err = strconv.Atoi(strVal)
	return
}

type Block struct {
	*world.Block
	kind *BlockType
}

func BlockFromWorldBlock(wb *world.Block) *Block {
	b := &Block{
		wb,
		blockTypes[wb.Id],
	}
	return b
}

func TickFluid(b *Block) {

}

func TickConvert(b *Block) {

}
