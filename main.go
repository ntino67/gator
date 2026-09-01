package main

import (
	"errors"
	"fmt"
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

	cfg.CurrentUserName = "ntino67"
	cfg.DbURL = "postgres://example"

	if err := cfg.Save(); err != nil {
		log.Fatalf("error setting the user: %v", err)
	}

	fmt.Printf("%s, %s", cfg.CurrentUserName, cfg.DbURL)
}
