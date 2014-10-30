package main

import ()

type Config struct {
	Auth      string // auth database address
	Game      string // game db address
	Master    string // Master server address
	Version   string // The server version
	ServeGame bool   // Serves the game api if true
	ServeAuth bool   // Serves the auth api if true
	RunTicker bool   // Runs the world ticker if true
}

type WorldConfig struct {
	LayerSize   int
	WorldHeight int
	Name        string
}
