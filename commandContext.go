package main

import "github.com/bitknowstech/pokedexcli/internal/pokeapi"

type commandContext struct {
	mapNext string
	mapPrev string
	args    []string
	dex     pokedex
}

func (cc *commandContext) setArgs(args []string) {
	cc.args = args
}

func newCommandContext() commandContext {
	return commandContext{mapNext: pokeapi.LocationAreaURL + "?offset=0&limit=20"}
}
