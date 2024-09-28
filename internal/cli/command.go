package cli

import "fmt"

type CommandCallback func(*Cli, []string) error

type Command struct {
	Name        string
	Description string
	Callback    CommandCallback
}

func NewCommand(name, description string, callback CommandCallback) Command {
	return Command{Name: name, Description: description, Callback: callback}
}

func (c Command) Usage() string {
	return fmt.Sprintf("%s - %s", c.Name, c.Description)
}
