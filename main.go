package main

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"gitlab.com/vali19th/boot.dev_pokedex_cli/internal/pokeapi"
)

type cliCommand struct {
	callback    func(*config) error
	description string
}

type config struct {
	client   pokeapi.Client
	next_url *string
	prev_url *string
}

func main() {
	empty_str := ""
	cfg := &config{client: pokeapi.NewClient(), next_url: &empty_str, prev_url: nil}
	commands := get_commands()
	commands["h"].callback(cfg)

	for {
		fmt.Print("\nPokedex> ")
		var input string
		fmt.Scanln(&input)

		cmd, ok := commands[input]
		if !ok {
			fmt.Printf("Invalid command %q\n", input)
			continue
		}

		if err := cmd.callback(cfg); err != nil {
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
	}
}

func cmd_help(cfg *config) error {
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

func cmd_quit(cfg *config) error {
	os.Exit(0)
	return nil
}

func cmd_next(cfg *config) error {
	if cfg.next_url == nil {
		return errors.New("You are on the last page.")
	} else if *cfg.next_url == "" {
		cfg.next_url = nil
	}

	fmt.Println(cfg)
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

func cmd_prev(cfg *config) error {
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

