package common

import (
	"encoding/json"
	ferr "github.com/jonas747/fortia/error"
	"io/ioutil"
)

// Simple function for loading a json file
func LoadJsonFile(path string, out interface{}) ferr.FortiaError {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return ferr.Wrap(err, "")
	}

	err = json.Unmarshal(body, out)
	if err != nil {
		return ferr.Wrap(err, "")
	}
	return nil
}
