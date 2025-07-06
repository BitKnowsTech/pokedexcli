package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bitknowstech/pokedexcli/internal/pokeapi"
)

func commandMap(ctx *commandContext) error {
	if ctx == nil {
		return fmt.Errorf("no context passed")
	}

	if ctx.mapNext == "" {
		fmt.Println("you're on the last page")
		return nil
	}

	next, prev, err := mapReq(ctx.mapNext)
	if err != nil {
		return err
	}

	ctx.mapNext = next
	ctx.mapPrev = prev
	return nil
}

func commandMapb(ctx *commandContext) error {
	if ctx == nil {
		return fmt.Errorf("no context passed")
	}

	if ctx.mapPrev == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	next, prev, err := mapReq(ctx.mapPrev)
	if err != nil {
		return err
	}

	ctx.mapNext = next
	ctx.mapPrev = prev
	return nil
}

func mapReq(url string) (next, prev string, err error) {
	var body []byte

	v, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return "", "", err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return "", "", err
		}
		cache.Add(url, body)
	} else {
		body = v
	}

	var decBody pokeapi.NamedAPIResourceList[pokeapi.LocationArea]

	if err := json.Unmarshal(body, &decBody); err != nil {
		return "", "", err
	}

	if decBody.Results == nil {
		return "", "", fmt.Errorf("no results")
	}

	locations := decBody.Results

	for _, loc := range locations {
		fmt.Println(loc.Name)
	}

	return decBody.Next, decBody.Previous, nil
}

func init() {
	cmds.register("map", "Returns the next 20 map locations", commandMap)
	cmds.register("mapb", "Returns the previous 20 map locations", commandMapb)
}
