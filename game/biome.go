package game

import (
	"encoding/json"
	"github.com/jonas747/fortia/errorcodes"
	"github.com/jonas747/fortia/errors"
	"github.com/jonas747/fortia/messages"
	"io/ioutil"
)

func BiomesFromJson(data []byte) (*messages.WorldBiomes, errors.FortiaError) {
	var biomes *messages.WorldBiomes

	nErr := json.Unmarshal(data, biomes)
	if nErr != nil {
		return nil, errors.Wrap(nErr, errorcodes.ErrorCode_JsonDecodeErr, "", nil)
	}

	return biomes, nil
}

func BiomesFromFile(file string) (*messages.WorldBiomes, errors.FortiaError) {
	data, nErr := ioutil.ReadFile(file)
	if nErr != nil {
		return &messages.WorldBiomes{}, errors.Wrap(nErr, errorcodes.ErrorCode_FileReadErr, "", nil)
	}
	biomes, err := BiomesFromJson(data)
	return biomes, err
}

/*
func (b *BiomesInfo) GetBiomeFromId(id int) Biome {
	for _, v := range b.Biomes {
		if v.Id == id {
			return v
		}
	}
	return Biome{}
}
*/
