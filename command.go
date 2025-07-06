package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/bitknowstech/pokedexcli/internal/pokeapi"
	"github.com/bitknowstech/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*commandContext) error
}

type commands map[string]cliCommand

var cache pokecache.Cache = *pokecache.NewCache(time.Duration(5 * time.Second))
var cmds commands

func (c commands) register(name string, description string, callback func(*commandContext) error) {
	c[name] = cliCommand{name: name, description: description, callback: callback}
}

func commandExit(_ *commandContext) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *commandContext) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for k, v := range cmds {
		fmt.Println(k+":", v.description)
	}
	return nil
}

func commandExplore(ctx *commandContext) error {
	if ctx.args == nil {
		return fmt.Errorf("insufficient arguments. usage: explore <location>")
	}

	var body []byte

	v, ok := cache.Get(pokeapi.LocationAreaURL + "/" + ctx.args[0])
	if !ok {
		res, err := http.Get(pokeapi.LocationAreaURL + "/" + ctx.args[0])
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return fmt.Errorf("api responded %d", res.StatusCode)
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		cache.Add(pokeapi.LocationAreaURL+"/"+ctx.args[0], body)
	} else {
		body = v
	}

	var location pokeapi.LocationArea

	if err := json.Unmarshal(body, &location); err != nil {
		return err
	}

	for _, v := range location.PokemonEncounters {
		fmt.Println(v.Pokemon.Name)
	}

	return nil
}

func init() {
	cmds = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exits the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}

	cmds.register("explore", "Shows the pokemon in a location. usage: explore <location>", commandExplore)
}
