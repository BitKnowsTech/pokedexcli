package main

import "github.com/bitknowstech/pokedexcli/internal/pokeapi"

type pokedex map[string]pokeapi.Pokemon

func newPokedex() pokedex {
	return pokedex{}
}
