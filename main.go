package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"gitlab.com/vali19th/boot.dev_pokedex_cli/internal/pokeapi"
)

type cliCommand struct {
	callback    func(*config, ...string) error
	description string
}

type config struct {
	client   pokeapi.Client
	next_url *string
	prev_url *string
	pokedex  map[string]pokeapi.Pokemon
}

func main() {
	empty_str := ""
	cfg := &config{
		client:   pokeapi.NewClient(time.Hour),
		next_url: &empty_str,
		prev_url: nil,
		pokedex:  make(map[string]pokeapi.Pokemon),
	}
	commands := get_commands()
	commands["h"].callback(cfg)

	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nPokedex> ")
		reader.Scan()
		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		name := words[0]
		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}

		cmd, ok := commands[name]
		if !ok {
			fmt.Printf("Invalid command %q\n", name)
			continue
		}

		if err := cmd.callback(cfg, args...); err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	}
}

func get_commands() map[string]cliCommand {
	return map[string]cliCommand{
		"h": {callback: cmd_help, description: "Display a help message"},
		"q": {callback: cmd_quit, description: "Quit"},
		"n": {callback: cmd_next, description: "Display the next page of locations"},
		"p": {callback: cmd_prev, description: "Display the previous page of locations"},
		"e": {callback: cmd_explore, description: "Display pokemons in a location area"},
		"c": {callback: cmd_catch, description: "Attempt to catch a pokemon and add it to the pokedex"},
		"i": {callback: cmd_inspect, description: "View information about a caught pokemon"},
		"l": {callback: cmd_list_caught_pokemons, description: "List all caught pokemons"},
	}
}

func cmd_help(cfg *config, args ...string) error {
	commands := get_commands()

	// Sort by name
	sorted := make([]string, 0, len(commands))
	for name := range commands {
		sorted = append(sorted, name)
	}
	sort.Strings(sorted)

	// Generate the help message
	help := "Available commands:\n"
	for _, name := range sorted {
		cmd := commands[name]
		help += fmt.Sprintf("    %s: %s\n", name, cmd.description)
	}

	fmt.Print(help)
	return nil
}

func cmd_quit(cfg *config, args ...string) error {
	os.Exit(0)
	return nil
}

func cmd_next(cfg *config, args ...string) error {
	if cfg.next_url == nil {
		return errors.New("You are on the last page.")
	} else if *cfg.next_url == "" {
		cfg.next_url = nil
	}

	resp, err := cfg.client.GetLocationAreas(cfg.next_url)
	if err != nil {
		return err
	}

	cfg.next_url = resp.Next
	cfg.prev_url = resp.Previous

	msg := "Location areas:\n"
	for _, area := range resp.Results {
		msg += fmt.Sprintf("- %s\n", area.Name)
	}

	fmt.Print(msg)
	return nil
}

func cmd_prev(cfg *config, args ...string) error {
	if cfg.prev_url == nil {
		return errors.New("You are on the first page.")
	}

	resp, err := cfg.client.GetLocationAreas(cfg.prev_url)
	if err != nil {
		return err
	}

	cfg.next_url = resp.Next
	cfg.prev_url = resp.Previous

	msg := "Location areas:\n"
	for _, area := range resp.Results {
		msg += fmt.Sprintf("- %s\n", area.Name)
	}

	fmt.Print(msg)
	return nil
}

func cmd_explore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("No location area provided")
	}

	name := args[0]
	locArea, err := cfg.client.GetLocationArea(name)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Pokemons in %s:\n", name)
	for _, pokemon := range locArea.PokemonEncounters {
		msg += fmt.Sprintf("- %s\n", pokemon.Pokemon.Name)
	}

	fmt.Print(msg)
	return nil
}

func cmd_catch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("No pokemon name provided")
	}

	name := args[0]
	pokemon, err := cfg.client.GetPokemon(name)
	if err != nil {
		return err
	}

	const threshold = 50
	n := rand.Intn(pokemon.BaseExperience)
	if n > threshold {
		return fmt.Errorf("failed to catch %s", name)
	}

	cfg.pokedex[name] = pokemon
	fmt.Printf("%s was caught!\n", name)
	return nil
}

func cmd_inspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("No pokemon name provided")
	}

	name := args[0]
	pokemon, ok := cfg.pokedex[name]
	if !ok {
		return errors.New("You haven't caught this pokemon yet")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" - %s: %v\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf(" - %s\n", t.Type.Name)
	}

	return nil
}

func cmd_list_caught_pokemons(cfg *config, args ...string) error {
	fmt.Println("Pokemons in pokedex:")
	for _, pokemon := range cfg.pokedex {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

