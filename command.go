package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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
	if ctx == nil {
		return fmt.Errorf("no context passed")
	}

	if len(ctx.args) < 2 {
		return fmt.Errorf("insufficient arguments. usage: %s <location>", ctx.args[0])
	}

	var body []byte

	v, ok := cache.Get(pokeapi.LocationAreaURL + "/" + ctx.args[1])
	if !ok {
		res, err := http.Get(pokeapi.LocationAreaURL + "/" + ctx.args[1])
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

		cache.Add(pokeapi.LocationAreaURL+"/"+ctx.args[1], body)
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

func commandCatch(ctx *commandContext) error {
	if ctx == nil {
		return fmt.Errorf("no context passed")
	}

	if len(ctx.args) < 2 {
		return fmt.Errorf("insufficient arguments. usage: %s <pokemon>", ctx.args[0])
	}

	req := pokeapi.PokemonURL + "/" + ctx.args[1]

	var body []byte

	v, ok := cache.Get(req)
	if !ok {
		res, err := http.Get(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		by, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		cache.Add(req, by)

		body = by

	} else {
		body = v
	}

	if string(body) == "Not Found" {
		return fmt.Errorf("pokemon not found")
	}

	var pkm pokeapi.Pokemon

	if err := json.Unmarshal(body, &pkm); err != nil {
		return err
	}

	fmt.Println("Throwing a Pokeball at", pkm.Name+"...")

	if (rand.Intn(1000) - pkm.BaseExperience) > 500 {
		fmt.Println(pkm.Name, "was caught!")
		ctx.dex[pkm.Name] = pkm
	} else {
		fmt.Println(pkm.Name, "escaped!")
	}

	return nil
}

func commandInspect(ctx *commandContext) error {
	if ctx == nil {
		return fmt.Errorf("no context passed")
	}

	if len(ctx.args) < 2 {
		return fmt.Errorf("insufficient arguments. usage: %s <pokemon>", ctx.args[0])
	}

	pkm, ok := ctx.dex[ctx.args[1]]
	if !ok {
		fmt.Println("you have not caught a", ctx.args[1])
		return nil
	}

	fmt.Println("#", pkm.Id)
	fmt.Println("Name:", pkm.Name)
	fmt.Println("Height:", pkm.Height)
	fmt.Println("Weight:", pkm.Weight)
	fmt.Println("Stats:")
	for _, v := range pkm.Stats {
		fmt.Println("  -", v.Stat.Name, ":", v.Base_stat)
	}

	return nil
}

func commandPokedex(ctx *commandContext) error {
	if len(ctx.dex) < 1 {
		fmt.Println("You have no pokemon in your pokedex")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for _, v := range ctx.dex {
		fmt.Println(" -", v.Name)
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
		"pokedex": {
			name:        "pokedex",
			description: "Display your pokedex",
			callback:    commandPokedex,
		},
	}

	cmds.register("explore", "Shows the pokemon in a location. usage: explore <location>", commandExplore)
	cmds.register("catch", "Attempt to catch a pokemon. usage: catch <pokemon-name>", commandCatch)
	cmds.register("inspect", "Inspects a caught pokemon. usage: inspect <pokemon>", commandInspect)
}
