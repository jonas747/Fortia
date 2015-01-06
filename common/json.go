package common

import (
	"encoding/json"
	"github.com/jonas747/fortia/errorcodes"
	"github.com/jonas747/fortia/errors"
	"io/ioutil"
)

// Simple function for loading a json file
func LoadJsonFile(path string, out interface{}) errors.FortiaError {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New(errorcodes.ErrorCode_FileReadErr, "Error readind file %s, %s", path, err.Error())
	}

	err = json.Unmarshal(body, out)
	if err != nil {
		return errors.New(errorcodes.ErrorCode_JsonDecodeErr, "Error decoding json: %s", err.Error())
	}
	return nil
}
