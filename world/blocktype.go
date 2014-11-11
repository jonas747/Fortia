package world

import (
	"encoding/json"
	ferr "github.com/jonas747/fortia/error"
	//	"github.com/jonas747/fortia/vec"
	"io/ioutil"
	"strconv"
	"strings"
)

// The probability which this block can spawn
type BlockProbability struct {
	Everywhere int
	Outside    int
	Inside     int
	Biomes     map[string]int
}

// Represents a block type
type BlockType struct {
	Id        int
	Name      string
	Flags     []string
	Biomes    []string
	AllBiomes bool
	Type      string
	Layer     string

	Probability BlockProbability

	LayerStart   int
	LayerEnd     int
	LayerOutSide bool

	// Additional properties
	Properties map[string]interface{}
}

// Returns BlockTypes from json byte slice
func BlockTypesFromJson(data []byte) (BlockTypes, ferr.FortiaError) {
	// Decode the json
	var btypes []BlockType
	err := json.Unmarshal(data, &btypes)
	if err != nil {
		return nil, ferr.Wrap(err, "")
	}

	finalOut := make([]*BlockType, len(btypes))
	for i, v := range btypes {
		if v.LayerStart == 0 && v.LayerEnd == 0 {
			if v.Layer == "outside" {
				v.LayerOutSide = true
			} else if v.Layer == "inside" {
				v.LayerOutSide = false
			} else if v.Layer == "*" {
				v.LayerEnd = 1000
			} else if strings.Contains(v.Layer, "-") {
				split := strings.Split(v.Layer, "-")
				start, err := strconv.Atoi(split[0])
				if err != nil {
					return nil, ferr.Wrap(err, "")
				}
				end, err := strconv.Atoi(split[1])
				if err != nil {
					return nil, ferr.Wrap(err, "")
				}
				v.LayerStart = start
				v.LayerEnd = end
			}
		}

		if len(v.Biomes) == 0 || v.Biomes[0] == "*" {
			v.AllBiomes = true
		}
		finalOut[i] = &v
	}
	return BlockTypes(finalOut), nil
}

func BlockTypesFromFile(file string) (BlockTypes, ferr.FortiaError) {
	data, nErr := ioutil.ReadFile(file)
	if nErr != nil {
		return nil, ferr.Wrap(nErr, "")
	}
	btypes, err := BlockTypesFromJson(data)
	return btypes, err
}

// Returns an error with message 404 if not found
func (j *BlockType) GetPropertyInt(key string) (int, ferr.FortiaError) {
	interfaceVal, ok := j.Properties[key]
	if !ok {
		return 0, ferr.New("404")
	}

	value, ok := interfaceVal.(int)
	if !ok {
		return 0, ferr.New("Type mismatch")
	}
	return value, nil
}

func (j *BlockType) GetPropertyStr(key string) (string, ferr.FortiaError) {
	interfaceVal, ok := j.Properties[key]
	if !ok {
		return "", ferr.New("404")
	}

	value, ok := interfaceVal.(string)
	if !ok {
		return "", ferr.New("Type mismatch")
	}
	return value, nil
}

type BlockTypes []*BlockType

func (t BlockTypes) Get(index int) *BlockType {
	if index >= len(t) {
		return nil
	}
	return t[index]
}

/*
Returns a filtered BlockTypes
Available filters:

type			string					Filter by type
flags			[]string				By flags
properties		map[string]interface{}	By properties
haspropterties	[]string				By Property Keys
biomes			[]string				By biomes
layer			int						By layer
outside			bool					By outside
*/
func (t BlockTypes) Filter(filter map[string]interface{}) BlockTypes {
	out := make([]*BlockType, 0)
	for _, v := range []*BlockType(t) {
	FILTERLOOP:
		for filtername, filterval := range filter {
			switch filtername {
			case "type":
				str, _ := filterval.(string)
				if v.Type == str {
					out = append(out, v)
				}
			case "flags":
				slice, _ := filterval.([]string)
				for _, filterFlag := range slice {
					found := false
					for _, typeFlag := range v.Flags {
						if typeFlag == filterFlag {
							found = true
							break
						}
					}
					if !found {
						continue FILTERLOOP
					}
				}
				out = append(out, v)
			case "propterties":
				fMap := filterval.(map[string]interface{})
				for fKey, fVal := range fMap {
					pVal, ok := v.Properties[fKey]
					if !ok {
						continue FILTERLOOP
					}
					if pVal != fVal {
						continue FILTERLOOP
					}
				}
				out = append(out, v)
			case "hasProperties":
				slice := filterval.([]string)
				for _, fKey := range slice {
					_, ok := v.Properties[fKey]
					if !ok {
						continue FILTERLOOP
					}
				}
				out = append(out, v)
			case "biomes":
				slice := filterval.([]string)
				for _, biome := range slice {
					found := false

					for _, tBiome := range v.Biomes {
						if tBiome == biome {
							found = true
							break
						}
					}

					if !found {
						continue FILTERLOOP
					}
				}
				out = append(out, v)
			case "layer":
				lInt := filterval.(int)
				lstart := v.LayerStart
				lend := v.LayerEnd

				if lstart > lend {
					lstart = lend
					lend = v.LayerStart
				}

				if lInt <= lend || lInt >= lstart {
					out = append(out, v)
				}
			case "outside":
				isOutside := filterval.(bool)
				if isOutside == v.LayerOutSide {
					out = append(out, v)
				}
			}
		}
	}

	return BlockTypes(out)
}
