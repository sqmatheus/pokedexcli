package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/peterh/liner"
)

type Cli struct {
	ctx        context.Context
	cancel     context.CancelFunc
	CommandMap map[string]Command
}

func NewCli() *Cli {
	ctx, cancel := context.WithCancel(context.Background())
	return &Cli{
		ctx:        ctx,
		cancel:     cancel,
		CommandMap: make(map[string]Command),
	}
}

func (c *Cli) RegisterCommand(cmd Command) {
	c.CommandMap[cmd.Name] = cmd
}

func (c *Cli) executeCommand(input string) error {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return errors.New("provide a command! use 'help'")
	}

	cmd, exists := c.CommandMap[fields[0]]
	if !exists {
		return errors.New("invalid command! use 'help'")
	}

	if err := cmd.Callback(c, fields[1:]); err != nil {
		return err
	}
	return nil
}

func (c *Cli) Exit() {
	c.cancel()
}

func (c *Cli) Start() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	line.SetCompleter(func(line string) (comp []string) {
		for name := range c.CommandMap {
			if strings.HasPrefix(name, strings.ToLower(line)) {
				comp = append(comp, name)
			}
		}
		return
	})

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			cmd, err := line.Prompt("pokedex > ")
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return
			}

			input := strings.TrimSpace(cmd)
			if input != "" {
				line.AppendHistory(input)
			}
			err = c.executeCommand(input)
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
		}
	}
}
