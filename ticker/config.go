package main

import (
	"encoding/json"
	"github.com/jonas747/fortia/common"
	ferr "github.com/jonas747/fortia/error"
	"io/ioutil"
)

type Config struct {
	GameDb    string
	LogServer string
}

func loadBlockTypes(fpath string) ([]*BlockType, ferr.FortiaError) {
	var rawTypes []BlockType
	err := common.LoadJsonFile("blocktypes.json", out)

	// Make sure we get them int he right order
	blockTypes := make([]*BlockType, len(rawTypes)+1)
	for i, v := range rawTypes {
		blockTypes[i] = &v
	}

	return blockTypes, nil
}
