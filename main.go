package main

import (
	"fmt"
	"log"

	"github.com/ntino67/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading the config file: %v", err)
	}

	if cfg.CurrentUserName == "" {
		if err := cfg.Set(); err != nil {
			log.Fatalf("error setting the user: %v", err)
		}
	}

	fmt.Printf("%s, %s", cfg.CurrentUserName, cfg.DbURL)
}
