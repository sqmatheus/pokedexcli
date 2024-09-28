package main

import (
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/sqmatheus/pokedexcli/internal/api"
	"github.com/sqmatheus/pokedexcli/internal/cli"
	"github.com/sqmatheus/pokedexcli/internal/util"
)

const factor = 20.0

var (
	mapState                        = util.NewState[api.LocationEntry](api.LocationUrl, "")
	pokedex  map[string]api.Pokemon = make(map[string]api.Pokemon)
)

func handleHelp(c *cli.Cli, _ []string) error {
	fmt.Println("Avaliable commands:")
	for _, cmd := range c.CommandMap {
		fmt.Printf("\t%s\n", cmd.Usage())
	}
	fmt.Println()
	return nil
}

func handleExit(c *cli.Cli, _ []string) error {
	fmt.Println("exited")
	c.Exit()
	return nil
}

func handleMap(_ *cli.Cli, _ []string) error {
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
}

func handleMapb(_ *cli.Cli, _ []string) error {
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
}

func handleExplore(_ *cli.Cli, params []string) error {
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
}

func handleCatch(_ *cli.Cli, params []string) error {
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
}

func handleInspect(_ *cli.Cli, params []string) error {
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
}

func handlePokedex(_ *cli.Cli, _ []string) error {
	if len(pokedex) == 0 {
		return errors.New("your pokedex is empty")
	}

	fmt.Println("Your Pokedex:")
	for name := range pokedex {
		fmt.Printf("\t- %s\n", name)
	}
	return nil
}

func main() {
	c := cli.NewCli()

	c.RegisterCommand(cli.NewCommand("help", "Displays a help message", handleHelp))
	c.RegisterCommand(cli.NewCommand("exit", "Exit the pokedex", handleExit))
	c.RegisterCommand(cli.NewCommand("map", "Displays the names of 20 location areas in the Pokemon world", handleMap))
	c.RegisterCommand(cli.NewCommand("mapb", "Displays the names of 20 location areas in the Pokemon world", handleMapb))
	c.RegisterCommand(cli.NewCommand("explore", "Displays a list of all the Pokémon in a given area", handleExplore))
	c.RegisterCommand(cli.NewCommand("catch", "Catch Pokemon and adds them to the Pokedex", handleCatch))
	c.RegisterCommand(cli.NewCommand("inspect", "Displays the details about a Pokemon in your pokedex", handleInspect))
	c.RegisterCommand(cli.NewCommand("pokedex", "Displays a list of all the Pokemon in your pokedex", handlePokedex))

	c.Start()
}
