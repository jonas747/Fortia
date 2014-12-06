package world

import (
	"encoding/json"
	ferr "github.com/jonas747/fortia/error"
	"io/ioutil"
)

type BiomeProperties struct {
	Trees                     int8
	Roughness                 int8
	Soil                      int8
	Metals                    int8
	Wildlife                  int8
	TemperatureSeasonVariance int8
	Temperature               int8
	Rivers                    int8
	Water                     int8
	Caves                     int8
}

type Biome struct {
	Id          int
	Name        string
	Flags       []string
	Probability int8
	Properties  BiomeProperties
}

type BiomesInfo struct {
	BiomeFlags        []string // List of enabled biome flags
	DefaultProperties BiomeProperties
	Biomes            []Biome
}

func BiomesFromJson(data []byte) (BiomesInfo, ferr.FortiaError) {
	var bi BiomesInfo
	err := json.Unmarshal(data, &bi)
	if err != nil {
		return BiomesInfo{}, ferr.Wrap(err, "")
	}

	// Apply default propterties
	// Todo find an easier way to do this, reflection perhaps?
	for k, v := range bi.Biomes {
		if v.Properties.Caves == 0 {
			v.Properties.Caves = bi.DefaultProperties.Caves
		}
		if v.Properties.Trees == 0 {
			v.Properties.Trees = bi.DefaultProperties.Trees
		}
		if v.Properties.Roughness == 0 {
			v.Properties.Roughness = bi.DefaultProperties.Roughness
		}
		if v.Properties.Soil == 0 {
			v.Properties.Soil = bi.DefaultProperties.Soil
		}
		if v.Properties.Metals == 0 {
			v.Properties.Metals = bi.DefaultProperties.Metals
		}
		if v.Properties.Wildlife == 0 {
			v.Properties.Wildlife = bi.DefaultProperties.Wildlife
		}
		bi.Biomes[k] = v
	}

	return bi, nil
}

func BiomesFromFile(file string) (BiomesInfo, ferr.FortiaError) {
	data, nErr := ioutil.ReadFile(file)
	if nErr != nil {
		return BiomesInfo{}, ferr.Wrap(nErr, "")
	}
	biomes, err := BiomesFromJson(data)
	return biomes, err
}

func (b *BiomesInfo) GetBiomeFromId(id int) Biome {
	for _, v := range b.Biomes {
		if v.Id == id {
			return v
		}
	}
	return Biome{}
}
