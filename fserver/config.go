package main

import ()

type Config struct {
	Auth       string // auth database address
	Game       string // game db address
	Master     string // Master server address
	Version    string // The server version
	ServeGame  bool
	ServerAuth bool
}

type WorldConfig struct {
	LayerSize   int
	WorldHeight int
	Name        string
}
