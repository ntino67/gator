package main

import (
	"errors"
	"fmt"

	"github.com/ntino67/gator/internal/config"
)

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

type state struct {
	cfg *config.Config
}

func (c *commands) run (s *state, cmd command) error {
	handler, ok := c.registeredCommands[cmd.name]
	if !ok {
		return errors.New("command not found")
	}
	return handler(s, cmd)
}

func (c *commands) register (name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	cfg := s.cfg
	cfg.CurrentUserName = cmd.args[0]

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("error saving: %w", err)
	}
	
	fmt.Printf("User has been set to: %s", cfg.CurrentUserName)
	return nil
}
