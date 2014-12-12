package world

import (
	"encoding/json"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/messages"
	"io/ioutil"
)

func BiomesFromJson(data []byte) (*messages.WorldBiomes, ferr.FortiaError) {
	var biomes *messages.WorldBiomes

	nErr := json.Unmarshal(data, biomes)
	if nErr != nil {
		return nil, ferr.Wrap(nErr, "")
	}

	return biomes, nil
}

func BiomesFromFile(file string) (*messages.WorldBiomes, ferr.FortiaError) {
	data, nErr := ioutil.ReadFile(file)
	if nErr != nil {
		return &messages.WorldBiomes{}, ferr.Wrap(nErr, "")
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
