package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var CommandMap map[string]Command = make(map[string]Command)

func RegisterCommand(c Command) {
	CommandMap[c.Name] = c
}

func executeCommand(input string) error {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return errors.New("provide a command! use 'help'")
	}

	cmd, exists := CommandMap[fields[0]]
	if !exists {
		return errors.New("invalid command! use 'help'")
	}

	if err := cmd.Callback(fields[1:]); err != nil {
		return err
	}
	return nil
}

func Start() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("pokedex > ")
		scanner.Scan()

		if scanner.Err() != nil {
			fmt.Printf("scanner error: %v\n", scanner.Err())
			break
		}

		input := strings.TrimSpace(scanner.Text())
		err := executeCommand(input)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}

}
