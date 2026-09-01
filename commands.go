package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/ntino67/gator/internal/config"
	"github.com/ntino67/gator/internal/database"
)

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

type state struct {
	cfg       *config.Config
	dbQueries *database.Queries
}

func (c *commands) run (s *state, cmd command) error {
	handler, ok := c.registeredCommands[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.name)
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
	queries := s.dbQueries
	ctx := context.Background()
	name := cmd.args[0]

	usr, err := queries.GetUser(ctx, name)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("no user with this name: %w", err)
	} else if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	cfg.SetUser(usr.Name)

	fmt.Printf("You are logged in as: %s", cfg.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	cfg := s.cfg
	queries := s.dbQueries
	ctx := context.Background()
	name := cmd.args[0]

	params := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
	}

	usr, err := queries.CreateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("couldn't create the user: %w", err)
	}

	cfg.SetUser(usr.Name)

	fmt.Printf("User %s was created\n", usr.Name)
	log.Printf("ID: %s, Name: %s", usr.ID, usr.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	fmt.Print("Resetted the database")
	return s.dbQueries.ResetUser(ctx)
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()
	
	usrs, err := s.dbQueries.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	for _, usr := range usrs {
		if usr.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", usr.Name)
			continue
		}
		fmt.Printf("* %s\n", usr.Name)
	}

	return nil
}
