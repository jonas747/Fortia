package main

import (
	"encoding/json"
	ferr "github.com/jonas747/fortia/error"
	"io/ioutil"
)

type Config struct {
	GameDb    string
	LogServer string
}

func loadConfig(fpath string) (*Config, ferr.FortiaError) {
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, ferr.Wrap(err, "Error loading config")
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, ferr.Wrap(err, "Error decoding config json")
	}
	return &config, nil
}

func loadBlockTypes(fpath string) ([]*BlockType, ferr.FortiaError) {
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, ferr.Wrap(err, "Error loading config")
	}

	var rawTypes []BlockType
	err = json.Unmarshal(data, &rawTypes)
	if err != nil {
		return nil, ferr.Wrap(err, "Error decoding config json")
	}

	// Make sure we get them int he right order
	blockTypes := make([]*BlockType, len(rawTypes)+1)
	for i, v := range rawTypes {
		blockTypes[i] = &v
	}

	return blockTypes, nil
}
