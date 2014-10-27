package main

import (
	"errors"
	"github.com/jonas747/fortia/world"
	"strconv"
)

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
