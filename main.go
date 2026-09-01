package main

import (
	"errors"
	"log"
	"os"

	"github.com/ntino67/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if errors.Is(err, os.ErrNotExist) {
		newCfg := config.Config{}
		if err := newCfg.Save(); err != nil {
			log.Fatalf("error creating new confing: %v", err)
		}
		cfg = newCfg
	} else if err != nil {
		log.Fatalf("error reading the config file: %v", err)
	}

	s := state{cfg: &cfg}

	cmds := commands{registeredCommands: make(map[string]func(*state, command) error)}
	
	// When we add a command, we will register it here.
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatalf("Additional argmuments needed, expected: %d, provided: %d", 2, len(os.Args))
	}

	name := os.Args[1]
	args := os.Args[2:]
	
	cmd := command{name: name, args: args}
	if err := cmds.run(&s, cmd); err != nil {
		log.Fatalf("error running %s: %v", cmd.name, err)
	}
}
