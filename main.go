package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	var ret []string
	for v := range strings.FieldsSeq(text) {
		ret = append(ret, strings.ToLower(v))
	}

	return ret
}

func main() {
	userInput := bufio.NewScanner(os.Stdin) // gets blocking input, returns *Scanner with contained input
	ctx := newCommandContext()
	dex := newPokedex()
	ctx.dex = dex

	for {
		fmt.Print("Pokedex > ")

		if !userInput.Scan() {
			break
		}

		input := cleanInput(userInput.Text())
		if len(input) == 0 {
			input = append(input, "")
		}

		command, ok := cmds[input[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		ctx.setArgs(input)

		if err := command.callback(&ctx); err != nil {
			fmt.Printf("Error in command %s: %v\n", command.name, err)
		}
	}
}
