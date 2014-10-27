package world

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
	Id           int
	Name         string
	Flags        []string
	Propbability int8
	Properties   BiomeProperties
}
