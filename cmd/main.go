package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/sqmatheus/pokedexcli/internal/api"
	"github.com/sqmatheus/pokedexcli/internal/cli"
)

const factor = 20.0

var mapState = newState[api.LocationEntry](api.LocationUrl, "")

var pokedex map[string]api.Pokemon = make(map[string]api.Pokemon)

type state[T any] struct {
	Next     string
	Previous string
}

func newState[T any](next string, previous string) *state[T] {
	return &state[T]{Next: next, Previous: previous}
}

func (s *state[T]) Update(p api.Pagination[T]) {
	if p.Next != nil {
		s.Next = *p.Next
	} else {
		s.Next = ""
	}

	if p.Previous != nil {
		s.Previous = *p.Previous
	} else {
		s.Previous = ""
	}
}

func main() {
	cli.RegisterCommand(cli.NewCommand("help", "Displays a help message", func(_ []string) error {
		fmt.Println("Avaliable commands:")
		for _, cmd := range cli.CommandMap {
			fmt.Printf("\t%s\n", cmd.Usage())
		}
		fmt.Println()
		return nil
	}))

	cli.RegisterCommand(cli.NewCommand("exit", "Exit the pokedex", func(_ []string) error {
		fmt.Println("exited")
		os.Exit(1)
		return nil
	}))

	cli.RegisterCommand(cli.NewCommand("map", "Displays the names of 20 location areas in the Pokemon world", func(_ []string) error {
		if mapState.Next == "" {
			return errors.New("there are no more maps")
		}

		locations, err := api.GetLocations(mapState.Next)
		if err != nil {
			return err
		}

		mapState.Update(locations)
		for _, entry := range locations.Results {
			fmt.Println(entry.Name)
		}
		return nil
	}))

	cli.RegisterCommand(cli.NewCommand("mapb", "Displays the names of 20 location areas in the Pokemon world", func(_ []string) error {
		if mapState.Previous == "" {
			return errors.New("there are no previous maps")
		}

		locations, err := api.GetLocations(mapState.Previous)
		if err != nil {
			return err
		}

		mapState.Update(locations)
		for _, entry := range locations.Results {
			fmt.Println(entry.Name)
		}
		return nil
	}))

	cli.RegisterCommand(cli.NewCommand("explore", "Displays a list of all the Pokémon in a given area", func(params []string) error {
		if len(params) == 0 {
			return errors.New("provide an area 'explore <area>'")
		}
		name := params[0]
		fmt.Printf("Exploring %s...\n", name)

		area, err := api.GetLocationArea(name)
		if err != nil {
			return err
		}

		for _, encounters := range area.PokemonEncounters {
			fmt.Printf("\t- %s\n", encounters.Pokemon.Name)
		}

		return nil
	}))

	cli.RegisterCommand(cli.NewCommand("catch", "Catch Pokemon and adds them to the Pokedex", func(params []string) error {
		if len(params) == 0 {
			return errors.New("provide a pokemon 'catch <pokemon>'")
		}

		name := params[0]
		if _, ok := pokedex[name]; ok {
			fmt.Printf("you already have %s in your pokedex!\n", name)
			return nil
		}

		pokemon, err := api.GetPokemon(name)
		if err != nil {
			return err
		}

		be := float64(pokemon.BaseExperience)
		fmt.Printf("Chance: %.2f%%\n", (factor*100.0)/be)

		fmt.Printf("Throwing a Pokeball at %s...\n", name)
		if (rand.Float64() * be) >= factor {
			fmt.Printf("%s escaped!\n", name)
			return nil
		}

		fmt.Printf("%s was caught!\n", name)
		pokedex[name] = pokemon
		return nil
	}))

	cli.RegisterCommand(cli.NewCommand("inspect", "Displays the details about a Pokemon in your pokedex", func(params []string) error {
		if len(params) == 0 {
			return errors.New("provide a pokemon 'inspect <pokemon>'")
		}

		name := params[0]
		pokemon, ok := pokedex[name]
		if !ok {
			return errors.New("you have not caught that pokemon")
		}

		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, t := range pokemon.Types {
			fmt.Printf("  - %s\n", t.Type.Name)
		}
		return nil
	}))

	cli.RegisterCommand(cli.NewCommand("pokedex", "Displays a list of all the Pokemon in your pokedex", func(_ []string) error {
		if len(pokedex) == 0 {
			return errors.New("your pokedex is empty")
		}

		fmt.Println("Your Pokedex:")
		for name := range pokedex {
			fmt.Printf("\t- %s\n", name)
		}
		return nil
	}))

	cli.Start()
}
