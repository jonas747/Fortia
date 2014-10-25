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
