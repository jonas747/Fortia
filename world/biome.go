package world

import (
	"encoding/json"
	ferr "github.com/jonas747/fortia/error"
	"io/ioutil"
)

type BiomeProperties struct {
	Trees                      int8
	Roughness                  int8
	Soil                       int8
	Metals                     int8
	Wildlife                   int8
	TempteratureSeasonVariance int8
	Temperature                int8
	Rivers                     int8
	Water                      int8
	Caves                      int8
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
